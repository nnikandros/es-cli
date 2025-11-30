package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/spf13/cobra"
)

func listIndicesAsSlice(es *elasticsearch.TypedClient) ([]cobra.Completion, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	s := make([]cobra.Completion, 0, 35)

	defer cancelFunc()

	r, err := es.Cat.Indices().Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("at getting the indices %w", err)
	}

	for _, indexRecord := range r {
		if IsIndexNameValid(*indexRecord.Index) {
			s = append(s, *indexRecord.Index)
		}
	}

	return s, nil

}

func ValidArgsFuncAutoCompletion(es *elasticsearch.TypedClient) cobra.CompletionFunc {
	s := func(cmd *cobra.Command, args []string, toComplete string) ([]cobra.Completion, cobra.ShellCompDirective) {
		s, err := listIndicesAsSlice(es)
		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return s, cobra.ShellCompDirectiveNoFileComp

	}

	return s

}
