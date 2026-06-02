package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	_ = godotenv.Load()
	_ = godotenv.Load("../../.env")

	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	dbName := os.Getenv("MONGODB_DB_NAME")
	if dbName == "" {
		dbName = "ufriend_cx"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(dbName)
	col := db.Collection("auth_sessions")

	var session bson.M
	err = col.FindOne(ctx, bson.M{
		"expires_at": bson.M{"$gt": time.Now()},
	}, options.FindOne().SetSort(bson.M{"expires_at": -1})).Decode(&session)
	if err != nil {
		// If no active session, try to find any session
		err = col.FindOne(ctx, bson.M{}).Decode(&session)
		if err != nil {
			fmt.Println("No session found in database.")
			return
		}
	}

	fmt.Printf("Active Session ID: %s\n", session["session_id"])
	fmt.Printf("Expires At: %v\n", session["expires_at"])
}
