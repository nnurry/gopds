package service_test

import (
	"fmt"
	"testing"

	"gopds/hyperbloom/internal/service"
)

func TestCreateBloom(t *testing.T) {
	// Define some initial IDs for testing
	ids := map[string][]string{
		"123": {"123", "456", "789"},
		"456": {"987", "654", "321"},
	}

	// Print the number of keys in the bloom filter list before hashing
	fmt.Println("Num of keys before hashing:", len(service.BloomList()))

	// Iterate over each key in 'ids'
	for key := range ids {
		// Iterate over each ID associated with the current key
		for _, id := range ids[key] {
			// Hash the 'id' using the 'key'
			service.BloomHash(key, id)
			// Print the approximate size of the bloom filter for the current 'key'
			outStr := fmt.Sprintf("New appr. size (%s) = %d", key, service.BloomGet(key).BloomCardinality())
			fmt.Println(outStr)
		}
	}

	// Merge all IDs into a single slice for further comparison
	mergedIds := []string{}
	for key := range ids {
		mergedIds = append(mergedIds, ids[key]...)
	}

	// Iterate over each key in 'ids' again
	for key := range ids {
		// Retrieve the bloom filter for the current 'key'
		b := service.BloomGet(key)
		// Check each merged ID against the current bloom filter
		for _, id := range mergedIds {
			fmt.Printf("(%s) ⪽ (%s) = %t\n", id, key, b.CheckExists(id))
		}
		// Print the cardinality of the bloom filter for the current 'key'
		fmt.Println("Bloom card:", b.BloomCardinality())
		// Print the cardinality of the hyperbloom for the current 'key'
		fmt.Println("Hyper card:", b.HyperCardinality())
	}

	// Print the number of keys in the bloom filter list after hashing
	fmt.Println("Num of keys after hashing:", len(service.BloomList()))
}
