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

func infoNodesCmdFunc(es *elasticsearch.TypedClient) NodesSubCmd {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "nodes information",
		Long:  `nodes information`,
		RunE:  runNodesSubCmd(es),
		Args:  cobra.NoArgs,
	}

	return cmd
}

func addInfoNodesFlags(nodesSub NodesSubCmd) NodesSubCmd {

	nodesSub.Flags().Bool("stats", false, "Get node statistics. Get statistics for nodes in a cluster. By default, all stats are returned.")
	nodesSub.Flags().Bool("usage", false, "Get feature usage information. ")

	return nodesSub
}

func runNodesSubCmd(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		stats, _ := cmd.Flags().GetBool("stats")
		usage, _ := cmd.Flags().GetBool("usage")

		switch {
		case stats:
			statsResponse, err := es.Nodes.Stats().Do(ctx)
			if err != nil {
				return fmt.Errorf("at getting nodes stats")
			}

			if err := json.NewEncoder(cmd.OutOrStdout()).Encode(statsResponse); err != nil {
				return serde.SerializingError(err)
			}

			return nil

		case usage:
			usageResponse, err := es.Nodes.Usage().Do(ctx)
			if err != nil {
				return fmt.Errorf("at getting nodes usage")
			}

			if err := json.NewEncoder(cmd.OutOrStdout()).Encode(usageResponse); err != nil {
				return serde.SerializingError(err)
			}

			return nil

		default:
			infoResponse, err := es.Nodes.Info().Do(ctx)
			if err != nil {
				return fmt.Errorf("at getting nodes info %w", err)
			}

			if err := json.NewEncoder(cmd.OutOrStdout()).Encode(infoResponse); err != nil {
				return serde.SerializingError(err)
			}

			return nil

		}

	}

}
