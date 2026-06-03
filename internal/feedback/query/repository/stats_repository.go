package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/feedback/query/dto"
	"ufriend-cx-dashboard-server/pkg/period"
)

func formatFeedbackRating(val float64) float64 {
	return float64(int(val*10)) / 10.0
}

// parseFeedbackPeriodFilter delegates to shared period.ParseFilter
func parseFeedbackPeriodFilter(p string) bson.M {
	return period.ParseFilter(p)
}

// --- Trend Config ---

// trendConfig กำหนดพารามิเตอร์สำหรับ CSAT trend chart ตามช่วงเวลาที่เลือก
type trendConfig struct {
	NumPoints  int       // จำนวนจุดบน chart
	StepMs     int64     // ระยะห่างระหว่างจุด (milliseconds)
	TrendStart time.Time // เวลาเริ่มต้นของ trend
}

// resolveTrendConfig แปลง period string เป็นค่า config สำหรับ CSAT trend pipeline
func resolveTrendConfig(period string) trendConfig {
	now := time.Now()

	switch period {
	case "7d":
		return trendConfig{
			NumPoints:  7,
			StepMs:     24 * 60 * 60 * 1000, // 1 day
			TrendStart: now.Add(-7 * 24 * time.Hour),
		}
	case "1m":
		return trendConfig{
			NumPoints:  4,
			StepMs:     7 * 24 * 60 * 60 * 1000, // 7 days
			TrendStart: now.Add(-28 * 24 * time.Hour),
		}
	case "3m":
		return trendConfig{
			NumPoints:  12,
			StepMs:     7 * 24 * 60 * 60 * 1000, // 7 days
			TrendStart: now.Add(-12 * 7 * 24 * time.Hour),
		}
	default: // All Time
		return trendConfig{
			NumPoints:  26,
			StepMs:     7 * 24 * 60 * 60 * 1000,          // 7 days
			TrendStart: now.Add(-26 * 7 * 24 * time.Hour), // ~6 months
		}
	}
}

// --- Sentiment Stats Pipeline ---

// sentimentResult ผลลัพธ์จาก aggregation สำหรับค่าเฉลี่ยและจำนวน sentiment
type sentimentResult struct {
	AvgRating float64
	Positive  int
	Neutral   int
	Negative  int
}

// querySentimentStats รัน aggregation pipeline สำหรับ avg rating + sentiment counts
func querySentimentStats(ctx context.Context, collection *mongo.Collection, matchStage bson.M) (*sentimentResult, error) {
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

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate feedback stats: %w", err)
	}
	defer cursor.Close(ctx)

	var results []struct {
		AvgRating float64 `bson:"avg_rating"`
		Positive  int     `bson:"positive"`
		Neutral   int     `bson:"neutral"`
		Negative  int     `bson:"negative"`
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("decode feedback stats: %w", err)
	}

	result := &sentimentResult{}
	if len(results) > 0 {
		result.AvgRating = results[0].AvgRating
		result.Positive = results[0].Positive
		result.Neutral = results[0].Neutral
		result.Negative = results[0].Negative
	}
	return result, nil
}

// --- CSAT Trend Pipeline ---

// queryCSATTrend รัน aggregation pipeline สำหรับ CSAT trend chart (weekly/daily based on period)
func queryCSATTrend(ctx context.Context, collection *mongo.Collection, branch string, cfg trendConfig, fallbackAvg float64) ([]float64, error) {
	match := bson.M{
		"created_at": bson.M{"$gte": cfg.TrendStart},
	}
	if branch != "" {
		match["branch"] = branch
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: match}},
		{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"$floor": bson.M{
					"$divide": bson.A{
						bson.M{"$subtract": bson.A{"$created_at", cfg.TrendStart}},
						cfg.StepMs,
					},
				},
			},
			"avg_rating": bson.M{"$avg": "$rating"},
		}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate csat trend: %w", err)
	}
	defer cursor.Close(ctx)

	var trendResults []struct {
		Index     int     `bson:"_id"`
		AvgRating float64 `bson:"avg_rating"`
	}
	if err := cursor.All(ctx, &trendResults); err != nil {
		return nil, fmt.Errorf("decode csat trend: %w", err)
	}

	// Initialize array กับ overall average เป็น default (graceful fallback)
	csatTrend := make([]float64, cfg.NumPoints)
	for i := range csatTrend {
		csatTrend[i] = formatFeedbackRating(fallbackAvg)
	}

	// Map ผลลัพธ์จริงลงใน array ตาม index
	for _, res := range trendResults {
		if res.Index >= 0 && res.Index < cfg.NumPoints {
			csatTrend[res.Index] = formatFeedbackRating(res.AvgRating)
		}
	}

	return csatTrend, nil
}

// --- Orchestrator ---

func (r *feedbackQueryRepository) GetStats(ctx context.Context, branch string, period string) (*dto.FeedbackStatsResponse, error) {
	// 1. สร้าง match filter
	matchStage := bson.M{}
	if branch != "" {
		matchStage["branch"] = branch
	}
	for k, v := range parseFeedbackPeriodFilter(period) {
		matchStage[k] = v
	}

	// 2. Query sentiment stats (avg rating + positive/neutral/negative counts)
	sentiment, err := querySentimentStats(ctx, r.collection, matchStage)
	if err != nil {
		return nil, err
	}

	// 3. Query CSAT trend (weekly/daily chart data)
	cfg := resolveTrendConfig(period)
	csatTrend, err := queryCSATTrend(ctx, r.collection, branch, cfg, sentiment.AvgRating)
	if err != nil {
		return nil, err
	}

	return &dto.FeedbackStatsResponse{
		AvgRating:     sentiment.AvgRating,
		PositiveCount: sentiment.Positive,
		NeutralCount:  sentiment.Neutral,
		NegativeCount: sentiment.Negative,
		WeeklyCSAT:    csatTrend,
	}, nil
}

