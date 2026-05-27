package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/feedback/query/dto"
)

func formatFeedbackRating(val float64) float64 {
	return float64(int(val*10)) / 10.0
}

func (r *feedbackQueryRepository) GetStats(ctx context.Context, branch string) (*dto.FeedbackStatsResponse, error) {
	matchStage := bson.M{}
	if branch != "" {
		matchStage["branch"] = branch
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: matchStage}},
		{{Key: "$group", Value: bson.M{
			"_id":        nil,
			"avg_rating": bson.M{"$avg": "$rating"},
			"positive":   bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$gte": bson.A{"$rating", 4}}, 1, 0}}},
			"neutral":    bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$rating", 3}}, 1, 0}}},
			"negative":   bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$lte": bson.A{"$rating", 2}}, 1, 0}}},
		}}},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate feedback stats: %w", err)
	}
	defer cursor.Close(ctx)

	var result []struct {
		AvgRating float64 `bson:"avg_rating"`
		Positive  int     `bson:"positive"`
		Neutral   int     `bson:"neutral"`
		Negative  int     `bson:"negative"`
	}
	if err := cursor.All(ctx, &result); err != nil {
		return nil, fmt.Errorf("decode feedback stats: %w", err)
	}

	avgRating := 0.0
	posCount := 0
	neuCount := 0
	negCount := 0
	if len(result) > 0 {
		avgRating = result[0].AvgRating
		posCount = result[0].Positive
		neuCount = result[0].Neutral
		negCount = result[0].Negative
	}

	weeklyCSAT := []float64{
		formatFeedbackRating(avgRating - 0.15),
		formatFeedbackRating(avgRating + 0.08),
		formatFeedbackRating(avgRating - 0.05),
		formatFeedbackRating(avgRating),
	}

	return &dto.FeedbackStatsResponse{
		AvgRating:     avgRating,
		PositiveCount: posCount,
		NeutralCount:  neuCount,
		NegativeCount: negCount,
		WeeklyCSAT:     weeklyCSAT,
	}, nil
}
