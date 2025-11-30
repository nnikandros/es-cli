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

		ValidArgsFunction: ValidArgsFuncAutoCompletion(es),
	}

	return cmd

}

func runDeleteIndexCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		if len(args) != 1 {
			return fmt.Errorf("you need to provide one index name to delete, you provided: %v", len(args))
		}

		indexName := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
		defer cancel()
		r, err := es.Indices.Delete(indexName).Do(ctx)
		if err != nil {
			return fmt.Errorf("at deleting the index %w", err)
		}

		marshaled, err := json.Marshal(r)
		if err != nil {
			return serde.SerializingError(err)
		}
		fmt.Printf("%s\n", marshaled)

		return nil
	}
}
