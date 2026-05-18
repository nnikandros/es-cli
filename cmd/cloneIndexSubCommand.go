package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nnikandros/serde"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

type CloneCmd = *cobra.Command

func cloneMigrateCmdFunc(es *elasticsearch.TypedClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "clone",
		Short:   "clone subcommand clones an index(ices) that you provide as args",
		Long:    "To be updated",
		RunE:    runCloneIndexCmdFunc(es),
		Args:    cobra.NoArgs,
		Example: `es clone -t <clone-index> -s <index-to-be-cloned>`,
	}

	return cmd

}

func addCloneFlags(cloneCmd *cobra.Command, es *elasticsearch.TypedClient) *cobra.Command {
	cloneCmd.Flags().StringP("target", "t", "", "Name of the clone")
	cloneCmd.Flags().StringP("source", "s", "", "Name of the index that will be cloned.")

	cloneCmd.RegisterFlagCompletionFunc("target", ValidIndexArgsAutoCompletion(es))
	cloneCmd.RegisterFlagCompletionFunc("source", ValidIndexArgsAutoCompletion(es))

	return cloneCmd

}

// To refactor
func runCloneIndexCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()

		targetIndex, _ := cmd.Flags().GetString("target")
		sourceIndex, _ := cmd.Flags().GetString("source")

		_, err := es.Indices.Close(sourceIndex).Do(ctx)

		if err != nil {
			return fmt.Errorf("at closing source index %v: %w", sourceIndex, err)
		}

		r2, err := es.Indices.Clone(sourceIndex, targetIndex).Do(ctx)
		if err != nil {
			return fmt.Errorf("at cloning the index: %v reason: %w", sourceIndex, err)
		}

		if err = json.NewEncoder(cmd.OutOrStdout()).Encode(r2); err != nil {
			return serde.SerializingError(err)
		}

		r3, err := es.Indices.Open(sourceIndex).Do(ctx)
		if err != nil {
			return fmt.Errorf("at opening the index %v: %w", sourceIndex, err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%+v", r3)
		return nil
	}

}
