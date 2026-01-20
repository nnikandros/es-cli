package cmd

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

func infoIndexCmdFunc(es *elasticsearch.TypedClient) ClusterCmd {
	cmd := &cobra.Command{
		Use:   "index",
		Short: "information about indices",
		Long:  `indicies`,
		// RunE:  runInfoClusterCmdFunc(es),
		Args: cobra.NoArgs,
	}

	return cmd
}
