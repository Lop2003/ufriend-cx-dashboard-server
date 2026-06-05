package repository

import (
	"context"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ufriend-cx-dashboard-server/internal/customer/model"
	"ufriend-cx-dashboard-server/internal/customer/query/dto"
)

func (r *customerQueryRepository) FindAll(ctx context.Context, filter *dto.CustomerFilter) ([]*model.Customer, int64, error) {
	collection := r.db.Collection("customers")

	bsonFilter := bson.M{}
	if filter.Branch != "" {
		bsonFilter["branch"] = filter.Branch
	}
	if filter.Status != "" {
		bsonFilter["status"] = filter.Status
	}
	if filter.Search != "" {
		// CRITICAL: On 10M+ records, $or across multiple fields causes full index scans → timeout.
		// Solution: detect input type and route to a SINGLE indexed field per query.
		search := filter.Search

		if containsThai(search) {
			// Thai text → name search using TEXT INDEX ($text)
			// Word-based matching: "ใจดี" matches "สมชาย ใจดี" (finds last names)
			bsonFilter["$text"] = bson.M{"$search": search}
		} else if search[0] >= '0' && search[0] <= '9' {
			// Starts with digit → phone search with formatted ^prefix
			var digits strings.Builder
			for _, ch := range search {
				if ch >= '0' && ch <= '9' {
					digits.WriteRune(ch)
				}
			}
			d := digits.String()
			var formatted string
			switch {
			case len(d) <= 3:
				formatted = d
			case len(d) <= 6:
				formatted = d[:3] + "-" + d[3:]
			default:
				formatted = d[:3] + "-" + d[3:6] + "-" + d[6:]
			}
			bsonFilter["phone"] = bson.M{"$regex": "^" + regexp.QuoteMeta(formatted)}
		} else {
			// English/Latin → product ^prefix search, case-insensitive
			escapedSearch := regexp.QuoteMeta(search)
			bsonFilter["product"] = bson.M{"$regex": "^" + escapedSearch, "$options": "i"}
		}
	}

	// Dynamic Bounded Count to support smooth deep pagination while keeping DB performance safe
	var total int64
	var err error
	if len(bsonFilter) == 0 {
		total, err = collection.EstimatedDocumentCount(ctx)
	} else {
		page := filter.Page
		if page < 1 {
			page = 1
		}
		limit := filter.Limit
		if limit < 1 {
			limit = 10
		}
		if limit > 100 {
			limit = 100
		}

		// Cap count at 10,000 by default, or higher if the user is deep paginating
		countLimit := int64(10000)
		currentPageEnd := int64(page) * int64(limit)
		if currentPageEnd >= countLimit {
			countLimit = currentPageEnd + int64(limit*10)
		}

		// Use pipeline with limit to ensure fast execution on 10M+ records
		pipeline := mongo.Pipeline{
			{{Key: "$match", Value: bsonFilter}},
			{{Key: "$limit", Value: countLimit}},
			{{Key: "$count", Value: "count"}},
		}
		cursorCount, errCount := collection.Aggregate(ctx, pipeline)
		if errCount == nil {
			defer cursorCount.Close(ctx)
			var countResults []struct {
				Count int64 `bson:"count"`
			}
			if errCount = cursorCount.All(ctx, &countResults); errCount == nil && len(countResults) > 0 {
				total = countResults[0].Count
			}
		} else {
			err = errCount
		}
	}
	if err != nil {
		return nil, 0, err
	}

	opts := options.Find()

	sortBy := "created_at"
	sortOrder := -1

	allowedSortFields := map[string]string{
		"name":        "name",
		"product":     "product",
		"plan_months": "plan_months",
		"status":      "status",
		"created_at":  "created_at",
	}

	if field, ok := allowedSortFields[filter.SortBy]; ok {
		sortBy = field
	}

	if filter.SortOrder == "asc" {
		sortOrder = 1
	} else if filter.SortOrder == "desc" {
		sortOrder = -1
	}

	opts.SetSort(bson.D{{Key: sortBy, Value: sortOrder}})

	// Pagination — บังคับ limit เพื่อป้องกัน response ขนาดใหญ่เกินไป
	// Default limit = 20 (consistent กับ feedback domain)
	page := filter.Page
	if page < 1 {
		page = 1
	}
	limit := filter.Limit
	if limit < 1 {
		limit = 20 // default limit (consistent กับ feedback)
	}
	if limit > 100 {
		limit = 100 // hard cap
	}
	opts.SetSkip(int64((page - 1) * limit))
	opts.SetLimit(int64(limit))

	cursor, err := collection.Find(ctx, bsonFilter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var customers = []*model.Customer{}
	if err := cursor.All(ctx, &customers); err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

// containsThai checks if a string contains any Thai character (Unicode range U+0E00–U+0E7F).
// Used to detect input type for optimal search index routing.
func containsThai(s string) bool {
	for _, r := range s {
		if r >= 0x0E00 && r <= 0x0E7F {
			return true
		}
	}
	return false
}
