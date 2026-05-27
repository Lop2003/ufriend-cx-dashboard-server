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
	// Step 1: Aggregate customers by branch (extremely fast, no join)
	pipelineCustomers := mongo.Pipeline{
		{{Key: "$group", Value: bson.M{
			"_id":            "$branch",
			"customer_count": bson.M{"$sum": 1},
			"overdue_count":  bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$status", "overdue"}}, 1, 0}}},
		}}},
		{{Key: "$sort", Value: bson.M{"_id": 1}}},
	}

	cursorCust, err := r.db.Collection("customers").Aggregate(ctx, pipelineCustomers)
	if err != nil {
		return nil, fmt.Errorf("aggregate customers by branch: %w", err)
	}
	defer cursorCust.Close(ctx)

	var custStats []struct {
		Branch        string `bson:"_id"`
		CustomerCount int    `bson:"customer_count"`
		OverdueCount  int    `bson:"overdue_count"`
	}
	if err := cursorCust.All(ctx, &custStats); err != nil {
		return nil, fmt.Errorf("decode customer branch stats: %w", err)
	}

	// Step 2: Aggregate average feedback ratings by branch (extremely fast because it starts from feedbacks collection)
	pipelineFeedbacks := mongo.Pipeline{
		{{Key: "$lookup", Value: bson.M{
			"from":         "customers",
			"localField":   "customer_id",
			"foreignField": "_id",
			"as":           "customer",
		}}},
		{{Key: "$unwind", Value: "$customer"}},
		{{Key: "$group", Value: bson.M{
			"_id":        "$customer.branch",
			"avg_rating": bson.M{"$avg": "$rating"},
		}}},
	}

	cursorFeed, err := r.db.Collection("feedbacks").Aggregate(ctx, pipelineFeedbacks)
	if err != nil {
		return nil, fmt.Errorf("aggregate feedbacks by branch: %w", err)
	}
	defer cursorFeed.Close(ctx)

	var feedStats []struct {
		Branch    string  `bson:"_id"`
		AvgRating float64 `bson:"avg_rating"`
	}
	if err := cursorFeed.All(ctx, &feedStats); err != nil {
		return nil, fmt.Errorf("decode feedback branch stats: %w", err)
	}

	// Step 3: Merge stats in memory (very small number of branches, N <= 10)
	ratingsMap := make(map[string]float64)
	for _, f := range feedStats {
		ratingsMap[f.Branch] = f.AvgRating
	}

	var stats []*dto.BranchStat
	for _, c := range custStats {
		avgRating := ratingsMap[c.Branch]
		stats = append(stats, &dto.BranchStat{
			Branch:        c.Branch,
			CustomerCount: c.CustomerCount,
			AvgRating:     avgRating,
			OverdueCount:  c.OverdueCount,
		})
	}

	return stats, nil
}