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

func infoclusterCmdFunc(es *elasticsearch.TypedClient) ClusterCmd {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "cluster wide information",
		Long:  `cluster-wide information`,
		RunE:  runInfoClusterCmdFunc(es),
		Args:  cobra.NoArgs,
	}

	return cmd
}

func addInfoClusterFlags(infoSub InfoSubCmd) ClusterCmd {

	// infoSub.Flags().StringP("target", "t", "_all", "target for cluster information")
	infoSub.Flags().Bool("settings", false, "get cluster settings")
	infoSub.Flags().Bool("info", true, "get cluster info")
	infoSub.Flags().Bool("stats", false, "get cluster stats")

	// infoSub.RegisterFlagCompletionFunc("target", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// 	levels := []string{"_all", "http", "ingest", "thread_pool", "script"}
	// 	return levels, cobra.ShellCompDirectiveNoFileComp
	// })

	return infoSub
}

func runInfoClusterCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		// target, _ := cmd.Flags().GetString("target")
		settings, _ := cmd.Flags().GetBool("settings")
		info, _ := cmd.Flags().GetBool("info")
		stats, _ := cmd.Flags().GetBool("stats")
		r, err := es.Cluster.Info("_all").Do(ctx)

		if err != nil {
			return fmt.Errorf("at getting cluster info %w", err)
		}

		if err = json.NewEncoder(cmd.OutOrStdout()).Encode(r); err != nil {
			return serde.SerializingError(err)
		}

		if !settings {
			return nil
		}

		s, err := es.Cluster.GetSettings().Do(ctx)
		if err != nil {
			return fmt.Errorf("at getting cluster settings %w", err)
		}

		if err = json.NewEncoder(cmd.OutOrStdout()).Encode(s); err != nil {
			return serde.SerializingError(err)
		}

		return nil
	}

}
