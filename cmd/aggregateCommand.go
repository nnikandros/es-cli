package cmd

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

type AggregateCmd = *cobra.Command

func aggregateCmdFunc(es *elasticsearch.TypedClient) AggregateCmd {
	cmd := &cobra.Command{
		Use:   "tasks",
		Short: "tasks",
		Long:  `A longer description that spans multiple lines and likely contains`,

		// RunE: runAggregateCmdFunc(es),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		// SilenceUsage:      true,
		// ValidArgsFunction: ValidArgsFuncAutoCompletion(es),

	}

	return cmd

}

func addAggregateFlags(listIndices AggregateCmd) AggregateCmd {
	listIndices.Flags().StringSliceP("field", "f", []string{}, "add desc")
	return listIndices

}

// func runAggregateCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
// 	return func(cmd *cobra.Command, args []string) error {
// 		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
// 		field, _ := cmd.Flags().GetStringSlice("field")

// 		index := ParseArgsIntoString(cmd, args)
// 		defer cancelFunc()

// 	}
// }
