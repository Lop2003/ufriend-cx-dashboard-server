package customer

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	queryHandler "ufriend-cx-dashboard-server/internal/customer/query/handler"
	queryRepository "ufriend-cx-dashboard-server/internal/customer/query/repository"
	queryUsecase "ufriend-cx-dashboard-server/internal/customer/query/usecase"
	"ufriend-cx-dashboard-server/internal/customer/router/http"
)

type CustomerDomain struct {
	queryHandler *queryHandler.CustomerQueryHandler
}

func NewCustomerDomain(db *mongo.Database) *CustomerDomain {
	// สร้าง index แบบ blocking — ต้องสร้างเสร็จก่อนรับ traffic
	// ถ้า index มีอยู่แล้ว MongoDB จะ skip ทันที (idempotent, < 1ms)
	// ใช้ timeout 2 นาทีเพื่อรองรับ initial index creation บนข้อมูลขนาดใหญ่
	createCustomerIndexes(db)

	queryRepo := queryRepository.NewCustomerQueryRepository(db)
	queryUC := queryUsecase.NewCustomerQueryUsecase(queryRepo)
	queryH := queryHandler.NewCustomerQueryHandler(queryUC)

	return &CustomerDomain{
		queryHandler: queryH,
	}
}

// createCustomerIndexes สร้าง index สำหรับ customers collection
// Blocking call — ถ้า fail จะ log error แต่ไม่ crash server (graceful degradation)
func createCustomerIndexes(db *mongo.Database) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	customersCol := db.Collection("customers")

	_, err := customersCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "branch", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
		{Keys: bson.D{{Key: "name", Value: 1}}},
		{Keys: bson.D{{Key: "phone", Value: 1}}},
		{Keys: bson.D{{Key: "product", Value: 1}}},
		// Text index on name for $text word-based search (supports last-name matching on 10M+ records)
		{Keys: bson.D{{Key: "name", Value: "text"}}},
	})

	if err != nil {
		slog.Error("failed to create customer indexes",
			"collection", "customers",
			"duration", time.Since(start).String(),
			"error", err,
		)
	} else {
		slog.Info("customer indexes verified",
			"collection", "customers",
			"indexes", []string{"branch", "status", "created_at", "name", "phone", "product"},
			"duration", time.Since(start).String(),
		)
	}

	// Index สำหรับ $lookup join — ไม่มีจะ scan ทั้ง collection
	createRelatedIndexes(db)
}

// createRelatedIndexes สร้าง index บน feedbacks/follow_ups สำหรับ $lookup
func createRelatedIndexes(db *mongo.Database) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// feedbacks.customer_id — ใช้ใน $lookup จาก customer detail
	if _, err := db.Collection("feedbacks").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "customer_id", Value: 1}},
	}); err != nil {
		slog.Error("failed to create feedbacks.customer_id index", "error", err)
	}

	// follow_ups.customer_id — ใช้ใน $lookup จาก customer detail
	if _, err := db.Collection("follow_ups").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "customer_id", Value: 1}},
	}); err != nil {
		slog.Error("failed to create follow_ups.customer_id index", "error", err)
	}

	slog.Info("related indexes verified",
		"indexes", []string{"feedbacks.customer_id", "follow_ups.customer_id"},
	)
}

func (d *CustomerDomain) RegisterRoutes(router fiber.Router) {
	http.RegisterCustomerHTTPRoutes(router, d.queryHandler)
}