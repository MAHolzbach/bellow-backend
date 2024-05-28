// package main

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	"github.com/joho/godotenv"
// 	"go.mongodb.org/mongo-driver/bson"
// 	"go.mongodb.org/mongo-driver/mongo"
// 	"go.mongodb.org/mongo-driver/mongo/options"
// )

// func hello(w http.ResponseWriter, req *http.Request) {
// 	fmt.Fprintf(w, "hello\n")
// }

// func headers(w http.ResponseWriter, req *http.Request) {
// 	for name, headers := range req.Header {
// 		for _, h := range headers {
// 			fmt.Fprintf(w, "%v: %v\n", name, h)
// 		}
// 	}
// }

// func getAll(w http.ResponseWriter, req *http.Request) {
// 	err := godotenv.Load()

// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	db_String := os.Getenv("DB_STRING")
// 	fmt.Println("db_String-->", db_String)

// 	type JsonBody struct {
// 		collection string
// 		database   string
// 		dataSource string
// 	}

// 	bodyJson := JsonBody{
// 		collection: "cards",
// 		database:   "cards",
// 		dataSource: "Cluster0",
// 	}

// 	jsonData, jsonErr := json.Marshal(bodyJson)
// 	fmt.Println("jsonData-->", jsonData, jsonErr)

// 	if jsonErr != nil {
// 		panic(jsonErr)
// 	}

// 	resp, err := http.Post(db_String, "application/json", bytes.NewBuffer(jsonData))
// 	fmt.Println("resp-->", resp)

// 	if err != nil {
// 		panic(err)
// 	}

// 	defer resp.Body.Close()

// 	fmt.Println("Response status:", resp.Status)
// }

// var client *mongo.Client

// type JobCard struct {
// 	title       string `json:"title"`
// 	description string `json:"description"`
// }

// func getAllHandler(w http.ResponseWriter, r *http.Request) {
// 	collection := client.Database("cards").Collection("cards")
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	cursor, err := collection.Find(ctx, bson.M{})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer cursor.Close(ctx)

// 	var results []JobCard
// 	if err = cursor.All(ctx, &results); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(results)
// }

// func main() {
// 	err := godotenv.Load()

// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	MONGO_URI := os.Getenv("MONGO_URI")

// 	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
// 	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

// 	opts := options.Client().ApplyURI(MONGO_URI).SetServerAPIOptions(serverAPI)

// 	// Create a new client and connect to the server
// 	client, err := mongo.Connect(context.TODO(), opts)
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer func() {
// 		if err = client.Disconnect(context.TODO()); err != nil {
// 			panic(err)
// 		}
// 	}()

// 	// Send a ping to confirm a successful connection
// 	if err := client.Database("cards").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

// 	http.HandleFunc("/", hello)
// 	http.HandleFunc("/get-all", getAll)

// 	log.Fatal(http.ListenAndServe(":8090", nil))
// }

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	// "go.mongodb.org/mongo-driver/bson"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Job struct {
	Title       string
	Description string
}

func main() {
	envErr := godotenv.Load()

	if envErr != nil {
		log.Fatal("Failed to load env file -->", envErr)
	}

	MONGO_URI := os.Getenv("MONGO_URI")
	// Set client options
	clientOptions := options.Client().ApplyURI(MONGO_URI)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("cards").Collection("cards")

	newJobCard := Job{"New Job", "Test New Job"}

	insertResult, insertOneErr := collection.InsertOne(context.TODO(), newJobCard)

	if insertOneErr != nil {
		log.Fatal(insertOneErr)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
}
