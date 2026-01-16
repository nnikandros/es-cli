package test

import (
	"encoding/json"
	"escobra/cmd"
	"fmt"
	"log"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

func hello(a ...int) {

	fmt.Println(a)

}

func TestFilterQueryBuilders(t *testing.T) {

	f := cmd.BuildTermLevelQuery("LEVEL", []string{"ERROR"})
	f2 := cmd.BuildTermLevelQuery("APP_NAME", []string{"castor", "lion"})
	f3 := cmd.BuildTermLevelQuery("TIMESTAMP", []string{"castor", "lion"})

	filterQuery := cmd.BuildFilterAndQuery([]types.Query{*f, *f2, *f3})
	b, err := json.MarshalIndent(filterQuery, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)

}
