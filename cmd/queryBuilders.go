package cmd

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
)

// For searching based on elasticsearch id use this one, since we have the build in IdsQuery
func BuildTermIdQuery(ids []string) *types.Query {

	switch len(ids) {
	case 0:
		return nil
	default:

		r := types.NewIdsQuery()
		r.Values = ids

		q := &types.Query{Ids: r}
		return q

	}

}

// build a simple term/terms query for a field. PAss a field, and the value(s)
func BuildTermLevelQuery(field string, values []string) *types.Query {

	switch len(values) {
	case 0:
		return nil

	case 1:

		query := &types.Query{
			Term: map[string]types.TermQuery{
				field: {Value: values[0]},
			},
		}

		return query

	default:
		m := make(map[string]types.TermsQueryField)
		m[field] = values

		query := &types.Query{
			Terms: &types.TermsQuery{TermsQuery: m},
		}

		return query
	}
}
