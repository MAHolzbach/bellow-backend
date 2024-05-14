package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	db_String := os.Getenv("DB_STRING")

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

	if jsonErr != nil {
		panic(jsonErr)
	}

	resp, err := http.Post(db_String, "application/json", bytes.NewBuffer(jsonData))

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	fmt.Println("Response status:", resp.Status)
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/getAll", getAll)

	http.ListenAndServe(":8090", nil)
}
