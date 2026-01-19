package cmd

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

func infoNodesCmdFunc(es *elasticsearch.TypedClient) NodesSubCmd {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "cluster wide information",
		Long:  `cluster-wide information`,
		RunE:  runNodesCmdFunc(es),
		Args:  cobra.NoArgs,
	}

	return cmd
}

func runNodesCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}

}

// func runNodesSubCmd(es *elasticsearch.TypedClient) RunEFunc {
// 	return func(cmd *cobra.Command, args []string) error {
// 		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 		defer cancel()

// 		infoResponse, err := es.Nodes.Info().Do(ctx)
// 		if err != nil {
// 			return fmt.Errorf("at getting nodes info %w", err)
// 		}

// 		b, err := json.Marshal(infoResponse)
// 		if err != nil {
// 			return serde.SerializingError(err)
// 		}
// 		fmt.Fprintf(cmd.OutOrStdout(), "%s", b)

// 		return nil
// 	}

// }
