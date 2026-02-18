package test

import (
	"context"
	"encoding/json"
	"escobra/cmd"
	"log"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func TestAggegations(t *testing.T) {
	typedClient, err := cmd.NewElasticTypedClient()
	if err != nil {
		log.Fatal(err)
	}

	aggs := map[string]types.Aggregations{"case_codes": {Terms: &types.TermsAggregation{Field: new("CASE_CODE")}}}
	// typedClient.Search().Aggregations(aggs)

	r, err := typedClient.Search().Index("castor-test-alldecidionspublic-with-coj-v1").Aggregations(aggs).Size(0).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// if err = json.NewEncoder(t.Output()).Encode(r); err != nil {
	// 	log.Fatal(err)
	// }

	c := r.Aggregations["case_codes"]

	if err = json.NewEncoder(t.Output()).Encode(c); err != nil {
		log.Fatal(err)
	}

	// if c.StringTerms == nil {
	// 	t.Errorf("not a StringTermsAggregate")
	// }

	// _, ok := c.(types.StringTermsAggregate)

	// if !ok {
	// 	t.Errorf("not ok")
	// }

	// fmt.Println(ok)

	// fmt.Printf("%+v", v)
}
