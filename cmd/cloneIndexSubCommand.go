package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

func cloneIndexCmdFunc(es *elasticsearch.TypedClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "clone",
		Short:             "clone subcommand clones an index(ices) that you provide as args",
		Long:              `clones the given index provided as an argument (Careful with the use of this command )`,
		RunE:              runCloneIndexCmdFunc(es),
		ValidArgsFunction: ValidArgsFuncAutoCompletion(es),
	}

	return cmd

}

func addCloneFlags(countCmd *cobra.Command) *cobra.Command {
	countCmd.Flags().StringP("target", "t", "", "adds the number of documents in the index pattern with a timestamp and displays it in a table format")
	return countCmd

}

func runCloneIndexCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		sourceIndex := args[0]
		targetIndex, _ := cmd.Flags().GetString("target")

		_, err := es.Indices.Close(sourceIndex).Do(ctx)

		if err != nil {
			return fmt.Errorf("at closing %w", err)
		}

		_, err = es.Indices.Clone(sourceIndex, targetIndex).Do(ctx)
		if err != nil {
			return fmt.Errorf("at cloning the index: %v reason: %w", sourceIndex, err)
		}

		r3, err := es.Indices.Open(sourceIndex).Do(ctx)
		if err != nil {
			return fmt.Errorf("at opening the index %v: %w", sourceIndex, err)
		}

		fmt.Println(r3)
		return nil
	}

}
