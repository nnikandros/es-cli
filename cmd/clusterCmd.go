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

type ClusterCmd = *cobra.Command

type PingSubCmd = *cobra.Command
type InfoSubCmd = *cobra.Command
type NodesSubCmd = *cobra.Command

func clusterCmdFunc(es *elasticsearch.TypedClient) ClusterCmd {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "actions about the cluster",
		Long:  `stuff abouit the cluster`,
		RunE:  runCusterCmdFunc(es),
		Args:  cobra.NoArgs,
	}

	return cmd
}

func ClusterCmdFunc(es *elasticsearch.TypedClient) ClusterCmd {
	clusterCmd := clusterCmdFunc(es)

	pingSubCmd := pingClusterCmdFunc(es)
	infoSubCmd := infoClusterCmdFunc(es)

	nodesSubCmd := nodesClusterCmdFunc(es)

	clusterCmd.AddCommand(pingSubCmd)
	clusterCmd.AddCommand(infoSubCmd)
	clusterCmd.AddCommand(nodesSubCmd)

	return clusterCmd
}

func pingClusterCmdFunc(es *elasticsearch.TypedClient) PingSubCmd {
	cmd := &cobra.Command{
		Use:   "ping",
		Short: "pings",
		Long:  `ping ping`,
		RunE:  runPingClusterCmdFunc(es),
		Args:  cobra.NoArgs,
	}

	return cmd

}

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

func nodesClusterCmdFunc(es *elasticsearch.TypedClient) NodesSubCmd {
	cmd := &cobra.Command{
		Use:   "nodes",
		Short: "nodes",
		Long:  `nodes nodes`,
		RunE:  runNodesSubCmd(es),
		Args:  cobra.NoArgs,
	}

	return cmd

}

func runPingClusterCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		h, err := es.Ping().Do(ctx)
		if err != nil {
			return fmt.Errorf("pinging the cluster %w", err)
		}

		if h {
			fmt.Println(h)
		} else {
			return fmt.Errorf("ping returned false. Check connection to the cluster, credentials, hosts etc")
		}

		return nil
	}
}

func runCusterCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}

}

func runInfoClusterSubcmd(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		r, err := es.Cluster.Info("_all").Do(ctx)

		if err != nil {
			return fmt.Errorf("at getting cluster info %w", err)
		}

		b, err := json.Marshal(r)
		if err != nil {
			return serde.SerializingError(err)
		}

		fmt.Fprintf(cmd.OutOrStdout(), "%s", b)

		return nil
	}

}

func runNodesSubCmd(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		infoResponse, err := es.Nodes.Info().Do(ctx)
		if err != nil {
			return fmt.Errorf("at getting nodes info %w", err)
		}

		b, err := json.Marshal(infoResponse)
		if err != nil {
			return serde.SerializingError(err)
		}
		fmt.Fprintf(cmd.OutOrStdout(), "%s", b)

		return nil
	}

}
