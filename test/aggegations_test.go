package test

import (
	"context"
	"encoding/json"
	"escobra/cmd"
	"log"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

type Bucket struct {
	DocCount int    `json:"doc_count"`
	Key      string `json:"key"`
}

type Aggregations struct {
	Buckets                 []Bucket `json:"buckets"`
	DocCountErrorUpperBound int      `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int      `json:"sum_other_doc_count"`
}

func TestAggegations(t *testing.T) {
	typedClient, err := cmd.NewElasticTypedClient()
	if err != nil {
		log.Fatal(err)
	}

	field := "CASE_CODE"
	field2 := "COJ_DOC_TYPE"

	// aggs := map[string]types.Aggregations{field: {Terms: &types.TermsAggregation{Field: &field}}}
	// aggs := map[string]types.Aggregations{field: {Terms: &types.TermsAggregation{Field: &field}}, field2: {Terms: &types.TermsAggregation{Field: &field2}}}
	// typedClient.Search().Aggregations(aggs)

	q := []types.MultiTermLookup{{Field: field}, {Field: field2}}
	m := types.MultiTermsAggregation{Terms: q}
	aggs2 := map[string]types.Aggregations{"case_code_and_coj": {MultiTerms: &m}}

	r, err := typedClient.Search().Index("castor-test-alldecidionspublic-with-coj-v1").Aggregations(aggs2).Size(0).Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	// s, ok := r.Aggregations[field]
	// if !ok {
	// 	t.Error(ok)
	// }
	// b, err := json.Marshal(s)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// var a = Aggregations{}

	// err = json.Unmarshal(b, &a)
	// if err != nil {
	// 	t.Errorf("unmarashal %v", err)
	// }
	// var uniqueKeys []string
	// for _, k := range a.Buckets {
	// 	uniqueKeys = append(uniqueKeys, k.Key)
	// }

	// fmt.Println(uniqueKeys)
	// var aggsResults Aggregations
	// json.Unmarshal()
	if err = json.NewEncoder(t.Output()).Encode(r); err != nil {
		log.Fatal(err)
	}
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

// {"aggregations":{"CASE_CODE":{"buckets":[{"doc_count":3,"key":"C-162/15 P"},{"doc_count":3,"key":"C-209/78"},{"doc_count":3,"key":"C-278/00"},{"doc_count":2,"key":"C-10/71"},{"doc_count":2,"key":"C-100/80"},{"doc_count":2,"key":"C-101/15 P"},{"doc_count":2,"key":"C-102/77"},{"doc_count":2,"key":"C-103/84"},{"doc_count":2,"key":"C-105/04 P"},{"doc_count":2,"key":"C-105/14"}],"doc_count_error_upper_bound":4,"sum_other_doc_count":1976}},"hits":{"hits":[],"total":{"relation":"eq","value":1999}},"_shards":{"failed":0,"skipped":0,"successful":2,"total":2},"timed_out":false,"took":1}
