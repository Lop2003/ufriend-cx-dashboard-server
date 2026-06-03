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
		escapedSearch := regexp.QuoteMeta(filter.Search)

		// Create a flexible phone search regex by extracting digits
		var phoneDigits strings.Builder
		for _, ch := range filter.Search {
			if ch >= '0' && ch <= '9' {
				phoneDigits.WriteRune(ch)
			}
		}

		var phoneRegex string
		if phoneDigits.Len() > 0 {
			var sb strings.Builder
			for i, r := range phoneDigits.String() {
				if i > 0 {
					sb.WriteString(`[^\d]*`)
				}
				sb.WriteRune(r)
			}
			phoneRegex = sb.String()
		} else {
			phoneRegex = escapedSearch
		}

		bsonFilter["$or"] = []bson.M{
			{"name": bson.M{"$regex": escapedSearch, "$options": "i"}},
			{"phone": bson.M{"$regex": phoneRegex, "$options": "i"}},
			{"product": bson.M{"$regex": escapedSearch, "$options": "i"}},
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
