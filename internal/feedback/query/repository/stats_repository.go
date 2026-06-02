package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/feedback/query/dto"
)

func formatFeedbackRating(val float64) float64 {
	return float64(int(val*10)) / 10.0
}

// parsePeriodFilter converts a period string ("7d", "1m") into a bson.M filter on created_at.
// Returns empty bson.M{} when period is empty or unrecognized → zero regression.
func parseFeedbackPeriodFilter(period string) bson.M {
	var since time.Time
	now := time.Now()

	switch period {
	case "7d":
		since = now.AddDate(0, 0, -7)
	case "1m":
		since = now.AddDate(0, -1, 0)
	case "3m":
		since = now.AddDate(0, -3, 0)
	default:
		return bson.M{}
	}

	return bson.M{"created_at": bson.M{"$gte": since}}
}

func (r *feedbackQueryRepository) GetStats(ctx context.Context, branch string, period string) (*dto.FeedbackStatsResponse, error) {
	matchStage := bson.M{}
	if branch != "" {
		matchStage["branch"] = branch
	}

	// Merge period filter into matchStage
	periodFilter := parseFeedbackPeriodFilter(period)
	for k, v := range periodFilter {
		matchStage[k] = v
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

	// Calculate dynamic CSAT trend (weekly or daily based on period)
	var numPoints int
	var stepMs int64
	var trendStart time.Time

	now := time.Now()

	switch period {
	case "7d":
		numPoints = 7
		stepMs = 24 * 60 * 60 * 1000 // 1 day in ms
		trendStart = now.Add(-7 * 24 * time.Hour)
	case "1m":
		numPoints = 4
		stepMs = 7 * 24 * 60 * 60 * 1000 // 7 days in ms
		trendStart = now.Add(-28 * 24 * time.Hour)
	case "3m":
		numPoints = 12
		stepMs = 7 * 24 * 60 * 60 * 1000 // 7 days in ms
		trendStart = now.Add(-12 * 7 * 24 * time.Hour)
	default: // All Time (empty)
		numPoints = 26
		stepMs = 7 * 24 * 60 * 60 * 1000 // 7 days in ms
		trendStart = now.Add(-26 * 7 * 24 * time.Hour) // 182 days (6 months)
	}

	weeklyMatch := bson.M{
		"created_at": bson.M{"$gte": trendStart},
	}
	if branch != "" {
		weeklyMatch["branch"] = branch
	}

	weeklyPipeline := mongo.Pipeline{
		{{Key: "$match", Value: weeklyMatch}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"$floor": bson.M{
					"$divide": bson.A{
						bson.M{"$subtract": bson.A{"$created_at", trendStart}},
						stepMs,
					},
				},
			},
			"avg_rating": bson.M{"$avg": "$rating"},
		}}},
	}

	weeklyCursor, err := r.collection.Aggregate(ctx, weeklyPipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate weekly csat: %w", err)
	}
	defer weeklyCursor.Close(ctx)

	var weeklyResults []struct {
		Index     int     `bson:"_id"`
		AvgRating float64 `bson:"avg_rating"`
	}
	if err := weeklyCursor.All(ctx, &weeklyResults); err != nil {
		return nil, fmt.Errorf("decode weekly csat: %w", err)
	}

	// Initialize CSAT array with overall average rating as a graceful default
	weeklyCSAT := make([]float64, numPoints)
	for i := 0; i < numPoints; i++ {
		weeklyCSAT[i] = formatFeedbackRating(avgRating)
	}

	// Map results to correct index
	for _, res := range weeklyResults {
		if res.Index >= 0 && res.Index < numPoints {
			weeklyCSAT[res.Index] = formatFeedbackRating(res.AvgRating)
		}
	}

	return &dto.FeedbackStatsResponse{
		AvgRating:     avgRating,
		PositiveCount: posCount,
		NeutralCount:  neuCount,
		NegativeCount: negCount,
		WeeklyCSAT:    weeklyCSAT,
	}, nil
}
