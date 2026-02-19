//go:build integration

package test

import (
	"context"
	"encoding/json"
	"escobra/cmd"
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

	// field := "CASE_CODE"
	// field2 := "COJ_DOC_TYPE"

	// aggs := map[string]types.Aggregations{field: {Terms: &types.TermsAggregation{Field: &field}}, field2: {Terms: &types.TermsAggregation{Field: &field2}}}
	// typedClient.Search().Aggregations(aggs)

	q := []types.MultiTermLookup{{Field: "CASE_CODE"}, {Field: "COJ_DOC_TYPE"}}

	m := types.MultiTermsAggregation{Terms: q}
	aggs2 := map[string]types.Aggregations{"case_code_and_coj": {MultiTerms: &m}}

	r, err := typedClient.Search().Index("castor-test-alldecidionspublic-with-coj-v1").Aggregations(aggs2).Size(0).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	if err = json.NewEncoder(t.Output()).Encode(r); err != nil {
		log.Fatal(err)
	}
	// x := Aggregations{}

	// c := r.Aggregations[field]
	// b, err := json.Marshal(c)
	// if err != nil {
	// 	t.Error(err)
	// }
	// // if err = json.NewEncoder(t.Output()).Encode(c); err != nil {
	// // 	log.Fatal(err)
	// // }

	// // var x map[string]any

	// err = json.Unmarshal(b, &x)
	// if err != nil {
	// 	t.Error(err)
	// }

	// fmt.Printf("%+v\n", x)

	// keys := make([]string, 0)

	// for _, v := range x.Buckets {
	// 	keys = append(keys, v.Key)
	// }

	// fmt.Println(keys)

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
