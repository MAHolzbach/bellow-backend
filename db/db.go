package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func connectToMongo() (*mongo.Client, error) {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Failed to load env -->", envErr)
	}

	MONGO_DB_URI := os.Getenv("MONGO_DB_URI")

	clientOptions := options.Client().ApplyURI(MONGO_DB_URI)

	// username := os.Getenv("MONGO_DB_USERNAME")
	// password := os.Getenv("MONGO_DB_PASSWORD")

	// clientOptions.SetAuth(options.Credential{
	// 	Username: username,
	// 	Password: password,
	// })

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	return client, nil
}
