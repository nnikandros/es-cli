package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"serde"
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
		// Args:  cobra.A,
	}

	return cmd
}

func addInfoIndexFlags(indexInfoSub IndexSubCmd) IndexSubCmd {
	indexInfoSub.Flags().BoolP("all", "a", false, "display all info for all indices, including monitor")

	return indexInfoSub

}

func runIndexInfoCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		t, _ := cmd.Flags().GetBool("all")

		indexName := ParseArgsIntoString(cmd, args)

		if len(indexName) > 0 {
			ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelFunc()

			r, err := es.Cat.Indices().Index(indexName).Do(ctx)
			if err != nil {
				return fmt.Errorf("at getting index info %w", err)
			}

			if err := json.NewEncoder(cmd.OutOrStdout()).Encode(r); err != nil {
				return serde.SerializingError(err)
			}

			return nil
		}

		switch {
		default:
			if err := listIndicesTabular(es, cmd.OutOrStderr()); err != nil {
				return fmt.Errorf("at listing indices in a table format %w", err)
			}

			return nil

		case t:
			ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelFunc()
			r, err := es.Cat.Indices().Do(ctx)
			if err != nil {
				return fmt.Errorf("at doing request to get the indices %w", err)
			}

			if err := json.NewEncoder(cmd.OutOrStdout()).Encode(r); err != nil {
				return serde.SerializingError(err)
			}

			return nil

		}

	}

}

func listIndicesTabular(es *elasticsearch.TypedClient, w io.Writer) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFunc()

	r, err := es.Cat.Indices().Do(ctx)
	if err != nil {
		return fmt.Errorf("at doing request to get the indices %w", err)
	}

	tbW := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tbW, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", "index", "health", "docs.count", "docs.deleted", "dataset.size", "status", "primary_shards", "primary_size")
	for _, indexRecord := range r {
		if IsIndexNameValid(*indexRecord.Index) {
			fmt.Fprintf(tbW, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", *indexRecord.Index, *indexRecord.Health, *indexRecord.DocsCount, *indexRecord.DocsDeleted, *indexRecord.DatasetSize, *indexRecord.Status, *indexRecord.Pri, *indexRecord.PriStoreSize)

		}
	}

	tbW.Flush()

	return nil

	// }

	// // this wil print all indices including the ones from elasticsearch
	// func listIndicesTabularaAll(es *elasticsearch.TypedClient, w io.Writer) error {
	// 	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	// 	defer cancelFunc()

	// 	r, err := es.Cat.Indices().Do(ctx)
	// 	if err != nil {
	// 		return fmt.Errorf("at doing request to get the indices %w", err)
	// 	}

	// 	tbW := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	// 	fmt.Fprintf(tbW, "%s\t%s\t%s\t%s\t%s\t%s\n", "index", "health", "docs.count", "docs.deleted", "dataset.size", "status")
	// 	for _, indexRecord := range r {
	// 		fmt.Fprintf(tbW, "%s\t%s\t%s\t%s\t%s\t%s\n", *indexRecord.Index, *indexRecord.Health, *indexRecord.DocsCount, *indexRecord.DocsDeleted, *indexRecord.DatasetSize, *indexRecord.Status)
	// 	}

	// 	tbW.Flush()

	// 	return nil

}
