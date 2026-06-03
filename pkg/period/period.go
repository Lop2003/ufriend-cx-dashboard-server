package period

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

// ParseFilter converts a period string ("7d", "1m", "3m") into a bson.M filter on created_at.
// Returns empty bson.M{} when period is empty or unrecognized → zero regression.
// ใช้ร่วมกันระหว่าง customer และ feedback repositories เพื่อไม่ duplicate logic
func ParseFilter(period string) bson.M {
	var since time.Time
	now := time.Now()

	switch period {
	case "7d":
		since = now.AddDate(0, 0, -7)
	case "1m":
		since = now.AddDate(0, -1, 0)
	case "3m":
		since = now.AddDate(0, -3, 0)
	default:
		return bson.M{}
	}

	return bson.M{"created_at": bson.M{"$gte": since}}
}
