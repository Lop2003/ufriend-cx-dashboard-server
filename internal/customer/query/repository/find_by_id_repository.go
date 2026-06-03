package repository

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"ufriend-cx-dashboard-server/internal/customer/query/dto"
)

type customerDetailAggResult struct {
	Id         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	Phone      string             `bson:"phone"`
	Product    string             `bson:"product"`
	Branch     string             `bson:"branch"`
	PlanMonths int                `bson:"plan_months"`
	Status     string             `bson:"status"`
	CreatedAt  primitive.DateTime `bson:"created_at"`
	Feedbacks  []struct {
		Id        primitive.ObjectID `bson:"_id"`
		Rating    int                `bson:"rating"`
		Comment   string             `bson:"comment"`
		Category  string             `bson:"category"`
		Sentiment string             `bson:"sentiment"`
		CreatedAt primitive.DateTime `bson:"created_at"`
	} `bson:"feedbacks"`
	FollowUps []struct {
		Id        primitive.ObjectID `bson:"_id"`
		Type      string             `bson:"type"`
		Note      string             `bson:"note"`
		Status    string             `bson:"status"`
		CreatedAt primitive.DateTime `bson:"created_at"`
	} `bson:"follow_ups"`
}

func (r *customerQueryRepository) FindByID(ctx context.Context, id string) (*dto.CustomerDetailResponse, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid customer id: %w", err)
	}

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": objId}}},
		// Pipeline $lookup กับ $sort + $limit เพื่อป้องกัน response ขนาดใหญ่
		// คืนเฉพาะ 50 รายการล่าสุด (ถ้าต้องการดูทั้งหมดให้ใช้ API ของ feedback/follow_up โดยตรง)
		{{Key: "$lookup", Value: bson.M{
			"from": "feedbacks",
			"let":  bson.M{"cid": "$_id"},
			"pipeline": bson.A{
				bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$customer_id", "$$cid"}}}},
				bson.M{"$sort": bson.M{"created_at": -1}},
				bson.M{"$limit": 50},
			},
			"as": "feedbacks",
		}}},
		{{Key: "$lookup", Value: bson.M{
			"from": "follow_ups",
			"let":  bson.M{"cid": "$_id"},
			"pipeline": bson.A{
				bson.M{"$match": bson.M{"$expr": bson.M{"$eq": bson.A{"$customer_id", "$$cid"}}}},
				bson.M{"$sort": bson.M{"created_at": -1}},
				bson.M{"$limit": 50},
			},
			"as": "follow_ups",
		}}},
	}

	cursor, err := r.db.Collection("customers").Aggregate(ctx, pipeline)
	if err != nil {
		return nil, fmt.Errorf("aggregate customer: %w", err)
	}
	defer cursor.Close(ctx)

	var results []customerDetailAggResult
	if err := cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("decode customer: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("customer not found")
	}

	r0 := results[0]
	resp := &dto.CustomerDetailResponse{
		Id:         r0.Id.Hex(),
		Name:       r0.Name,
		Phone:      r0.Phone,
		Product:    r0.Product,
		Branch:     r0.Branch,
		PlanMonths: r0.PlanMonths,
		Status:     r0.Status,
		CreatedAt:  r0.CreatedAt.Time(),
		Feedbacks:  make([]dto.FeedbackSub, 0, len(r0.Feedbacks)),
		FollowUps:  make([]dto.FollowUpSub, 0, len(r0.FollowUps)),
	}

	for _, f := range r0.Feedbacks {
		resp.Feedbacks = append(resp.Feedbacks, dto.FeedbackSub{
			Id: f.Id.Hex(), Rating: f.Rating, Comment: f.Comment,
			Category: f.Category, Sentiment: f.Sentiment, CreatedAt: f.CreatedAt.Time(),
		})
	}
	for _, f := range r0.FollowUps {
		resp.FollowUps = append(resp.FollowUps, dto.FollowUpSub{
			Id: f.Id.Hex(), Type: f.Type, Note: f.Note,
			Status: f.Status, CreatedAt: f.CreatedAt.Time(),
		})
	}

	return resp, nil
}