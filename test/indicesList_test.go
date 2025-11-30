package test

import (
	"context"
	"escobra/cmd"
	"fmt"
	"log"
	"os"
	"testing"
	"text/tabwriter"
)

// type Response []types.IndicesRecord

func TestList(t *testing.T) {

	es, _ := cmd.NewElasticTypedClient()

	r, err := es.Cat.Indices().Do(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(r)

	// for _, indexRecord := range r {
	// 	indexName := *indexRecord.Index
	// 	if strings.HasPrefix(indexName, ".") || strings.Contains(indexName, "connector") {
	// 		continue
	// 	}
	// 	fmt.Printf("%v\n", indexName)
	// }

	// b, _ := json.MarshalIndent(r, "", " ")
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	// fmt.Printf("%s", b)
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", "dataset.size", "docs.count", "docs.deleted", "health", "index", "status")
	for _, indexRecord := range r {
		if cmd.IsIndexNameValid(*indexRecord.Index) {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", *indexRecord.DatasetSize, *indexRecord.DocsCount, *indexRecord.DocsDeleted, *indexRecord.Health, *indexRecord.Index, *indexRecord.Status)

		}
	}

	w.Flush()
}
