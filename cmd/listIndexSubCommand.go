package cmd

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

// type Response []types.IndicesRecord

type ListSubCmd = *cobra.Command

func listIndicesCmdFunc(es *elasticsearch.TypedClient) ListSubCmd {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list subcommand gets information about available indices",
		Long: `The list subcommand desplays the index names. The --all flag will show all indices, including hidden ones managed by elasticsearch itself.
Enabling tha flag --tab will enable move verbose information such as number of documents, health status etc in a table format.		
`,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		// Run: func(cmd *cobra.Command, args []string) { },
		// Run: runIndices,
		RunE:         runListCmdFunc(es),
		SilenceUsage: true,
		Args:         cobra.NoArgs,
		Example: `es index list
es index list
es index list -a`,
	}

	return cmd

}

func addListFlags(listIndices ListSubCmd) ListSubCmd {
	// listIndices.Flags().BoolP("tab", "t", false, "Display detailed information about each index")
	listIndices.Flags().BoolP("all", "a", false, "List all indices, including hidden or system ones")
	return listIndices

}

func runListCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		all, _ := cmd.Flags().GetBool("all")

		switch {

		case all:
			return listAllIndices(es, cmd.OutOrStdout())
		default:
			return listIndices(es, cmd.OutOrStdout())
		}

	}

}

func listIndices(es *elasticsearch.TypedClient, w io.Writer) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFunc()

	r, err := es.Cat.Indices().Do(ctx)
	if err != nil {
		return fmt.Errorf("at doing request to get the indices %w", err)
	}

	for _, indexRecord := range r {
		if IsIndexNameValid(*indexRecord.Index) {
			fmt.Fprintf(w, "%v\n", *indexRecord.Index)
		}
	}

	return nil

}

// will list all indices incliusing the hidden ones from elasticsearch
func listAllIndices(es *elasticsearch.TypedClient, w io.Writer) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFunc()

	r, err := es.Cat.Indices().Do(ctx)
	if err != nil {
		return fmt.Errorf("at doing request to get the indices %w", err)
	}

	for _, indexRecord := range r {
		fmt.Fprintf(w, "%v\n", *indexRecord.Index)
	}

	return nil

}
