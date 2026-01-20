package cmd

import (
	"context"
	"fmt"
	"io"
	"text/tabwriter"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

func infoIndexCmdFunc(es *elasticsearch.TypedClient) IndexSubCmd {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "information about indices",
		Long:  `indicies`,
		RunE:  runIndexInfoCmdFunc(es),
		Args:  cobra.NoArgs,
	}

	return cmd
}

func addInfoIndexFlags() {

}

func runIndexInfoCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		return nil

	}

}

// not finished
func listIndicesTabular(es *elasticsearch.TypedClient, w io.Writer) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFunc()

	r, err := es.Cat.Indices().Do(ctx)
	if err != nil {
		return fmt.Errorf("at doing request to get the indices %w", err)
	}

	tbW := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tbW, "%s\t%s\t%s\t%s\t%s\t%s\n", "index", "health", "docs.count", "docs.deleted", "dataset.size", "status")
	for _, indexRecord := range r {
		if IsIndexNameValid(*indexRecord.Index) {
			fmt.Fprintf(tbW, "%s\t%s\t%s\t%s\t%s\t%s\n", *indexRecord.Index, *indexRecord.Health, *indexRecord.DocsCount, *indexRecord.DocsDeleted, *indexRecord.DatasetSize, *indexRecord.Status)

		}
	}

	tbW.Flush()

	return nil

}

// this wil print all indices including the ones from elasticsearch
func listIndicesTabularaAll(es *elasticsearch.TypedClient, w io.Writer) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFunc()

	r, err := es.Cat.Indices().Do(ctx)
	if err != nil {
		return fmt.Errorf("at doing request to get the indices %w", err)
	}

	tbW := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tbW, "%s\t%s\t%s\t%s\t%s\t%s\n", "index", "health", "docs.count", "docs.deleted", "dataset.size", "status")
	for _, indexRecord := range r {
		fmt.Fprintf(tbW, "%s\t%s\t%s\t%s\t%s\t%s\n", *indexRecord.Index, *indexRecord.Health, *indexRecord.DocsCount, *indexRecord.DocsDeleted, *indexRecord.DatasetSize, *indexRecord.Status)
	}

	tbW.Flush()

	return nil

}
