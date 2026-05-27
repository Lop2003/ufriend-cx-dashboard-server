package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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
	Branch     string             `bson:"branch"`
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

// Slices for random generation
var firstNames = []string{
	"สมชาย", "วิภา", "นภา", "ธนา", "มานี", "สุนีย์", "อรุณ", "จารุณี", "ชัชชาติ", "นัยนา",
	"ปราชญา", "รพีพร", "สิทธิศักดิ์", "ทิพย์", "วัฒนา", "กิตติ", "ธีรวัฒน์", "อภิชาติ", "สุรศักดิ์", "ประเสริฐ",
	"มนัส", "ปรีชา", "สมเกียรติ", "สมพงษ์", "วิทูรย์", "วิชัย", "เกรียงไกร", "เกียรติศักดิ์", "อนันต์", "สมบัติ",
	"ประภาส", "อำนวย", "สุรเดช", "มานะ", "ชูชาติ", "สมคิด", "พิพัฒน์", "วีระ", "อดิศักดิ์", "ดนัย",
	"พัชรา", "ยุพา", "วรรณา", "นงลักษณ์", "ศิริพร", "เพ็ญศรี", "สุดา", "สุวรรณา", "อารีย์", "วิไล",
}

var lastNames = []string{
	"วงศ์ดี", "ศรีสุข", "ใจดี", "รักไทย", "สดใส", "ศรีวงษ์", "สุขสวัสดิ์", "อุดมศรี", "บุญเรือง", "เพ็ญศรี",
	"ตรีศิริ", "จินดานนท์", "พลอยแจ่ม", "สวามิตร", "ธารสิงห์", "เจริญยศ", "รุ่งเรือง", "งามดี", "แสงทอง", "ดีเลิศ",
	"มั่นคง", "ทองดี", "รักษาธรรม", "พึ่งบุญ", "ทรัพย์สิน", "สุนทร", "สมบูรณ์", "ประเสริฐสุข", "บุญเกิด", "เลิศวิจิตร",
	"สุขเจริญ", "พงษ์พานิช", "วัฒนพานิช", "เกียรติขจร", "รักษ์ดี", "สิริวัฒนา", "โสภณ", "เรืองโรจน์", "ศิริวัฒน์",
}

var products = []string{
	"iPhone 16 Pro Max", "iPad Air", "iPhone 15", "iPhone 16", "iPad Pro",
	"MacBook Pro", "Apple Watch", "iPhone 15 Pro", "iPad Mini", "AirPods Max",
	"iPhone 16 Pro", "MacBook Air", "Apple Watch Ultra",
}

var branches = []string{
	"วงเวียนใหญ่", "รังสิต", "ลาดพร้าว", "สยาม", "บางนา", "ปิ่นเกล้า", "พระราม 9", "ฟิวเจอร์พาร์ค", "เมกาบางนา",
}

var planMonths = []int{6, 8, 10, 12, 20, 24}

var positiveComments = []string{
	"บริการดีเยี่ยมมากครับ", "พนักงานพูดจาดี น่ารักมาก", "อนุมัติไวมาก แนะนำเลยครับ",
	"สินค้าคุณภาพดีมาก ไม่มีปัญหาเลย", "สาขาบริการรวดเร็วทันใจ ประทับใจมาก",
	"ประทับใจระบบผ่อน สะดวกมาก", "เจ้าหน้าที่ตอบคำถามชัดเจนดีมากครับ", "ชอบบริการสาขานี้มาก รวดเร็วดี",
}

var neutralComments = []string{
	"บริการอยู่ในระดับปานกลาง", "พอใช้ได้ ไม่มีปัญหาอะไรพิเศษ", "ดอกเบี้ยปานกลาง พอรับได้",
	"ระบบมีหน่วงๆ บ้างบางช่วง แต่ใช้งานได้", "ตามมาตรฐานทั่วไป ไม่มีอะไรพิเศษ",
	"โอเค ไม่มีปัญหาอะไร ปกติดี", "พอใจปานกลาง ทำงานปกติ",
}

var negativeComments = []string{
	"บริการแย่มาก ติดต่อยากสุดๆ", "ดอกเบี้ยแพงเกินไป ผ่อนไม่ไหวแล้ว", "พนักงานไม่มีความสุภาพเลย",
	"ระบบล่มบ่อยมาก ทำรายการไม่ได้เลย", "ค้างชำระเพราะระบบมีปัญหา", "ติดต่อเจ้าหน้าที่ยากมาก",
	"ดอกเบี้ยค่อนข้างแพงไปหน่อย ไม่เข้าใจสัญญา", "พนักงานบริการช้า รอนานมาก",
}

var categories = []string{"service", "payment", "product", "branch"}

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

func randomPhone() string {
	prefixes := []string{"081", "082", "083", "084", "085", "086", "087", "088", "089", "095", "096", "097"}
	prefix := prefixes[rand.Intn(len(prefixes))]
	part1 := rand.Intn(900) + 100  // 100-999
	part2 := rand.Intn(9000) + 1000 // 1000-9999
	return fmt.Sprintf("%s-%d-%d", prefix, part1, part2)
}

