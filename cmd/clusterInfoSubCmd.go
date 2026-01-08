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

func addInfoClusterFlags(infoSub InfoSubCmd) InfoSubCmd {

	infoSub.Flags().StringP("target", "t", "_all", "target for cluster information")

	infoSub.RegisterFlagCompletionFunc("target", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		levels := []string{"_all", "http", "ingest", "thread_pool", "script"}
		return levels, cobra.ShellCompDirectiveNoFileComp
	})

	return infoSub
}

func runInfoClusterSubcmd(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		target, _ := cmd.Flags().GetString("target")
		r, err := es.Cluster.Info(target).Do(ctx)

		if err != nil {
			return fmt.Errorf("at getting cluster info %w", err)
		}

		if err = json.NewEncoder(cmd.OutOrStdout()).Encode(r); err != nil {
			return serde.SerializingError(err)
		}

		return nil
	}

}
