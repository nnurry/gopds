package api

import "net/http"

// ServeHyperBloom registers HTTP request handlers for specific endpoints related to HyperBloom operations.
func ServeHyperBloom(mux *http.ServeMux) {
	// Register various HTTP request handlers for specific endpoints

	// Handler for hashing a value and adding it to the Bloom filter
	mux.HandleFunc("/hyperbloom/hash", bloomHash)

	// Handler for checking if a value exists in the Bloom filter
	mux.HandleFunc("/hyperbloom/exists", bloomExists)

	// Handler for bitwise existence check in Bloom filters associated with multiple keys
	mux.HandleFunc("/hyperbloom/exists/bitwise", bloomBitwiseExists)

	// Handler for chaining existence check in Bloom filters associated with multiple keys
	mux.HandleFunc("/hyperbloom/exists/chaining", bloomChainingExists)

	// Handler for computing approximate cardinality of a Bloom filter and HyperLogLog for a given key
	mux.HandleFunc("/hyperbloom/card", bloomCard)

	// Handler for calculating Jaccard similarity using Bloom filters for two different keys
	mux.HandleFunc("/hyperbloom/sim", bloomSim)
}

// Serve is a wrapper function that calls ServeHyperBloom to register HTTP request handlers.
// It provides a convenient way to initialize the server with the desired handlers.
func Serve(mux *http.ServeMux) {
	ServeHyperBloom(mux)
}
