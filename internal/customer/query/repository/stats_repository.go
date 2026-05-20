package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/customer/query/dto"
)

func (r *customerQueryRepository) GetSummary(ctx context.Context) (*dto.SummaryResponse, error) {
	customers := r.db.Collection("customers")

	total, err := customers.CountDocuments(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("count customers: %w", err)
	}

	overdueCount, err := customers.CountDocuments(ctx, bson.M{"status": "overdue"})
	if err != nil {
		return nil, fmt.Errorf("count overdue: %w", err)
	}

	activeCount, err := customers.CountDocuments(ctx, bson.M{"status": "active"})
	if err != nil {
		return nil, fmt.Errorf("count active: %w", err)
	}

	completedCount, err := customers.CountDocuments(ctx, bson.M{"status": "completed"})
	if err != nil {
		return nil, fmt.Errorf("count completed: %w", err)
	}

	// Average rating from feedbacks
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{"_id": nil, "avg_rating": bson.M{"$avg": "$rating"}}}},
	}
	cursor, err := r.db.Collection("feedbacks").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("avg rating: %w", err)
	}
	defer cursor.Close(ctx)

	var avgResult []struct {
		AvgRating float64 `bson:"avg_rating"`
	}
	if err := cursor.All(ctx, &avgResult); err != nil {
		return nil, fmt.Errorf("decode avg rating: %w", err)
	}

	avgRating := 0.0
	if len(avgResult) > 0 {
		avgRating = avgResult[0].AvgRating
	}

	return &dto.SummaryResponse{
		TotalCustomers: int(total),
		AvgRating:      avgRating,
		OverdueCount:   int(overdueCount),
		ActiveCount:    int(activeCount),
		CompletedCount: int(completedCount),
	}, nil
}

func (r *customerQueryRepository) GetByBranch(ctx context.Context) ([]*dto.BranchStat, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.M{
			"from":         "feedbacks",
			"localField":   "_id",
			"foreignField": "customer_id",
			"as":           "feedbacks",
		}}},
		{{Key: "$unwind", Value: bson.M{"path": "$feedbacks", "preserveNullAndEmptyArrays": true}}},
		{{Key: "$group", Value: bson.M{
			"_id":            "$branch",
			"customer_count": bson.M{"$addToSet": "$_id"},
			"avg_rating":     bson.M{"$avg": "$feedbacks.rating"},
			"overdue_list":   bson.M{"$addToSet": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$status", "overdue"}}, "$_id", nil}}},
		}}},
		{{Key: "$project", Value: bson.M{
			"_id":            1,
			"customer_count": bson.M{"$size": "$customer_count"},
			"avg_rating":     bson.M{"$ifNull": bson.A{"$avg_rating", 0}},
			"overdue_count":  bson.M{"$size": bson.M{"$filter": bson.M{"input": "$overdue_list", "cond": bson.M{"$ne": bson.A{"$$this", nil}}}}},
		}}},
		{{Key: "$sort", Value: bson.M{"_id": 1}}},
	}

	cursor, err := r.db.Collection("customers").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate by branch: %w", err)
	}
	defer cursor.Close(ctx)

	var stats []*dto.BranchStat
	if err := cursor.All(ctx, &stats); err != nil {
		return nil, fmt.Errorf("decode branch stats: %w", err)
	}

	return stats, nil
}