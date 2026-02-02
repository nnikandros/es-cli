package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"maps"
	"serde"
	"slices"
	"strings"
	"text/tabwriter"
	"time"

	_ "embed"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"github.com/spf13/cobra"
	"go.yaml.in/yaml/v2"
)

type SearchCmd = *cobra.Command

//go:embed es_fields.yaml
var fileByte []byte

type EsFieldsConfig struct {
	Name         string   `yaml:"name"`
	DefaultValue []string `yaml:"default"`
	ValidArgs    []string `yaml:"valid-args"`
	Usage        string   `yaml:"usage"`
	Value        []string
	Date         bool `yaml:"date"`
}

type EsFields struct {
	Fields []EsFieldsConfig `yaml:"fields"`
}

var e EsFields

func searchCmdFunc(es *elasticsearch.TypedClient) SearchCmd {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "search API",
		Long:  "Running a search query against an index\nThe arguments can be an index a list of indexes separated by space or index name with wildcard. Can also be _all to search all indices",
		RunE:  runSearchCmdFunc(es),
		Example: `es search <test-index-*>
es search <index-1 index-2 > (Do an empty search against the provided two indices)
es search <index-1,index-2> --size 100 (Do an empty search against the provided two indices but increase the size of the results to 100)
es search <index-1 index-2> -s 100 (Do an empty search against the provided two indices but increase the size of the results to 100 but in shorthand notation)
es search <index> --sort-by TIMESTAMP (Do an empty search but sort the results with newsest on top based on the the field TIMESTAMP. If u want autocompletion you have to add it on the es_fields.yaml)
`,
		ValidArgsFunction: ValidArgsFuncAutoCompletion(es),
	}

	return cmd

}

func addSearchFlags(searchCmd SearchCmd) SearchCmd {
	var sortByArgs []string

	searchCmd.Flags().IntP("size", "s", 10, "size of search")
	searchCmd.Flags().StringSliceP("fields", "f", []string{}, "source  fields to return")
	searchCmd.Flags().StringSlice("sort-by", []string{}, "sort by given date field, newest first")
	searchCmd.Flags().BoolP("reverse", "r", false, "reverse the order of results, i.e. show newest in the bottom")

	searchCmd.Flags().Bool("tab", false, "display the output of --fields in a table format")

	searchCmd.Flags().Bool("filter", false, "make a filter query")

	searchCmd.Flags().StringSlice("id", []string{}, "do a search based on elasticsearch internal _id. If you provide one id it will be a term search. If you provide more than one, it will be a terms search")

	err := yaml.Unmarshal(fileByte, &e)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range e.Fields {
		if !f.Date {
			searchCmd.Flags().StringSlice(f.Name, f.DefaultValue, f.Usage)
			searchCmd.RegisterFlagCompletionFunc(f.Name, cobra.FixedCompletions(f.ValidArgs, cobra.ShellCompDirectiveNoFileComp))
		} else {
			sortByArgs = append(sortByArgs, f.Name)

		}

	}

	searchCmd.RegisterFlagCompletionFunc("sort-by", cobra.FixedCompletions(sortByArgs, cobra.ShellCompDirectiveNoFileComp))

	return searchCmd

}

func runSearchCmdFunc(es *elasticsearch.TypedClient) RunEFunc {
	return func(cmd *cobra.Command, args []string) error {

		ParsedFlags, err := ParsedFlagsFromCmd(cmd)
		if err != nil {
			return fmt.Errorf("at parsing flags %w", err)
		}

		indexName := ParseArgsIntoString(cmd, args)

		r, err := searchWithFlags(es, indexName, ParsedFlags)
		if err != nil {
			return fmt.Errorf("at doing search with flags, %w", err)
		}

		if err = processResponse(r, ParsedFlags, cmd.OutOrStdout()); err != nil {
			return fmt.Errorf("at processing the response %w", err)
		}

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

// query builder
func buildQuery(es *elasticsearch.TypedClient, indexName string, flags SearchFlags) *search.Search {
	sort := []types.SortCombinations{}
	if len(flags.SortBy) > 0 {

		for _, datetimeField := range flags.SortBy {
			sort = append(sort, types.SortOptions{
				SortOptions: map[string]types.FieldSort{
					datetimeField: {
						Order: &sortorder.Desc,
					},
				},
			})
		}

	}
	searchReq := es.Search().Index(indexName).Size(flags.Size).Sort(sort...)

	if flags.Id != nil {
		if q := BuildIdQuery(flags.Id); q != nil {
			searchReq = searchReq.Query(q)
		}
	}

	var queries []types.Query

	if flags.Filter {

		for _, i := range flags.FieldsTermsMap {
			if q := BuildTermsQuery(i.Name, i.Value); q != nil {
				queries = append(queries, *q)
			}

		}

		q := BuildFilterAndQuery(queries)
		searchReq = searchReq.Query(q)

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

	hits := make([]types.Hit, 0, len(r.Hits.Hits))

	switch flags.Reverse {
	case true:
		hits = Reverse(r.Hits.Hits)
	default:
		hits = append(hits, r.Hits.Hits...)

	}

	// if flags.Reverse {

	// 	reversedHits := Reverse(r.Hits.Hits)

	// 	hits = reversedHits
	// } else {
	// 	hits = r.Hits.Hits
	// }

	if len(flags.Fields) > 0 {
		results := make([]map[string]json.RawMessage, 0, flags.Size)

		for _, hit := range hits {
			results = append(results, hit.Fields)
		}

		if len(results) == 0 {
			return fmt.Errorf("no results to process")
		}

		if !flags.Tabular {
			if err := json.NewEncoder(w).Encode(results); err != nil {
				return serde.SerializingError(err)
			}

			return nil

		}

		if err := processTab(results, w); err != nil {
			return fmt.Errorf("at post-processing for table format  %w", err)
		}

		return nil
	}

	if err := json.NewEncoder(w).Encode(r); err != nil {
		return serde.SerializingError(err)
	}

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
