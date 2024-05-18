package main

import (
	"bytes"
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

func hello(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func getAll(w http.ResponseWriter, req *http.Request) {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db_String := os.Getenv("DB_STRING")
	fmt.Println("db_String-->", db_String)

	type JsonBody struct {
		collection string
		database   string
		dataSource string
	}

	bodyJson := JsonBody{
		collection: "cards",
		database:   "cards",
		dataSource: "Cluster0",
	}

	jsonData, jsonErr := json.Marshal(bodyJson)
	fmt.Println("jsonData-->", jsonData)

	if jsonErr != nil {
		panic(jsonErr)
	}

	resp, err := http.Post(db_String, "application/json", bytes.NewBuffer(jsonData))
	fmt.Println("resp-->", resp)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
}

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MONGO_URI := os.Getenv("MONGO_URI")

	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(MONGO_URI).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	// Send a ping to confirm a successful connection
	if err := client.Database("cards").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	http.HandleFunc("/", hello)
	http.HandleFunc("/get-all", getAll)

	http.ListenAndServe(":8090", nil)
}
