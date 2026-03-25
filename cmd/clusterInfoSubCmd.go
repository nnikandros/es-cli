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
	infoSub.Flags().Bool("settings", false, "Get cluster-wide settings. By default, it returns only settings that have been explicitly defined.")
	infoSub.Flags().Bool("stats", false, "Get cluster statistics. Get basic index metrics (shard numbers, store size, memory usage) and information about the current nodes that form the cluster (number, roles, os, jvm versions, memory usage, cpu and installed plugins")
	infoSub.Flags().Bool("ping", false, "pings cluster")

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

		settings, _ := cmd.Flags().GetBool("settings")
		stats, _ := cmd.Flags().GetBool("stats")
		ping, _ := cmd.Flags().GetBool("ping")

		switch {
		case settings:
			s, err := es.Cluster.GetSettings().IncludeDefaults(true).Do(ctx)
			if err != nil {
				return fmt.Errorf("at getting cluster settings %w", err)
			}

			if err = json.NewEncoder(cmd.OutOrStdout()).Encode(s); err != nil {
				return serde.SerializingError(err)
			}

			return nil

		case stats:
			s, err := es.Cluster.Stats().Do(ctx)
			if err != nil {
				return fmt.Errorf("at getting cluster stats %w", err)
			}

			if err = json.NewEncoder(cmd.OutOrStdout()).Encode(s); err != nil {
				return serde.SerializingError(err)
			}

			return nil
		case ping:
			h, err := es.Ping().Do(ctx)
			if err != nil {
				return fmt.Errorf("pinging the cluster %w", err)
			}

			if h {
				fmt.Fprintln(cmd.OutOrStdout(), h)
			} else {
				return fmt.Errorf("ping returned false. Check connection to the cluster, credentials, hosts etc")
			}

			return nil

		default:
			r, err := es.Cluster.Info("_all").Do(ctx)
			if err != nil {
				return fmt.Errorf("at getting cluster info %w", err)
			}

			if err = json.NewEncoder(cmd.OutOrStdout()).Encode(r); err != nil {
				return serde.SerializingError(err)
			}

			return nil

		}

	}

}
