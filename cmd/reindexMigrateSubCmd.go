package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"serde"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/spf13/cobra"
)

func reindexMigrateCmdFunc(es *elasticsearch.TypedClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reindex",
		Short: "reindex API",
		Long:  "To be updated",
		RunE:  runReIndexCmdFunc(es),
		// ValidArgsFunction: ValidArgsFuncAutoCompletion(es),
		Example: `es migrate reindex -s <clone-index-1> -s <source-index-2>  -t <target-index>
es migrate reindex --source=<index-1>,<index-2> -t <target-index>`,
	}

	return cmd

}

func addReindexFlags(countCmd *cobra.Command) *cobra.Command {
	countCmd.Flags().StringP("target", "t", "", "Name of clone index that will be created.")
	countCmd.Flags().StringSliceP("source", "s", []string{}, "Name of clone index that will be created.")
	// countCmd.Flags().Int("size", 0, "number of docs to reindex")
	return countCmd

}

func runReIndexCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		sourcIndices, _ := cmd.Flags().GetStringSlice("source")
		targetIndex, _ := cmd.Flags().GetString("target")

		// size, _ := cmd.Flags().GetInt64("size")

		source := &types.ReindexSource{Index: sourcIndices}

		dest := &types.ReindexDestination{Index: targetIndex}

		reindexRequuest := es.Core.Reindex().Source(source).Dest(dest)
		// if size != 0 {
		// 	reindexRequuest = reindexRequuest.MaxDocs(size)
		// }

		resp, err := reindexRequuest.Do(ctx)
		if err != nil {
			return fmt.Errorf("reindx request %w", err)
		}

		if err = json.NewEncoder(cmd.OutOrStdout()).Encode(resp); err != nil {
			return serde.SerializingError(err)
		}

		return nil
	}

}
