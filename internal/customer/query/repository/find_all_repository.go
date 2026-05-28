package repository

import (
	"context"
	"regexp"

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
		bsonFilter["$or"] = []bson.M{
			{"name": bson.M{"$regex": "^" + escapedSearch}},
			{"phone": bson.M{"$regex": "^" + escapedSearch}},
			{"product": bson.M{"$regex": "^" + escapedSearch}},
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

	// Pagination — ใช้เฉพาะเมื่อ Limit > 0
	// ถ้า Limit = 0 → return ทุก document (ใช้สำหรับ form dropdown / branch map)
	if filter.Limit > 0 {
		page := filter.Page
		if page < 1 {
			page = 1
		}
		limit := filter.Limit
		if limit > 100 {
			limit = 100 // hard cap
		}
		opts.SetSkip(int64((page - 1) * limit))
		opts.SetLimit(int64(limit))
	}

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
