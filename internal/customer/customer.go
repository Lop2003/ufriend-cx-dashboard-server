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
	// สร้างดัชนี (Indexes) ทั้งหมดแบบ Background ในรูปแบบ Asynchronous เพื่อเพิ่มสปีดการเสิร์ชระดับ 2 ล้านเรคคอร์ด
	go func() {
		start := time.Now()
		// ให้เวลาเพียงพอสำหรับการสร้าง Index บนข้อมูลขนาดใหญ่ระดับ 10 ล้านเรคคอร์ด
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
		defer cancel()

		customersCol := db.Collection("customers")

		_, err := customersCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
			{Keys: bson.D{{Key: "branch", Value: 1}}},
			{Keys: bson.D{{Key: "status", Value: 1}}},
			{Keys: bson.D{{Key: "created_at", Value: -1}}},
			{Keys: bson.D{{Key: "name", Value: 1}}},
			{Keys: bson.D{{Key: "phone", Value: 1}}},
			{Keys: bson.D{{Key: "product", Value: 1}}},
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
	}()

	queryRepo := queryRepository.NewCustomerQueryRepository(db)
	queryUC := queryUsecase.NewCustomerQueryUsecase(queryRepo)
	queryH := queryHandler.NewCustomerQueryHandler(queryUC)

	return &CustomerDomain{
		queryHandler: queryH,
	}
}

func (d *CustomerDomain) RegisterRoutes(router fiber.Router) {
	http.RegisterCustomerHTTPRoutes(router, d.queryHandler)
}