func main() {
	log.Println("⚡ Starting database seeder...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("ufriend_cx")

	// Drop collections
	log.Println("🗑️ Dropping existing collections...")
	db.Collection("customers").Drop(ctx)
	db.Collection("feedbacks").Drop(ctx)
	db.Collection("follow_ups").Drop(ctx)

	totalCustomers := 2000000
	batchSize := 10000

	customersCol := db.Collection("customers")
	feedbacksCol := db.Collection("feedbacks")
	followUpsCol := db.Collection("follow_ups")

	rand.Seed(time.Now().UnixNano())

	startTime := time.Now()

	var customersBatch []interface{}
	var feedbacksBatch []interface{}
	var followUpsBatch []interface{}

	var createdCustomersCount int
	var createdFeedbacksCount int
	var createdFollowUpsCount int

	log.Printf("🌱 Generating %d customers, comments, and follow-ups...", totalCustomers)

	for i := 1; i <= totalCustomers; i++ {
		// 1. Generate Customer
		cID := primitive.NewObjectID()
		name := firstNames[rand.Intn(len(firstNames))] + " " + lastNames[rand.Intn(len(lastNames))]
		phone := randomPhone()
		product := products[rand.Intn(len(products))]
		branch := branches[rand.Intn(len(branches))]
		plan := planMonths[rand.Intn(len(planMonths))]

		// Distribute status: 70% active, 15% completed, 15% overdue
		statusRand := rand.Float64()
		status := "active"
		if statusRand < 0.15 {
			status = "completed"
		} else if statusRand < 0.30 {
			status = "overdue"
		}

		// CreatedAt in the last 365 days
		createdAt := time.Now().AddDate(0, 0, -rand.Intn(365))

		customer := Customer{
			Id:         cID,
			Name:       name,
			Phone:      phone,
			Product:    product,
			Branch:     branch,
			PlanMonths: plan,
			Status:     status,
			CreatedAt:  createdAt,
		}
		customersBatch = append(customersBatch, customer)

		// 2. Generate Feedback (approx 30% chance for any customer)
		hasFeedback := rand.Float64() < 0.30
		var feedbackRating int
		var feedbackDate time.Time

		if hasFeedback {
			// Rating bias: completed is positive, overdue is negative, active is random
			var rating int
			if status == "completed" {
				// Bias 4-5 stars
				rating = 4 + rand.Intn(2)
			} else if status == "overdue" {
				// Bias 1-3 stars
				rating = 1 + rand.Intn(3)
			} else {
				// Branch-specific bias
				r := rand.Float64()
				switch branch {
				case "สยาม", "เมกาบางนา":
					// High satisfaction bias (mostly 4-5 stars)
					if r < 0.05 {
						rating = 1
					} else if r < 0.10 {
						rating = 2
					} else if r < 0.20 {
						rating = 3
					} else if r < 0.55 {
						rating = 4
					} else {
						rating = 5
					}
				case "ลาดพร้าว", "บางนา":
					// Low satisfaction bias (mostly 1-3 stars)
					if r < 0.35 {
						rating = 1
					} else if r < 0.60 {
						rating = 2
					} else if r < 0.80 {
						rating = 3
					} else if r < 0.95 {
						rating = 4
					} else {
						rating = 5
					}
				default:
					// Average bias
					if r < 0.10 {
						rating = 1
					} else if r < 0.20 {
						rating = 2
					} else if r < 0.35 {
						rating = 3
					} else if r < 0.65 {
						rating = 4
					} else {
						rating = 5
					}
				}
			}
			feedbackRating = rating

			var comment string
			if rating >= 4 {
				comment = positiveComments[rand.Intn(len(positiveComments))]
			} else if rating == 3 {
				comment = neutralComments[rand.Intn(len(neutralComments))]
			} else {
				comment = negativeComments[rand.Intn(len(negativeComments))]
			}

			category := categories[rand.Intn(len(categories))]
			sentiment := sentimentFromRating(rating)
			feedbackDate = createdAt.Add(time.Duration(rand.Intn(24*7)) * time.Hour) // Feedback within a week

			feedback := Feedback{
				Id:         primitive.NewObjectID(),
				CustomerId: cID,
				Branch:     branch,
				Rating:     rating,
				Comment:    comment,
				Category:   category,
				Sentiment:  sentiment,
				CreatedAt:  feedbackDate,
			}
			feedbacksBatch = append(feedbacksBatch, feedback)
		}

		// 3. Generate Follow Up
		// High probability if overdue (80%), or if they left a bad feedback (rating <= 2) (70%)
		shouldFollowUp := false
		followUpType := "general"
		followUpNote := ""

		if status == "overdue" && rand.Float64() < 0.80 {
			shouldFollowUp = true
			followUpType = "payment_remind"
			notes := []string{
				"โทรแจ้งยอดค้างชำระ ลูกค้าสัญญาว่าจะจ่ายภายในสัปดาห์นี้",
				"ส่ง SMS แจ้งเตือนยอดค้างชำระเรียบร้อย",
				"โทรติดต่อไม่ได้ อยู่ระหว่างรอติดต่อกลับ",
				"ส่งอีเมลแจ้งเตือนการค้างชำระงวดล่าสุด",
				"เสนอปรับโครงสร้างหนี้และลดดอกเบี้ยให้ลูกค้า",
			}
			followUpNote = notes[rand.Intn(len(notes))]
		} else if hasFeedback && feedbackRating <= 2 && rand.Float64() < 0.70 {
			shouldFollowUp = true
			followUpType = "feedback_reply"
			notes := []string{
				"ติดต่อลูกค้าเพื่อสอบถามปัญหารายละเอียดเรื่องร้องเรียนและเสนอแนวทางแก้ไข",
				"โทรขอโทษลูกค้าเรื่องการบริการของพนักงานและประสานงานสาขาตรวจสอบ",
				"ส่งอีเมลชี้แจงการแก้ไขข้อผิดพลาดระบบและมอบโค้ดส่วนลดพิเศษชดเชย",
			}
			followUpNote = notes[rand.Intn(len(notes))]
		}

		if shouldFollowUp {
			followUpStatus := "pending"
			if rand.Float64() < 0.50 {
				followUpStatus = "done"
			}

			var followUpDate time.Time
			if !feedbackDate.IsZero() {
				followUpDate = feedbackDate.Add(time.Duration(rand.Intn(24*3)) * time.Hour) // Followup within 3 days of feedback
			} else {
				followUpDate = createdAt.Add(time.Duration(rand.Intn(24*14)) * time.Hour) // Followup within 2 weeks of creation
			}

			followUp := FollowUp{
				Id:         primitive.NewObjectID(),
				CustomerId: cID,
				Type:       followUpType,
				Note:       followUpNote,
				Status:     followUpStatus,
				CreatedAt:  followUpDate,
			}
			followUpsBatch = append(followUpsBatch, followUp)
		}

		// Insert batches to avoid running out of memory
		if len(customersBatch) >= batchSize {
			_, err = customersCol.InsertMany(ctx, customersBatch)
			if err != nil {
				log.Fatal("Error inserting customers batch:", err)
			}
			createdCustomersCount += len(customersBatch)
			customersBatch = nil

			if len(feedbacksBatch) > 0 {
				_, err = feedbacksCol.InsertMany(ctx, feedbacksBatch)
				if err != nil {
					log.Fatal("Error inserting feedbacks batch:", err)
				}
				createdFeedbacksCount += len(feedbacksBatch)
				feedbacksBatch = nil
			}

			if len(followUpsBatch) > 0 {
				_, err = followUpsCol.InsertMany(ctx, followUpsBatch)
				if err != nil {
					log.Fatal("Error inserting followups batch:", err)
				}
				createdFollowUpsCount += len(followUpsBatch)
				followUpsBatch = nil
			}

			log.Printf("⏳ Seeding progress: %d / %d customers...", createdCustomersCount, totalCustomers)
		}
	}

	// Insert any remaining items in the batches
	if len(customersBatch) > 0 {
		_, err = customersCol.InsertMany(ctx, customersBatch)
		if err != nil {
			log.Fatal("Error inserting final customers batch:", err)
		}
		createdCustomersCount += len(customersBatch)
	}
	if len(feedbacksBatch) > 0 {
		_, err = feedbacksCol.InsertMany(ctx, feedbacksBatch)
		if err != nil {
			log.Fatal("Error inserting final feedbacks batch:", err)
		}
		createdFeedbacksCount += len(feedbacksBatch)
	}
	if len(followUpsBatch) > 0 {
		_, err = followUpsCol.InsertMany(ctx, followUpsBatch)
		if err != nil {
			log.Fatal("Error inserting final followups batch:", err)
		}
		createdFollowUpsCount += len(followUpsBatch)
	}

	log.Println("⚡ Creating indexes for high-performance queries...")
	// Index on customers
	_, _ = customersCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "branch", Value: 1}}},
		{Keys: bson.D{{Key: "status", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: -1}}},
	})
	
	// Index on feedbacks
	_, _ = feedbacksCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "branch", Value: 1}}},
		{Keys: bson.D{{Key: "category", Value: 1}}},
		{Keys: bson.D{{Key: "rating", Value: 1}}},
		{Keys: bson.D{{Key: "customer_id", Value: 1}}},
	})

	elapsed := time.Since(startTime)
	log.Println("──────────────────────────────────────────────────")
	log.Printf("🎉 Database seeding completed successfully in %s!", elapsed)
	log.Printf("👥 Total Customers Seeded: %d", createdCustomersCount)
	log.Printf("💬 Total Feedbacks Seeded: %d", createdFeedbacksCount)
	log.Printf("📞 Total Follow-ups Seeded: %d", createdFollowUpsCount)
	log.Println("──────────────────────────────────────────────────")
}
