package test

import (
	"context"
	"encoding/json"
	"escobra/cmd"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
)

func Piki(args ...int) {

	fmt.Println(args)

}

func TestSearchAFter(t *testing.T) {

	typedClient, err := cmd.NewElasticTypedClient()
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now().Unix()

	fmt.Println(now)
	// time.Parse(time.DateOnly, u.FinishedDate)

	// time1, err := time.Parse(time.RFC3339, now)
	// time1 := time.Now().Format(time.RFC3339)
	time1 := time.Now()
	fmt.Println(time1)

	// z, err := time.Parse(time.RFC3339, "2025-10-03T09:34:03.402Z")

	sortOptions := types.SortOptions{
		SortOptions: map[string]types.FieldSort{
			"TIMESTAMP": {
				Order: &sortorder.Asc,
			},
		},
	}
	r, err := typedClient.Search().Index("castor-test-log-v8").Sort(sortOptions).SearchAfter([]types.FieldValue{1759493965636}...).Do(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	b, _ := json.MarshalIndent(r, "", " ")
	fmt.Printf("%s", b)

}
