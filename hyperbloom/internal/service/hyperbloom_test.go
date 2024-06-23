package service_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/nnurry/gopds/hyperbloom/internal/database/postgres"
	"github.com/nnurry/gopds/hyperbloom/internal/service"
)

func TestCreateBloom(t *testing.T) {
	var err error
	createQuery := "CREATE TABLE IF NOT EXISTS bloom_filters (key VARCHAR, bloombyte BYTEA, hyperbyte BYTEA)"
	_, err = postgres.DbClient.Exec(createQuery)
	if err != nil {
		log.Fatal("Can't create table bloom_filters", err)
	}
	ids := map[string][]string{
		"123": {"123", "456", "789"},
		"456": {"987", "654", "321"},
	}

	fmt.Println("Num of keys before hashing:", len(service.BloomList()))

	for key := range ids {
		for _, id := range ids[key] {
			service.BloomHash(key, id)
			outStr := fmt.Sprintf("New appr. size (%s) = %d", key, service.BloomGet(key).BloomCardinality())
			fmt.Println(outStr)
		}
	}

	mergedIds := []string{}

	for key := range ids {
		mergedIds = append(mergedIds, ids[key]...)
	}

	for key := range ids {
		b := service.BloomGet(key)
		for _, id := range mergedIds {
			fmt.Printf("(%s) ⪽ (%s) = %t\n", id, key, b.CheckExists(id))
		}
		fmt.Println("Bloom card:", b.BloomCardinality())
		fmt.Println("Hyper card:", b.HyperCardinality())
	}

	fmt.Println("Num of keys after hashing:", len(service.BloomList()))
}
