package api

import (
	"encoding/json"
	"fmt"
	"gopds/hyperbloom/internal/service"
	"io"
	"log"
	"net/http"
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
		// If there's an error decoding JSON, respond with a Bad Request status
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
		// If there's an error decoding JSON, respond with a Bad Request status
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

// bloomCard handles GET requests to compute approximate cardinality of the key.
// It expects query parameter "key" of type string.
func bloomCard(w http.ResponseWriter, r *http.Request) {
	// Log the request details (HTTP method, URL path, and Content-Type header)
	fmt.Println("[GET]", r.URL.Path, r.Header["Content-Type"])

	// Parse query parameters from the request URL
	queries := r.URL.Query()
	key := queries.Get("key")

	// Check if the 'key' query parameter is present and not empty
	if key != "" {
		// Call service to get the cardinality of the Bloom filter and HyperLogLog for the given key
		bCard, hCard := service.BloomCardinality(key)

		// Format the output string with the cardinality values
		output := fmt.Sprintf("Cardinality (bloom, hyperloglog) = (%d, %d)", bCard, hCard)

		// Write the formatted output string to the HTTP response
		w.Write([]byte(output))
	}
}

// bloomSim handles POST requests to calculate Bloom filter similarity.
// It expects a JSON body with "key_1" and "key_2" fields.
func bloomSim(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[POST]", r.URL.Path, r.Header["Content-Type"])

	// Read the request body
	bytebody, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	// Struct to unmarshal the JSON body
	jsonbody := &struct {
		Key1 string `json:"key_1"`
		Key2 string `json:"key_2"`
	}{}

	// Unmarshal the JSON body into the struct
	if err := json.Unmarshal(bytebody, &jsonbody); err != nil {
		// If there's an error decoding JSON, respond with a Bad Request status
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		log.Println("Error decoding JSON body:", err)
		return
	}

	// Calculate Bloom filter similarity using service function
	sim := service.BloomSimilarity(jsonbody.Key1, jsonbody.Key2)

	// Format the output string with the calculated similarity
	output := fmt.Sprintf("Jaccard similarity = %f", sim)

	// Write the formatted output string to the HTTP response
	w.Write([]byte(output))
}

// bloomBitwiseExists handles POST requests to check bitwise existence in Bloom filters.
// It expects a JSON body with "keys", "value", and "operator" fields.
func bloomBitwiseExists(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[POST]", r.URL.Path, r.Header["Content-Type"])

	// Read the request body
	bytebody, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	// Struct to unmarshal the JSON body
	jsonbody := &struct {
		Keys     []string `json:"keys"`
		Value    string   `json:"value"`
		Operator string   `json:"operator"`
	}{}

	// Unmarshal the JSON body into the struct
	if err := json.Unmarshal(bytebody, &jsonbody); err != nil {
		// If there's an error decoding JSON, respond with a Bad Request status
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		log.Println("Error decoding JSON body:", err)
		return
	}

	// Call service to determine bitwise existence
	bitResult := service.BloomBitwiseExists(
		jsonbody.Keys,
		jsonbody.Value,
		jsonbody.Operator,
	)

	// Prepare output based on bitwise result
	output := fmt.Sprintf("%s bitwise exists = %t", jsonbody.Operator, bitResult)

	// Write response to the client
	w.Write([]byte(output))
}

// bloomChainingExists handles POST requests to check chaining existence in Bloom filters.
// It expects a JSON body with "keys", "value", and "operator" fields.
func bloomChainingExists(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[POST]", r.URL.Path, r.Header["Content-Type"])

	// Read the request body
	bytebody, _ := io.ReadAll(r.Body)
	defer r.Body.Close()

	// Struct to unmarshal the JSON body
	jsonbody := &struct {
		Keys     []string `json:"keys"`
		Value    string   `json:"value"`
		Operator string   `json:"operator"`
	}{}

	// Unmarshal the JSON body into the struct
	if err := json.Unmarshal(bytebody, &jsonbody); err != nil {
		// If there's an error decoding JSON, respond with a Bad Request status
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		log.Println("Error decoding JSON body:", err)
		return
	}

	// Call service to check existence of value in Bloom filters associated with keys
	bitResult := service.BloomChainingExists(
		jsonbody.Keys,
		jsonbody.Value,
		jsonbody.Operator,
	)

	// Format the output string with the calculated result
	output := fmt.Sprintf("%s chaining exists = %t", jsonbody.Operator, bitResult)

	// Write the formatted output string to the HTTP response
	w.Write([]byte(output))
}
