package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nnurry/gopds/hyperbloom/internal/service"
)

// bloomHash handles POST requests for hashing a value and adding it to the Bloom filter.
// It expects a JSON body with "key" and "value" fields.
func bloomHash(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[POST]", r.URL.Path, r.Header["Content-Type"])
	// Read the request body
	bytebody, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	// Struct to unmarshal the JSON body
	jsonbody := &struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{}

	// Unmarshal the JSON body into the struct
	if err := json.Unmarshal(bytebody, &jsonbody); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		log.Println("Error decoding JSON body:", err)
		return
	}

	// Add the value to the Bloom filter using the provided key
	service.BloomHash(jsonbody.Key, jsonbody.Value)

	// Get the cardinality of the Bloom filter and HyperLogLog
	bCard, hCard := service.BloomCardinality(jsonbody.Key)

	// Format the output string
	output := fmt.Sprintf("Cardinality (bloom, hyperloglog) = (%d, %d)", bCard, hCard)

	// Write the output string to the response
	w.Write([]byte(output))
}

// bloomExists handles POST requests to check if a value exists in the Bloom filter.
// It expects a JSON body with "key" and "value" fields.
func bloomExists(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[POST]", r.URL.Path, r.Header["Content-Type"])
	// Read the request body
	bytebody, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	// Struct to unmarshal the JSON body
	jsonbody := &struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}{}

	// Unmarshal the JSON body into the struct
	if err := json.Unmarshal(bytebody, &jsonbody); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		log.Println("Error decoding JSON body:", err)
		return
	}

	// Check if the value exists in the Bloom filter using the provided key
	exists := service.BloomExists(jsonbody.Key, jsonbody.Value)

	// Format the output string
	output := fmt.Sprintf(
		"(%s) ⪽ (%s) = %t\n",
		jsonbody.Value,
		jsonbody.Key,
		exists,
	)

	// Write the output string to the response
	w.Write([]byte(output))
}

// main sets up the HTTP server and routes
func main() {
	var err error

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Register the bloomHash handler for the /hyperbloom/hash endpoint
	mux.HandleFunc("/hyperbloom/hash", bloomHash)
	// Register the bloomExists handler for the /hyperbloom/exists endpoint
	mux.HandleFunc("/hyperbloom/exists", bloomExists)
	// Register a default handler that prints the requested URL path
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
	})

	// Start the HTTP server on port 5000
	err = http.ListenAndServe(":5000", mux)
	log.Fatal(err) // Log fatal error if the server fails to start
}
