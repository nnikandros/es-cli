package cmd

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

type InfoCmd = *cobra.Command

type ClusterCmd = *cobra.Command
type PingSubCmd = *cobra.Command
type InfoSubCmd = *cobra.Command
type NodesSubCmd = *cobra.Command
type IndexSubCmd = *cobra.Command

// func clusterCmdFunc(es *elasticsearch.TypedClient) ClusterCmd {
// 	cmd := &cobra.Command{
// 		Use:   "cluster",
// 		Short: "actions about the cluster",
// 		Long:  `stuff abouit the cluster`,
// 		RunE:  runClusterCmdFunc(es),
// 		Args:  cobra.NoArgs,
// 	}

// 	return cmd
// }

func infoCmdFunc(es *elasticsearch.TypedClient) InfoCmd {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "info about the es",
		Long:  `stuff abouit the cluster`,
		RunE:  runInfoCmdFunc(es),
		Args:  cobra.NoArgs,
	}

	return cmd
}

func runInfoCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}

}

func InfoCmdFunc(es *elasticsearch.TypedClient) InfoCmd {
	infoCmd := infoCmdFunc(es)
	clusterSubCmd := addInfoClusterFlags(infoclusterCmdFunc(es))
	nodesSubCmd := addInfoNodesFlags(infoNodesCmdFunc(es))

	// pingSubCmd := pingClusterCmdFunc(es)
	// infoSubCmd := addInfoClusterFlags(infoClusterCmdFunc(es))

	// nodesSubCmd := nodesClusterCmdFunc(es)

	// clusterCmd.AddCommand(pingSubCmd)
	// clusterCmd.AddCommand(infoSubCmd)
	// clusterCmd.AddCommand(nodesSubCmd)
	infoCmd.AddCommand(clusterSubCmd)
	infoCmd.AddCommand(nodesSubCmd)

	return infoCmd
}

// func pingClusterCmdFunc(es *elasticsearch.TypedClient) PingSubCmd {
// 	cmd := &cobra.Command{
// 		Use:   "ping",
// 		Short: "pings",
// 		Long:  `ping ping`,
// 		RunE:  runPingClusterCmdFunc(es),
// 		Args:  cobra.NoArgs,
// 	}

// 	return cmd

// }

// func nodesClusterCmdFunc(es *elasticsearch.TypedClient) NodesSubCmd {
// 	cmd := &cobra.Command{
// 		Use:   "nodes",
// 		Short: "nodes",
// 		Long:  `nodes nodes`,
// 		RunE:  runNodesSubCmd(es),
// 		Args:  cobra.NoArgs,
// 	}

// 	return cmd

// }

// func runPingClusterCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
// 	return func(cmd *cobra.Command, args []string) error {
// 		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 		defer cancel()
// 		h, err := es.Ping().Do(ctx)
// 		if err != nil {
// 			return fmt.Errorf("pinging the cluster %w", err)
// 		}

// 		if h {
// 			fmt.Println(h)
// 		} else {
// 			return fmt.Errorf("ping returned false. Check connection to the cluster, credentials, hosts etc")
// 		}

// 		return nil
// 	}
// }

// func runClusterCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
// 	return func(cmd *cobra.Command, args []string) error {
// 		return nil
// 	}

// }
