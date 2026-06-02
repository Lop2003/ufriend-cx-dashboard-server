package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
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

type weightedItem struct {
	name   string
	weight float64
}

var weightedBranches = []weightedItem{
	{"สยาม", 0.25},
	{"ลาดพร้าว", 0.18},
	{"เมกาบางนา", 0.15},
	{"ฟิวเจอร์พาร์ค", 0.12},
	{"รังสิต", 0.10},
	{"บางนา", 0.08},
	{"ปิ่นเกล้า", 0.05},
	{"พระราม 9", 0.04},
	{"วงเวียนใหญ่", 0.03},
}

func selectWeightedBranch() string {
	r := rand.Float64()
	var cumulative float64
	for _, item := range weightedBranches {
		cumulative += item.weight
		if r <= cumulative {
			return item.name
		}
	}
	return "สยาม" // fallback
}

var weightedProducts = []weightedItem{
	{"iPhone 16 Pro Max", 0.20},
	{"iPhone 16", 0.18},
	{"iPhone 15", 0.15},
	{"iPhone 16 Pro", 0.12},
	{"iPhone 15 Pro", 0.10},
	{"iPad Air", 0.07},
	{"MacBook Air", 0.06},
	{"Apple Watch", 0.05},
	{"iPad Pro", 0.03},
	{"MacBook Pro", 0.02},
	{"Apple Watch Ultra", 0.01},
	{"AirPods Max", 0.007},
	{"iPad Mini", 0.003},
}

func selectWeightedProduct() string {
	r := rand.Float64()
	var cumulative float64
	for _, item := range weightedProducts {
		cumulative += item.weight
		if r <= cumulative {
			return item.name
		}
	}
	return "iPhone 16 Pro Max" // fallback
}

type weightedPlan struct {
	months int
	weight float64
}

var weightedPlans = []weightedPlan{
	{12, 0.40},
	{24, 0.25},
	{6, 0.15},
	{10, 0.10},
	{20, 0.06},
	{8, 0.04},
}

func selectWeightedPlan() int {
	r := rand.Float64()
	var cumulative float64
	for _, item := range weightedPlans {
		cumulative += item.weight
		if r <= cumulative {
			return item.months
		}
	}
	return 12 // fallback
}


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
	part1 := rand.Intn(900) + 100   // 100-999
	part2 := rand.Intn(9000) + 1000 // 1000-9999
	return fmt.Sprintf("%s-%d-%d", prefix, part1, part2)
}

