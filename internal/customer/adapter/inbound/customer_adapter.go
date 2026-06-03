package inbound

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// CustomerAdapter — expose ข้อมูล customer ให้ domain อื่นเรียกใช้
// ตาม adapter pattern: domain อื่นไม่ import customer internal โดยตรง
// ถ้าวันนึงแยกเป็น microservice → เปลี่ยน implementation เป็น HTTP/gRPC client
type CustomerAdapter interface {
	// GetBranch คืน branch ของ customer — ใช้ denormalize ลง feedback
	GetBranch(ctx context.Context, customerID primitive.ObjectID) (string, error)
	// Exists เช็คว่า customer มีอยู่จริง — ใช้ validate ก่อนสร้าง follow_up
	Exists(ctx context.Context, customerID primitive.ObjectID) (bool, error)
}

// mongoCustomerAdapter — implementation สำหรับ monolith (อ่านจาก MongoDB ตรง)
type mongoCustomerAdapter struct {
	collection *mongo.Collection
}

func NewCustomerAdapter(db *mongo.Database) CustomerAdapter {
	return &mongoCustomerAdapter{collection: db.Collection("customers")}
}

func (a *mongoCustomerAdapter) GetBranch(ctx context.Context, customerID primitive.ObjectID) (string, error) {
	var result struct {
		Branch string `bson:"branch"`
	}
	err := a.collection.FindOne(ctx, bson.M{"_id": customerID}).Decode(&result)
	if err != nil {
		return "", err
	}
	return result.Branch, nil
}

func (a *mongoCustomerAdapter) Exists(ctx context.Context, customerID primitive.ObjectID) (bool, error) {
	count, err := a.collection.CountDocuments(ctx, bson.M{"_id": customerID})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
