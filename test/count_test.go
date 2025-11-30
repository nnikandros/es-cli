package test

import (
	"context"
	"encoding/json"
	"escobra/cmd"
	"fmt"
	"os"
	"testing"
	"text/tabwriter"
)

// type Response []types.CountRecord
// type CountRecord struct {
// 	// Count the document count
// 	Count *string `json:"count,omitempty"`
// 	// Epoch seconds since 1970-01-01 00:00:00
// 	Epoch StringifiedEpochTimeUnitSeconds `json:"epoch,omitempty"`
// 	// Timestamp time in HH:MM:SS
// 	Timestamp *string `json:"timestamp,omitempty"`
// }

func TestCount(t *testing.T) {

	t.Run("count", func(t *testing.T) {
		es, _ := cmd.NewElasticTypedClient()

		r, _ := es.Cat.Count().Index("*").Do(context.Background())

		parsed, _ := json.Marshal(r)
		fmt.Printf("%s", parsed)
		response := r[0]

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

		fmt.Fprintf(w, "%s\t%s\t%s\n", "count", "epoch", "timestamp")
		fmt.Fprintf(w, "%s\t%s\t%s\n", *response.Count, response.Epoch, *response.Timestamp)

		// Flush the output
		w.Flush()

		// if err != nil {
		// 	log.Fatal(err)
		// }

		// b, _ := json.MarshalIndent(r, "", " ")
		// fmt.Printf("%s", b)
	})

}
