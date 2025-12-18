package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"serde"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/spf13/cobra"
)

type SearchCmd = *cobra.Command

func searchCmdFunc(es *elasticsearch.TypedClient) SearchCmd {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "search API",
		Long:  "Running a search query against an index\nThe arguments can be an index a list of indexes separated by space or index name with wildcard. Can also be _all to search all indices",
		RunE:  runSearchCmdFunc(es),
		Example: `es search <test-index-*>
es search <index-1 index-2 >
es search <index-1,index-2> --size 100
es search <index-1 index-2> -s 100
`,
		ValidArgsFunction: ValidArgsFuncAutoCompletion(es),
	}

	return cmd

}

func addSearchFlags(searchCmd SearchCmd) SearchCmd {

	searchCmd.Flags().IntP("size", "s", 10, "size of search")
	searchCmd.Flags().StringSliceP("fields", "f", []string{}, "source  fields to return")
	searchCmd.Flags().BoolP("time", "t", false, "sort by time, newest first")
	// searchCmd.Flags().BoolP("reverse", "r", false, "reverse order when sorting. So enabling will show oldest first")
	searchCmd.Flags().Bool("tab", false, "display the output of --fields in a table format")

	searchCmd.Flags().Bool("terms", false, "do a term/terms search.")

	// Fields to do a term/terms search against an index
	searchCmd.Flags().StringSlice("id", []string{}, "do a term/terms search based on elasticsearch internal _id. If you provide one id it will be a term search. If you provide more than one, it will be a terms search")
	searchCmd.Flags().StringSlice("LEVEL", []string{}, "do a term search for a LEVEL")

	// searchCmd.RegisterFlagCompletionFunc("LEVEL", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// 	levels := []string{"DEBUG", "ERROR", "INFO"}
	// 	return levels, cobra.ShellCompDirectiveNoFileComp
	// })
	searchCmd.RegisterFlagCompletionFunc("LEVEL", cobra.FixedCompletions([]string{"DEBUG", "ERROR", "INFO"}, cobra.ShellCompDirectiveNoFileComp))

	searchCmd.Flags().StringSlice("APP_NAME", []string{}, "do a term search for an APP_NAME")
	return searchCmd

}

func runSearchCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		ParsedFlags, err := ParsedFlagsFromCmd(cmd)
		if err != nil {
			return fmt.Errorf("at parsing flags %w", err)
		}
		indexName := ParseArgsIntoSingleString(args)

		r, err := searchWithFlags(es, indexName, ParsedFlags)
		if err != nil {
			return fmt.Errorf("at doing search with flags, %w", err)
		}

		if err = processResponse(r, ParsedFlags, cmd.OutOrStdout()); err != nil {
			return fmt.Errorf("at processing the response %w", err)
		}

		// fmt.Fprintf(cmd.OutOrStdout(), "%s\n", b)

		return nil
	}
}

func searchWithFlags(es *elasticsearch.TypedClient, indexName string, flags SearchFlags) (*search.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	r, err := buildQuery(es, indexName, flags).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("at doing search with buildQuery, %w", err)
	}

	return r, nil
}

func buildQuery(es *elasticsearch.TypedClient, indexName string, flags SearchFlags) *search.Search {
	sortMap := make(map[string]string)
	if flags.Time {

		sortMap["TIMESTAMP"] = "asc"
	}

	searchReq := es.Search().Index(indexName).Size(flags.Size).Sort(sortMap)

	if flags.Terms && flags.Id != nil {
		if q := BuildTermIdQuery(flags.Id); q != nil {
			searchReq = searchReq.Query(q)
		}
	}

	if flags.Terms && flags.LEVEL != nil {
		if q := BuildTermLevelQuery("LEVEL", flags.LEVEL); q != nil {
			searchReq = searchReq.Query(q)
		}
	}

	if flags.Terms && flags.APP_NAME != nil {
		if q := BuildTermLevelQuery("APP_NAME", flags.APP_NAME); q != nil {
			searchReq = searchReq.Query(q)
		}
	}

	if len(flags.Fields) > 0 {

		fields := make([]types.FieldAndFormat, 0, len(flags.Fields))
		for _, field := range flags.Fields {
			fields = append(fields, types.FieldAndFormat{Field: field})
		}
		searchReq = searchReq.Fields(fields...).Source_(false)
	}

	return searchReq
}

func processResponse(r *search.Response, flags SearchFlags, w io.Writer) error {
	if len(flags.Fields) > 0 {
		results := make([]map[string]json.RawMessage, 0, flags.Size)

		for _, hit := range r.Hits.Hits {
			results = append(results, hit.Fields)
		}

		if len(results) == 0 {
			return fmt.Errorf("no results to process")
		}

		if !flags.Tabular {
			b, err := json.MarshalIndent(results, "", " ")
			if err != nil {
				return serde.SerializingError(err)
			}

			fmt.Fprintf(w, "%s", b)
			return nil
		}

		if err := processTab(results, w); err != nil {
			return fmt.Errorf("at post-processing for table format  %w", err)
		}

		return nil
	}

	b, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return serde.SerializingError(err)
	}

	fmt.Fprintf(w, "%s", b)

	return nil

}

func processTab(results []map[string]json.RawMessage, w io.Writer) error {

	NumerOfESFilesRetrieved := len(results[0])

	unmarshaledResults := unmarshalValues(results)

	retrievedEsFields := KeysSorted(unmarshaledResults[0])

	tabWrite := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tabWrite, strings.Join(retrievedEsFields, "\t"))

	for _, result := range unmarshaledResults {
		row := make([]string, 0, NumerOfESFilesRetrieved)
		for _, k := range retrievedEsFields {
			row = append(row, strings.Join(result[k], ""))
		}
		fmt.Fprintln(tabWrite, strings.Join(row, "\t"))
	}

	tabWrite.Flush()

	return nil

}

func unmarshalValues(results []map[string]json.RawMessage) []map[string][]string {

	records := make([]map[string][]string, 0, len(results))

	for _, m := range results {
		rec := make(map[string][]string)
		for k, v := range m {
			var arr []string
			if err := json.Unmarshal(v, &arr); err != nil {
				rec[k] = []string{serde.DeserializingError(err).Error()}
			}

			rec[k] = arr
		}
		records = append(records, rec)
	}

	return records

}

func KeysSorted(m map[string][]string) []string {

	keys := slices.Collect(maps.Keys(m))
	slices.Sort(keys)

	return keys

}
