package repository

import (
	"context"
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/follow_up/query/dto"
)

type facetResult struct {
	Metadata []struct {
		Total int64 `bson:"total"`
	} `bson:"metadata"`
	Data []struct {
		Id         primitive.ObjectID `bson:"_id"`
		CustomerId primitive.ObjectID `bson:"customer_id"`
		Type       string             `bson:"type"`
		Note       string             `bson:"note"`
		Status     string             `bson:"status"`
		CreatedAt  primitive.DateTime `bson:"created_at"`
		Customer   *struct {
			Name    string `bson:"name"`
			Phone   string `bson:"phone"`
			Product string `bson:"product"`
			Branch  string `bson:"branch"`
		} `bson:"customer"`
	} `bson:"data"`
}

func (r *followUpQueryRepository) FindAll(ctx context.Context, filter *dto.FollowUpFilter) ([]*dto.FollowUpResponse, int64, error) {
	// Calculate pagination skip/limit
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
	skip := int64((page - 1) * limit)

	// Build match filter
	matchStage := bson.M{}
	if filter.Type != "" {
		matchStage["type"] = filter.Type
	}
	if filter.Status != "" {
		matchStage["status"] = filter.Status
	}
	if filter.Branch != "" {
		matchStage["customer.branch"] = filter.Branch
	}
	if filter.Search != "" {
		escapedSearch := regexp.QuoteMeta(filter.Search)
		matchStage["$or"] = []bson.M{
			{"customer.name": bson.M{"$regex": "^" + escapedSearch, "$options": "i"}},
			{"customer.phone": bson.M{"$regex": "^" + escapedSearch, "$options": "i"}},
			{"note": bson.M{"$regex": "^" + escapedSearch, "$options": "i"}},
		}
	}

	pipeline := mongo.Pipeline{
		// 1. Join with customers collection
		{{Key: "$lookup", Value: bson.M{
			"from":         "customers",
			"localField":   "customer_id",
			"foreignField": "_id",
			"as":           "customer",
		}}},
		// Unwind joined customer
		{{Key: "$unwind", Value: bson.M{
			"path": "$customer",
			"preserveNullAndEmptyArrays": true,
		}}},
		// 2. Apply filters
		{{Key: "$match", Value: matchStage}},
		// 3. Facet search for metadata count and data pagination
		{{Key: "$facet", Value: bson.M{
			"metadata": bson.A{
				bson.M{"$count": "total"},
			},
			"data": bson.A{
				bson.M{"$sort": bson.M{"created_at": -1}},
				bson.M{"$skip": skip},
				bson.M{"$limit": limit},
			},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, 0, fmt.Errorf("aggregate follow ups list: %w", err)
	}
	defer cursor.Close(ctx)

	var facetResults []facetResult
	if err := cursor.All(ctx, &facetResults); err != nil {
		return nil, 0, fmt.Errorf("decode follow ups list: %w", err)
	}

	total := int64(0)
	var items []*dto.FollowUpResponse = []*dto.FollowUpResponse{}

	if len(facetResults) > 0 {
		fRes := facetResults[0]
		if len(fRes.Metadata) > 0 {
			total = fRes.Metadata[0].Total
		}

		for _, item := range fRes.Data {
			res := &dto.FollowUpResponse{
				Id:         item.Id.Hex(),
				CustomerId: item.CustomerId.Hex(),
				Type:       item.Type,
				Note:       item.Note,
				Status:     item.Status,
				CreatedAt:  item.CreatedAt.Time(),
			}

			if item.Customer != nil {
				res.Customer = &dto.CustomerSub{
					Name:    item.Customer.Name,
					Phone:   item.Customer.Phone,
					Product: item.Customer.Product,
					Branch:  item.Customer.Branch,
				}
			}

			items = append(items, res)
		}
	}

	return items, total, nil
}
