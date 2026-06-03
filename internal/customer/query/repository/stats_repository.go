package repository

import (
	"context"
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/customer/query/dto"
	"ufriend-cx-dashboard-server/pkg/period"
)

// parsePeriodFilter delegates to shared period.ParseFilter
func parsePeriodFilter(p string) bson.M {
	return period.ParseFilter(p)
}

func (r *customerQueryRepository) GetSummary(ctx context.Context, period string) (*dto.SummaryResponse, error) {
	customers := r.db.Collection("customers")

	// Build pipeline: optional period filter + group by status
	periodMatch := parsePeriodFilter(period)

	pipeline := mongo.Pipeline{}
	if len(periodMatch) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: periodMatch}})
	}
	pipeline = append(pipeline, bson.D{{Key: "$group", Value: bson.M{
		"_id":   "$status",
		"count": bson.M{"$sum": 1},
	}}})

	// ใช้ Aggregate เพียง 1 ครั้งเพื่อคำนวณสถิติลูกค้าทั้งหมด ลดภาระฐานข้อมูลจากล้านๆ เรคคอร์ดอย่างมหาศาล
	cursorStats, err := customers.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate customer status stats: %w", err)
	}
	defer cursorStats.Close(ctx)

	var results []struct {
		Status string `bson:"_id"`
		Count  int    `bson:"count"`
	}
	if err := cursorStats.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("decode status stats: %w", err)
	}

	var total, overdueCount, activeCount, completedCount int64
	for _, res := range results {
		total += int64(res.Count)
		switch res.Status {
		case "overdue":
			overdueCount = int64(res.Count)
		case "active":
			activeCount = int64(res.Count)
		case "completed":
			completedCount = int64(res.Count)
		}
	}

	// Average rating from feedbacks (with same period filter)
	feedbackPipeline := mongo.Pipeline{}
	if len(periodMatch) > 0 {
		feedbackPipeline = append(feedbackPipeline, bson.D{{Key: "$match", Value: periodMatch}})
	}
	feedbackPipeline = append(feedbackPipeline, bson.D{{Key: "$group", Value: bson.M{"_id": nil, "avg_rating": bson.M{"$avg": "$rating"}}}})

	cursor, err := r.db.Collection("feedbacks").Aggregate(ctx, feedbackPipeline)
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

func (r *customerQueryRepository) GetByBranch(ctx context.Context, period string, branch string) ([]*dto.BranchStat, error) {
	periodMatch := parsePeriodFilter(period)

	// สร้าง match filter รวม period + branch (ถ้ามี)
	matchFilter := bson.M{}
	for k, v := range periodMatch {
		matchFilter[k] = v
	}
	if branch != "" {
		matchFilter["branch"] = branch
	}

	// Step 1: Aggregate customers by branch
	pipelineCustomers := mongo.Pipeline{}
	if len(matchFilter) > 0 {
		pipelineCustomers = append(pipelineCustomers, bson.D{{Key: "$match", Value: matchFilter}})
	}
	pipelineCustomers = append(pipelineCustomers,
		bson.D{{Key: "$group", Value: bson.M{
			"_id":            "$branch",
			"customer_count": bson.M{"$sum": 1},
			"overdue_count":  bson.M{"$sum": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$status", "overdue"}}, 1, 0}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.M{"_id": 1}}},
	)

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

	// Step 2: Aggregate average feedback ratings by branch (with same combined filter)
	pipelineFeedbacks := mongo.Pipeline{}
	if len(matchFilter) > 0 {
		pipelineFeedbacks = append(pipelineFeedbacks, bson.D{{Key: "$match", Value: matchFilter}})
	}
	pipelineFeedbacks = append(pipelineFeedbacks, bson.D{{Key: "$group", Value: bson.M{
		"_id":        "$branch",
		"avg_rating": bson.M{"$avg": "$rating"},
	}}})

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

func (r *customerQueryRepository) GetDailyStats(ctx context.Context, period string) ([]dto.DailyStatEntry, error) {
	periodMatch := parsePeriodFilter(period)

	pipeline := mongo.Pipeline{}
	if len(periodMatch) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: periodMatch}})
	}

	pipeline = append(pipeline,
		bson.D{{Key: "$project", Value: bson.M{
			"branch": 1,
			"date": bson.M{
				"$dateToString": bson.M{
					"format": "%Y-%m-%d",
					"date":   "$created_at",
				},
			},
		}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"date":   "$date",
				"branch": "$branch",
			},
			"count": bson.M{"$sum": 1},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "_id.date", Value: 1},
		}}},
	)

	cursor, err := r.db.Collection("customers").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate daily stats: %w", err)
	}
	defer cursor.Close(ctx)

	var dbResults []struct {
		Id struct {
			Date   string `bson:"date"`
			Branch string `bson:"branch"`
		} `bson:"_id"`
		Count int `bson:"count"`
	}
	if err := cursor.All(ctx, &dbResults); err != nil {
		return nil, fmt.Errorf("decode daily stats: %w", err)
	}

	// แปลงผลลัพธ์เป็น flat map format ที่ frontend คาดหวัง
	// format: { "date": "2026-01-01", "สยาม": 5, "ลาดพร้าว": 3 }
	rowMap := make(map[string]dto.DailyStatEntry)
	var dateKeys []string

	for _, res := range dbResults {
		date := res.Id.Date
		if date == "" {
			continue
		}

		entry, exists := rowMap[date]
		if !exists {
			entry = dto.DailyStatEntry{"date": date}
			rowMap[date] = entry
			dateKeys = append(dateKeys, date)
		}

		safeKey := strings.ReplaceAll(res.Id.Branch, " ", "_")
		entry[safeKey] = res.Count
	}

	result := make([]dto.DailyStatEntry, 0, len(dateKeys))
	for _, dateKey := range dateKeys {
		result = append(result, rowMap[dateKey])
	}

	return result, nil
}