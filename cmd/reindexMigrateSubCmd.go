package cmd

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

func reindexMigrateCmdFunc(es *elasticsearch.TypedClient) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reindex",
		Short: "reindex API",
		Long:  "To be updated",
		// RunE:              runCloneIndexCmdFunc(es),
		// ValidArgsFunction: ValidArgsFuncAutoCompletion(es),
		Example: `es clone -s <clone-index-1> -s <source-index-2>  -t <target-index>`,
	}

	return cmd

}

func addReindexFlags(countCmd *cobra.Command) *cobra.Command {
	countCmd.Flags().StringP("target", "t", "", "Name of clone index that will be created.")
	countCmd.Flags().StringSliceP("source", "s", []string{}, "Name of clone index that will be created.")
	return countCmd

}

func runReIndexCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		return nil
	}

}
