package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Customer struct {
	Id         primitive.ObjectID `bson:"_id"`
	Name       string             `bson:"name"`
	Phone      string             `bson:"phone"`
	Product    string             `bson:"product"`
	Branch     string             `bson:"branch"`
	PlanMonths int                `bson:"plan_months"`
	Status     string             `bson:"status"`
	CreatedAt  time.Time          `bson:"created_at"`
}

type Feedback struct {
	Id         primitive.ObjectID `bson:"_id"`
	CustomerId primitive.ObjectID `bson:"customer_id"`
	Rating     int                `bson:"rating"`
	Comment    string             `bson:"comment"`
	Category   string             `bson:"category"`
	Sentiment  string             `bson:"sentiment"`
	CreatedAt  time.Time          `bson:"created_at"`
}

type FollowUp struct {
	Id         primitive.ObjectID `bson:"_id"`
	CustomerId primitive.ObjectID `bson:"customer_id"`
	Type       string             `bson:"type"`
	Note       string             `bson:"note"`
	Status     string             `bson:"status"`
	CreatedAt  time.Time          `bson:"created_at"`
}

func sentimentFromRating(rating int) string {
	switch {
	case rating >= 4:
		return "positive"
	case rating == 3:
		return "neutral"
	default:
		return "negative"
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("connect error:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("ufriend_cx")

	// Drop collections เพื่อ seed ใหม่
	db.Collection("customers").Drop(ctx)
	db.Collection("feedbacks").Drop(ctx)
	db.Collection("follow_ups").Drop(ctx)

	// Customers (15 รายการ)
	c1 := primitive.NewObjectID()
	c2 := primitive.NewObjectID()
	c3 := primitive.NewObjectID()
	c4 := primitive.NewObjectID()
	c5 := primitive.NewObjectID()
	c6 := primitive.NewObjectID()
	c7 := primitive.NewObjectID()
	c8 := primitive.NewObjectID()
	c9 := primitive.NewObjectID()
	c10 := primitive.NewObjectID()
	c11 := primitive.NewObjectID()
	c12 := primitive.NewObjectID()
	c13 := primitive.NewObjectID()
	c14 := primitive.NewObjectID()
	c15 := primitive.NewObjectID()

	customers := []interface{}{
		Customer{Id: c1, Name: "สมชาย วงศ์ดี", Phone: "081-111-1111", Product: "iPhone 16 Pro Max", Branch: "วงเวียนใหญ่", PlanMonths: 10, Status: "active", CreatedAt: time.Now().AddDate(0, -2, 0)},
		Customer{Id: c2, Name: "วิภา ศรีสุข", Phone: "082-222-2222", Product: "iPad Air", Branch: "รังสิต", PlanMonths: 6, Status: "overdue", CreatedAt: time.Now().AddDate(0, -3, 0)},
		Customer{Id: c3, Name: "นภา ใจดี", Phone: "083-333-3333", Product: "iPhone 15", Branch: "วงเวียนใหญ่", PlanMonths: 12, Status: "completed", CreatedAt: time.Now().AddDate(0, -6, 0)},
		Customer{Id: c4, Name: "ธนา รักไทย", Phone: "084-444-4444", Product: "iPhone 16", Branch: "ลาดพร้าว", PlanMonths: 10, Status: "active", CreatedAt: time.Now().AddDate(0, -1, 0)},
		Customer{Id: c5, Name: "มานี สดใส", Phone: "085-555-5555", Product: "iPad Pro", Branch: "รังสิต", PlanMonths: 12, Status: "overdue", CreatedAt: time.Now().AddDate(0, -4, 0)},
		Customer{Id: c6, Name: "สุนีย์ ศรีวงษ์", Phone: "086-666-6666", Product: "MacBook Pro", Branch: "วงเวียนใหญ่", PlanMonths: 24, Status: "active", CreatedAt: time.Now().AddDate(0, -1, -10)},
		Customer{Id: c7, Name: "อรุณ สุขสวัสดิ์", Phone: "087-777-7777", Product: "Apple Watch", Branch: "ลาดพร้าว", PlanMonths: 8, Status: "active", CreatedAt: time.Now().AddDate(0, 0, -5)},
		Customer{Id: c8, Name: "จารุณี อุดมศรี", Phone: "088-888-8888", Product: "iPhone 15 Pro", Branch: "สยาม", PlanMonths: 12, Status: "overdue", CreatedAt: time.Now().AddDate(0, -5, 0)},
		Customer{Id: c9, Name: "ชัชชาติ บุญเรือง", Phone: "089-999-9999", Product: "iPad Mini", Branch: "รังสิต", PlanMonths: 10, Status: "completed", CreatedAt: time.Now().AddDate(0, -7, 0)},
		Customer{Id: c10, Name: "นัยนา เพ็ญศรี", Phone: "081-111-2222", Product: "AirPods Max", Branch: "วงเวียนใหญ่", PlanMonths: 6, Status: "active", CreatedAt: time.Now().AddDate(0, -2, -15)},
		Customer{Id: c11, Name: "ปราชญา ตรีศิริ", Phone: "082-222-3333", Product: "iPhone 16 Pro", Branch: "ลาดพร้าว", PlanMonths: 12, Status: "overdue", CreatedAt: time.Now().AddDate(0, -2, 0)},
		Customer{Id: c12, Name: "รพีพร จินดานนท์", Phone: "083-333-4444", Product: "MacBook Air", Branch: "สยาม", PlanMonths: 20, Status: "active", CreatedAt: time.Now().AddDate(0, -1, 0)},
		Customer{Id: c13, Name: "สิทธิศักดิ์ พลอยแจ่ม", Phone: "084-444-5555", Product: "iPad Air", Branch: "รังสิต", PlanMonths: 8, Status: "completed", CreatedAt: time.Now().AddDate(0, -8, 0)},
		Customer{Id: c14, Name: "ทิพย์ สวามิตร", Phone: "085-555-6666", Product: "iPhone 15", Branch: "วงเวียนใหญ่", PlanMonths: 10, Status: "active", CreatedAt: time.Now().AddDate(0, -1, -20)},
		Customer{Id: c15, Name: "วัฒนา ธารสิงห์", Phone: "086-666-7777", Product: "Apple Watch Ultra", Branch: "ลาดพร้าว", PlanMonths: 6, Status: "overdue", CreatedAt: time.Now().AddDate(0, -4, -10)},
	}

	_, err = db.Collection("customers").InsertMany(ctx, customers)
	if err != nil {
		log.Fatal("seed customers error:", err)
	}
	log.Println("✓ customers seeded")

	// Feedbacks (15 รายการ)
	feedbacks := []interface{}{
		Feedback{Id: primitive.NewObjectID(), CustomerId: c1, Rating: 5, Comment: "พนักงานพูดดี อธิบายชัดเจน", Category: "service", Sentiment: sentimentFromRating(5), CreatedAt: time.Now().AddDate(0, -1, 0)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c2, Rating: 2, Comment: "ดอกเบี้ยแพงไป ไม่เข้าใจตอนทำสัญญา ผ่อนแพง", Category: "payment", Sentiment: sentimentFromRating(2), CreatedAt: time.Now().AddDate(0, -2, 0)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c3, Rating: 4, Comment: "สินค้าดี ผ่อนครบแล้ว พอใจทั้งหมด", Category: "product", Sentiment: sentimentFromRating(4), CreatedAt: time.Now().AddDate(0, -1, -15)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c4, Rating: 3, Comment: "โอเค ไม่มีปัญหาอะไร ปกติดี", Category: "service", Sentiment: sentimentFromRating(3), CreatedAt: time.Now().AddDate(0, 0, -10)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c5, Rating: 1, Comment: "ติดต่อยาก ไม่มีคนรับโทรศัพท์ บริการแย่", Category: "branch", Sentiment: sentimentFromRating(1), CreatedAt: time.Now().AddDate(0, -3, 0)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c6, Rating: 5, Comment: "ทีมงานมืออาชีพ สินค้าคุณภาพดี แนะนำได้", Category: "service", Sentiment: sentimentFromRating(5), CreatedAt: time.Now().AddDate(0, -1, -10)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c7, Rating: 4, Comment: "กระบวนการผ่อนง่าย เหมาะสมมาก", Category: "product", Sentiment: sentimentFromRating(4), CreatedAt: time.Now().AddDate(0, 0, -5)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c8, Rating: 2, Comment: "ค้างชำระนาน เอกสารสับสน", Category: "payment", Sentiment: sentimentFromRating(2), CreatedAt: time.Now().AddDate(0, -2, -5)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c9, Rating: 5, Comment: "บริการท่องหาหมื่นครั้ง พอใจมากครับ", Category: "service", Sentiment: sentimentFromRating(5), CreatedAt: time.Now().AddDate(0, -1, 0)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c10, Rating: 3, Comment: "พอใจปานกลาง ทำงานปกติ", Category: "service", Sentiment: sentimentFromRating(3), CreatedAt: time.Now().AddDate(0, 0, -7)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c11, Rating: 1, Comment: "อัตราดอกเบี้ยสูง ไม่เหมาะสม ค้างชำระ", Category: "payment", Sentiment: sentimentFromRating(1), CreatedAt: time.Now().AddDate(0, -1, -10)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c12, Rating: 4, Comment: "ได้สินค้าตามที่ต้องการ บริการดี", Category: "product", Sentiment: sentimentFromRating(4), CreatedAt: time.Now().AddDate(0, 0, -3)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c13, Rating: 5, Comment: "โครงการผ่อนดี สมัครสำเร็จ พอใจมาก", Category: "service", Sentiment: sentimentFromRating(5), CreatedAt: time.Now().AddDate(0, -1, -15)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c14, Rating: 4, Comment: "สาขาเข้าใจดี ใช้บริการได้สะดวก", Category: "branch", Sentiment: sentimentFromRating(4), CreatedAt: time.Now().AddDate(0, 0, -8)},
		Feedback{Id: primitive.NewObjectID(), CustomerId: c15, Rating: 2, Comment: "ผ่อนติดขัด ติดต่อบอก ไม่มีหนทาง", Category: "payment", Sentiment: sentimentFromRating(2), CreatedAt: time.Now().AddDate(0, -2, -3)},
	}

	_, err = db.Collection("feedbacks").InsertMany(ctx, feedbacks)
	if err != nil {
		log.Fatal("seed feedbacks error:", err)
	}
	log.Println("✓ feedbacks seeded")

	// Follow ups
	followUps := []interface{}{
		FollowUp{Id: primitive.NewObjectID(), CustomerId: c2, Type: "payment_remind", Note: "โทรแจ้งค้างชำระ 2 งวด รอการติดต่อกลับ", Status: "pending", CreatedAt: time.Now().AddDate(0, 0, -5)},
		FollowUp{Id: primitive.NewObjectID(), CustomerId: c5, Type: "payment_remind", Note: "ส่ง SMS แจ้งเตือนค้างชำระแล้ว", Status: "done", CreatedAt: time.Now().AddDate(0, -1, 0)},
		FollowUp{Id: primitive.NewObjectID(), CustomerId: c2, Type: "feedback_reply", Note: "ขอโทษลูกค้าเรื่องดอกเบี้ย อธิบายเงื่อนไขเพิ่มเติม", Status: "done", CreatedAt: time.Now().AddDate(0, -2, 0)},
		FollowUp{Id: primitive.NewObjectID(), CustomerId: c8, Type: "payment_remind", Note: "โทรติดตามการชำระเงิน ปรึกษาแผนชำระ", Status: "pending", CreatedAt: time.Now().AddDate(0, 0, -3)},
		FollowUp{Id: primitive.NewObjectID(), CustomerId: c11, Type: "payment_remind", Note: "ส่ง Email แจ้งเตือนค้างชำระ งวดที่ 3", Status: "done", CreatedAt: time.Now().AddDate(0, -1, -5)},
		FollowUp{Id: primitive.NewObjectID(), CustomerId: c15, Type: "payment_remind", Note: "ติดต่อทำการชำระบางส่วน อยู่ระหว่างต่อรอง", Status: "pending", CreatedAt: time.Now().AddDate(0, 0, -2)},
	}

	_, err = db.Collection("follow_ups").InsertMany(ctx, followUps)
	if err != nil {
		log.Fatal("seed follow_ups error:", err)
	}
	log.Println("✓ follow_ups seeded")

	log.Println("seed เสร็จแล้ว!")
}
