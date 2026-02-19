package test

import (
	"context"
	"encoding/json"
	"escobra/cmd"
	"fmt"
	"log"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Aggregations struct {
	Buckets []struct {
		DocCount int    `json:"doc_count"`
		Key      string `json:"key"`
	} `json:"buckets"`
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
}

func TestAggegations(t *testing.T) {
	typedClient, err := cmd.NewElasticTypedClient()
	if err != nil {
		log.Fatal(err)
	}

	aggs := map[string]types.Aggregations{"CASE_CODE": {Terms: &types.TermsAggregation{Field: new("CASE_CODE")}}}
	// typedClient.Search().Aggregations(aggs)

	r, err := typedClient.Search().Index("castor-test-alldecidionspublic-with-coj-v1").Aggregations(aggs).Size(0).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// if err = json.NewEncoder(t.Output()).Encode(r); err != nil {
	// 	log.Fatal(err)
	// }
	x := Aggregations{}

	c := r.Aggregations["CASE_CODE"]
	b, err := json.Marshal(c)
	if err != nil {
		t.Error(err)
	}
	// if err = json.NewEncoder(t.Output()).Encode(c); err != nil {
	// 	log.Fatal(err)
	// }

	// var x map[string]any

	err = json.Unmarshal(b, &x)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", x)

	keys := make([]string, 0)

	for _, v := range x.Buckets {
		keys = append(keys, v.Key)
	}

	fmt.Println(keys)

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
