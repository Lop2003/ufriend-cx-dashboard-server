package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/feedback/model"
	"ufriend-cx-dashboard-server/internal/feedback/query/dto"
)

func (r *feedbackQueryRepository) FindAll(ctx context.Context, filter *dto.FeedbackFilter) ([]*model.Feedback, error) {
	// ถ้ามี branch filter → ใช้ aggregation pipeline เพื่อ lookup customers
	if filter.Branch != "" {
		return r.findAllWithBranch(ctx, filter)
	}

	bsonFilter := bson.M{}
	if filter.Category != "" {
		bsonFilter["category"] = filter.Category
	}
	if filter.Rating > 0 {
		bsonFilter["rating"] = filter.Rating
	}

	cursor, err := r.collection.Find(ctx, bsonFilter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var feedbacks []*model.Feedback
	if err := cursor.All(ctx, &feedbacks); err != nil {
		return nil, err
	}

	return feedbacks, nil
}

// findAllWithBranch — ใช้ $lookup join กับ customers เพื่อ filter ตาม branch
func (r *feedbackQueryRepository) findAllWithBranch(ctx context.Context, filter *dto.FeedbackFilter) ([]*model.Feedback, error) {
	pipeline := mongo.Pipeline{
		// Lookup customer เพื่อดึง branch
		{{Key: "$lookup", Value: bson.M{
			"from":         "customers",
			"localField":   "customer_id",
			"foreignField": "_id",
			"as":           "customer",
		}}},
		// Unwind customer array (1:1 relationship)
		{{Key: "$unwind", Value: bson.M{
			"path":                       "$customer",
			"preserveNullAndEmptyArrays": false,
		}}},
		// Match branch
		{{Key: "$match", Value: bson.M{"customer.branch": filter.Branch}}},
	}

	// เพิ่ม filter อื่นๆ ถ้ามี
	matchFilter := bson.M{}
	if filter.Category != "" {
		matchFilter["category"] = filter.Category
	}
	if filter.Rating > 0 {
		matchFilter["rating"] = filter.Rating
	}
	if len(matchFilter) > 0 {
		pipeline = append(pipeline, bson.D{{Key: "$match", Value: matchFilter}})
	}

	// Project กลับเฉพาะ fields ของ feedback (ลบ customer ออก)
	pipeline = append(pipeline, bson.D{{Key: "$project", Value: bson.M{
		"customer": 0,
	}}})

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var feedbacks []*model.Feedback
	if err := cursor.All(ctx, &feedbacks); err != nil {
		return nil, err
	}

	return feedbacks, nil
}