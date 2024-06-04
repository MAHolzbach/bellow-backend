package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Job struct {
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func getAll(w http.ResponseWriter, req *http.Request) {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Failed to load env file -->", envErr)
	}

	MONGO_URI := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(MONGO_URI)

	client, connectionErr := mongo.Connect(context.TODO(), clientOptions)

	if connectionErr != nil {
		log.Fatal(connectionErr)
	}

	var results []*Job

	collection := client.Database("cards").Collection("cards")

	cur, _ := collection.Find(context.TODO(), bson.D{{}})

	for cur.Next(context.TODO()) {

		var elem Job

		err := cur.Decode(&elem)

		if err != nil {
			log.Fatal(err)
		}

		results = append(results, &elem)

	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	cur.Close(context.TODO())

	fmt.Printf("Found multiple documents (array of pointers): %+v\n", results)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(results)
}

func postOne(w http.ResponseWriter, req *http.Request) {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Failed to load env file -->", envErr)
	}

	MONGO_URI := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(MONGO_URI)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("cards").Collection("cards")

	newJobCard := Job{"New Job", "Test New Job"}

	insertResult, insertOneErr := collection.InsertOne(context.TODO(), newJobCard)

	if insertOneErr != nil {
		log.Fatal(insertOneErr)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}

func main() {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Failed to load env file -->", envErr)
	}

	MONGO_URI := os.Getenv("MONGO_URI")

	clientOptions := options.Client().ApplyURI(MONGO_URI)

	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	http.HandleFunc("/", hello)
	http.HandleFunc("/post-one", postOne)
	http.HandleFunc("/get-all", getAll)

	log.Fatal(http.ListenAndServe(":8090", nil))
}
