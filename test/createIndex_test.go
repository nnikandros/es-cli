package test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// es.Indices.Create(arg).Mappings().Settings()

func TestCreateIndex(t *testing.T) {

	// es, _ := cmd.NewElasticTypedClient()

	m := "~/elasticops/settings_mappings/test-index/mappings.json"
	s := "~/elasticops/settings_mappings/test-index/settings.json"

	b, _ := os.ReadFile(m)

	b2, _ := os.ReadFile(s)

	var mappings types.TypeMapping
	var settings types.IndexSettings

	err := json.Unmarshal(b, &mappings)
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(b2, &settings)
	if err != nil {
		log.Fatal(err)
	}

	// r, err := es.Indices.Create("test-index-to-be-deleted").Mappings(&mappings).Settings(&settings).Do(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// es.Indices.Delete("test-index-to-be-deleted").Do(context.Background())

	// fmt.Printf("%+v\n", r)
}
