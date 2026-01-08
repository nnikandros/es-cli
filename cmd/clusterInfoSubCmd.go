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

func infoClusterCmdFunc(es *elasticsearch.TypedClient) InfoSubCmd {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "info",
		Long:  "info info",
		RunE:  runInfoClusterSubcmd(es),
		Args:  cobra.NoArgs,
	}

	return cmd

}

func runInfoClusterSubcmd(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		r, err := es.Cluster.Info("_all").Do(ctx)

		if err != nil {
			return fmt.Errorf("at getting cluster info %w", err)
		}

		err = json.NewEncoder(cmd.OutOrStdout()).Encode(r)
		if err != nil {
			return serde.SerializingError(err)
		}

		return nil
	}

}
