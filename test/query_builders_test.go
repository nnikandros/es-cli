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

func TestQueryBuilders(t *testing.T) {
	//

	// q := cmd.BuildTermLevelQuery("LEVEL", []string{"first,second"})
	// // q2 := cmd.BuildTermLevelQuery("APP_NAME", []string{"first app"})

	// b, err := json.Marshal(q)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	f := cmd.BuildTermLevelQueryV2("LEVEL", []string{"first", "second"})
	f2 := cmd.BuildTermLevelQueryV2("LEVEL", []string{"first", "second"})

	// b2, err := json.Marshal(f)

	m := cmd.MustQuery(f, f2)

	// if err := json.NewEncoder(os.Stdout).Encode(m); err != nil {
	// 	log.Fatal(err)
	// }

	b, err := json.MarshalIndent(m, "", " ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s", b)

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
