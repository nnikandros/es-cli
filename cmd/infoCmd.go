package cmd

import (
	"context"
	"encoding/json"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/info"
	"github.com/nnikandros/serde"
	"github.com/spf13/cobra"
)

type InfoCmd = *cobra.Command

type ClusterCmd = *cobra.Command
type PingSubCmd = *cobra.Command
type InfoSubCmd = *cobra.Command
type NodesSubCmd = *cobra.Command
type IndexSubCmd = *cobra.Command

func infoCmdFunc(es *elasticsearch.TypedClient) InfoCmd {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "you know, for search",
		Long:  `Basic information of cluster, including version. You know, for search`,
		RunE:  runInfoCmdFunc(es),
		Args:  cobra.NoArgs,
	}

	return cmd
}

func runInfoCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		r, err := info.New(es).Do(ctx)
		if err != nil {
			return err
		}

		if err := json.NewEncoder(cmd.OutOrStdout()).Encode(r); err != nil {
			return serde.SerializingError(err)
		}

		return nil
	}

}

func InfoCmdFunc(es *elasticsearch.TypedClient) InfoCmd {
	infoCmd := infoCmdFunc(es)
	clusterSubCmd := addInfoClusterFlags(infoclusterCmdFunc(es))
	nodesSubCmd := addInfoNodesFlags(infoNodesCmdFunc(es))
	indexSubCmd := addInfoIndexFlags(infoIndexCmdFunc(es))

	infoCmd.AddCommand(clusterSubCmd)
	infoCmd.AddCommand(nodesSubCmd)
	infoCmd.AddCommand(indexSubCmd)

	return infoCmd
}