func main() {
	log.Println("⚡ Starting database seeder...")

	// 1. Load env variables using godotenv
	_ = godotenv.Load()

	// 2. Safeguard check: APP_ENV MUST be 'development' or 'dev'
	appEnv := os.Getenv("APP_ENV")
	if appEnv != "development" && appEnv != "dev" {
		log.Fatalf("❌ CRITICAL ERROR: Seeder can ONLY be run in 'development' or 'dev' environment. Current APP_ENV: '%s'", appEnv)
	}

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	// 3. Double safeguard check: inspect MONGODB_URI to ensure it points to local environment
	isLocalhost := false
	localHosts := []string{"localhost", "127.0.0.1", "::1", "mongodb"}
	for _, host := range localHosts {
		if strings.Contains(mongoURI, host) {
			isLocalhost = true
			break
		}
	}
	if !isLocalhost {
		log.Fatalf("❌ CRITICAL ERROR: MONGODB_URI '%s' does not appear to be a local database. Seeder aborted to prevent data loss.", mongoURI)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal("Connection error:", err)
	}
	defer client.Disconnect(ctx)
	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		dbName = "ufriend_cx"
	}

	db := client.Database(dbName)

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
		product := selectWeightedProduct()
		branch := selectWeightedBranch()
		plan := selectWeightedPlan()

		// Distribute status based on branch weights for highly varied performance indicators
		statusRand := rand.Float64()
		status := "active"
		switch branch {
		case "สยาม", "เมกาบางนา":
			// Premium branches: high active, low overdue
			if statusRand < 0.04 {
				status = "overdue"
			} else if statusRand < 0.20 {
				status = "completed"
			}
		case "บางนา", "วงเวียนใหญ่":
			// Underperforming branches: higher overdue
			if statusRand < 0.30 {
				status = "overdue"
			} else if statusRand < 0.45 {
				status = "completed"
			}
		case "ลาดพร้าว", "พระราม 9":
			// Medium branches
			if statusRand < 0.20 {
				status = "overdue"
			} else if statusRand < 0.35 {
				status = "completed"
			}
		default:
			// Default distribution
			if statusRand < 0.15 {
				status = "overdue"
			} else if statusRand < 0.30 {
				status = "completed"
			}
		}

		// 1. Growth trend over the last 365 days (smooth growing trend using exponent 1.3 to avoid extreme spikes)
		rDays := rand.Float64()
		daysAgo := int(365.0 * math.Pow(rDays, 1.3))

		candidateDate := time.Now().AddDate(0, 0, -daysAgo)

		// 2. Weekly seasonality: Concentrate registrations on weekends (Fri, Sat, Sun)
		// Shift backward to Sunday, Saturday, or Friday to prevent future dates that clump on Today
		wd := candidateDate.Weekday()
		if wd >= time.Monday && wd <= time.Thursday && rand.Float64() < 0.35 {
			daysToSubtract := int(wd) + rand.Intn(3) // Monday (1) -> subtracts 1, 2, or 3 days -> Sun, Sat, Fri
			candidateDate = candidateDate.AddDate(0, 0, -daysToSubtract)
		}

		// Ensure no future dates (safeguard)
		if candidateDate.After(time.Now()) {
			candidateDate = time.Now()
		}
		createdAt := candidateDate

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

			// Category bias based on rating and branch for robust, insightful charts
			var category string
			catRand := rand.Float64()
			if rating <= 2 {
				// Complaints bias
				if branch == "บางนา" || branch == "วงเวียนใหญ่" {
					if catRand < 0.50 {
						category = "payment"
					} else if catRand < 0.80 {
						category = "service"
					} else if catRand < 0.90 {
						category = "branch"
					} else {
						category = "product"
					}
				} else if branch == "ลาดพร้าว" {
					// Service issues (queuing)
					if catRand < 0.60 {
						category = "service"
					} else if catRand < 0.80 {
						category = "branch"
					} else if catRand < 0.95 {
						category = "payment"
					} else {
						category = "product"
					}
				} else {
					if catRand < 0.40 {
						category = "service"
					} else if catRand < 0.75 {
						category = "payment"
					} else if catRand < 0.90 {
						category = "branch"
					} else {
						category = "product"
					}
				}
			} else {
				// Positive feedback bias
				if catRand < 0.45 {
					category = "product"
				} else if catRand < 0.80 {
					category = "service"
				} else if catRand < 0.90 {
					category = "branch"
				} else {
					category = "payment"
				}
			}

			sentiment := sentimentFromRating(rating)
			feedbackDate = createdAt.Add(time.Duration(rand.Intn(24*7)) * time.Hour) // Feedback within a week
			if feedbackDate.After(time.Now()) {
				feedbackDate = time.Now()
			}

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
			if followUpDate.After(time.Now()) {
				followUpDate = time.Now()
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
		{Keys: bson.D{{Key: "name", Value: 1}}},
		{Keys: bson.D{{Key: "phone", Value: 1}}},
		{Keys: bson.D{{Key: "product", Value: 1}}},
	})

	// Index on feedbacks
	_, _ = feedbacksCol.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "branch", Value: 1}}},
		{Keys: bson.D{{Key: "category", Value: 1}}},
		{Keys: bson.D{{Key: "rating", Value: 1}}},
		{Keys: bson.D{{Key: "customer_id", Value: 1}}},
		{Keys: bson.D{{Key: "created_at", Value: 1}}},                            // NEW index for CSAT aggregation
		{Keys: bson.D{{Key: "branch", Value: 1}, {Key: "created_at", Value: 1}}}, // NEW compound index for filtered CSAT aggregation
	})

	elapsed := time.Since(startTime)
	log.Println("──────────────────────────────────────────────────")
	log.Printf("🎉 Database seeding completed successfully in %s!", elapsed)
	log.Printf("👥 Total Customers Seeded: %d", createdCustomersCount)
	log.Printf("💬 Total Feedbacks Seeded: %d", createdFeedbacksCount)
	log.Printf("📞 Total Follow-ups Seeded: %d", createdFollowUpsCount)
	log.Println("──────────────────────────────────────────────────")
}
