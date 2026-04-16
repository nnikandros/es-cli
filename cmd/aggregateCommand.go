package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nnikandros/serde"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/spf13/cobra"
)

type AggregateCmd = *cobra.Command

type Bucket struct {
	DocCount int    `json:"doc_count"`
	Key      string `json:"key"`
}

type Aggregations struct {
	Buckets                 []Bucket `json:"buckets"`
	DocCountErrorUpperBound int      `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int      `json:"sum_other_doc_count"`
}

func aggregateCmdFunc(es *elasticsearch.TypedClient) AggregateCmd {
	cmd := &cobra.Command{
		Use:   "aggregate",
		Short: "aggegate functionality",
		Long:  `A longer description that spans multiple lines and likely contains`,

		RunE:              runAggregateCmdFunc(es),
		ValidArgsFunction: ValidIndexArgsAutoCompletion(es),
	}

	return cmd

}

func addAggregateFlags(aggsCmd AggregateCmd) AggregateCmd {
	aggsCmd.Flags().StringP("field", "f", "", "add desc")
	aggsCmd.Flags().BoolP("buckets", "b", false, "add desc")

	return aggsCmd

}

func runAggregateCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {
		ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)
		field, _ := cmd.Flags().GetString("field")
		buckets, _ := cmd.Flags().GetBool("buckets")

		index := ParseArgsIntoString(cmd, args)
		defer cancelFunc()

		aggs := map[string]types.Aggregations{field: {Terms: &types.TermsAggregation{Field: &field}}}
		r, err := es.Search().Index(index).Aggregations(aggs).Size(0).Do(ctx)
		if err != nil {
			return fmt.Errorf("at Search with aggregatations %w", err)
		}

		if buckets {
			if err = json.NewEncoder(cmd.OutOrStdout()).Encode(r.Aggregations); err != nil {
				return fmt.Errorf("at serializing aggrs %w", err)
			}

			return nil
		}
		a := Aggregations{}
		fieldAggs, ok := r.Aggregations[field]
		if !ok {
			return fmt.Errorf("r.Aggegations is missing key %v", field)
		}

		b, err := json.Marshal(fieldAggs)
		if err != nil {
			return serde.SerializingError(err)
		}

		if err = json.Unmarshal(b, &a); err != nil {
			return serde.DeserializingError(err)
		}

		var uniqueKeys []string
		for _, k := range a.Buckets {
			uniqueKeys = append(uniqueKeys, k.Key)
		}

		for _, value := range uniqueKeys {
			fmt.Fprintf(cmd.OutOrStdout(), "%v\n", value)
		}
		return nil
	}
}

func processAggregateResponse(r *search.Response) {

}

func createAggsMap(cmd *cobra.Command) map[string]types.Aggregations {
	return nil

}
