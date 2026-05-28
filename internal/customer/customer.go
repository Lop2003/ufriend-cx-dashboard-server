package customer

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		customersCol := db.Collection("customers")
		
		// ตั้งค่าการสร้างดัชนีเป็น Background index build (ไม่บล็อกการเขียน/อ่านฐานข้อมูล)
		indexOpts := options.CreateIndexes().SetMaxTime(25 * time.Second)

		_, err := customersCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
			{Keys: bson.D{{Key: "branch", Value: 1}}},
			{Keys: bson.D{{Key: "status", Value: 1}}},
			{Keys: bson.D{{Key: "created_at", Value: -1}}},
			{Keys: bson.D{{Key: "name", Value: 1}}},
			{Keys: bson.D{{Key: "phone", Value: 1}}},
			{Keys: bson.D{{Key: "product", Value: 1}}},
		}, indexOpts)

		if err != nil {
			log.Printf("⚠️ Warning: Failed to create MongoDB search indexes: %v", err)
		} else {
			log.Println("⚡ Successfully verified and created uFriend high-performance search indexes (Name, Phone, Product, Branch, Status)!")
		}
	}()

	queryRepo := queryRepository.NewCustomerQueryRepository(db)
	queryUC := queryUsecase.NewCustomerQueryUsecase(queryRepo)
	queryH := queryHandler.NewCustomerQueryHandler(queryUC)

	return &CustomerDomain{
		queryHandler: queryH,
	}
}

func (d *CustomerDomain) RegisterRoutes(app *fiber.App) {
	http.RegisterCustomerHTTPRoutes(app, d.queryHandler)
}