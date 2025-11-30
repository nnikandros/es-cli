package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

func countCmdFunc(es *elasticsearch.TypedClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count",
		Short: "count API",
		Long:  "Retrieves the number of documents in an index pattern. You can provide a pattern with wildcard, a comma-separated or white space separated list of indices\nTo target all indices you can pass no argument or the string _all",
		RunE:  runCountCmdFunc(es),
		Example: `es count
es count test-index-* -v
es count index-1 index-2 --tab
es count index-1,index-2 --tab
`,
	}

	return cmd

}

func addCountFlags(countCmd *cobra.Command) *cobra.Command {
	countCmd.Flags().BoolP("tab", "t", false, "adds the number of documents in the index pattern with a timestamp and displays it in a table format")
	return countCmd

}

func runCountCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		index := ParseArgsIntoSingleString(args)

		tabular, _ := cmd.Flags().GetBool("tab")

		switch tabular {
		case true:
			return getCountTabular(es, index)

		case false:
			return getCount(es, index)

		}

		return nil
	}

}

func getCount(es *elasticsearch.TypedClient, indexName string) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	r, err := es.Cat.Count().Index(indexName).Do(ctx)
	if err != nil {
		return fmt.Errorf("at getting count %w", err)
	}

	fmt.Printf("%s\n", *r[0].Count)

	return nil

}

func getCountTabular(es *elasticsearch.TypedClient, indexName string) error {
	response, err := es.Cat.Count().Index(indexName).Do(context.Background())
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintf(w, "%s\t%s\t%s\n", "count", "epoch", "timestamp")
	fmt.Fprintf(w, "%s\t%s\t%s\n", *response[0].Count, response[0].Epoch, *response[0].Timestamp)

	// Flush the output
	w.Flush()

	return nil

}
