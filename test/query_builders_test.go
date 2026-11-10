package test

import (
	"escobra/cmd"
	"fmt"
	"reflect"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func hello(a ...int) {

	fmt.Println(a)

}

func TestBuildFilterAndQuery(t *testing.T) {

	f := cmd.BuildTermLevelQuery("LEVEL", []string{"ERROR"})
	f2 := cmd.BuildTermLevelQuery("APP_NAME", []string{"castor", "lion"})
	f3 := cmd.BuildTermLevelQuery("TIMESTAMP", []string{"castor", "lion"})

	queries := []types.Query{*f, *f2, *f3}
	result := cmd.BuildFilterAndQuery(queries)

	if result == nil {
		t.Fatal("expected non-nil query")
	}

	if result.Bool == nil {
		t.Fatal("expected Bool query to be set")
	}

	if len(result.Bool.Filter) != len(queries) {
		t.Fatalf("expected %d filters, got %d",
			len(queries),
			len(result.Bool.Filter),
		)
	}

	for i, q := range queries {
		if !reflect.DeepEqual(q, result.Bool.Filter[i]) {
			t.Errorf("filter at index %d does not match input query", i)
		}
	}
}
