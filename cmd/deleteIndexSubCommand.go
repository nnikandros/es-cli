package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"serde"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

func deleteIndexCmdFunc(es *elasticsearch.TypedClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "delete subcommand deletes an index(ices) that you provide as args",
		Long:  `deletes the given index provided as an argument (Careful with the use of this command )`,
		RunE:  runDeleteIndexCmdFunc(es),

		ValidArgsFunction: ValidIndexArgsAutoCompletion(es),
	}

	return cmd

}

func runDeleteIndexCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		indexName := ParseArgsIntoString(cmd, args)

		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()
		r, err := es.Indices.Delete(indexName).Do(ctx)
		if err != nil {
			return fmt.Errorf("at deleting the index %w", err)
		}

		if err := json.NewEncoder(cmd.OutOrStdout()).Encode(r); err != nil {
			return serde.SerializingError(err)
		}

		return nil
	}
}
