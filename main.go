package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/organisasi/tubesbackend/routes"
)

var DB *mongo.Database

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	mongoURI := os.Getenv("MONGOSTRING")
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	DB = client.Database("tubesbackend")
	log.Println("Connected to MongoDB")
}

func main() {
	router := routes.SetupRoutes(DB)
	log.Fatal(router.Run(":8080")) // Run server on port 8080
}
