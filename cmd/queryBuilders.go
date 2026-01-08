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

func BuildTermLevelQueryV2(field string, values []string) *types.Query {

	switch len(values) {
	case 0:
		return nil

	default:
		m := make(map[string]types.TermsQueryField)
		m[field] = values

		query := &types.Query{
			Terms: &types.TermsQuery{TermsQuery: m},
		}

		return query
	}
}

func BuildFilterQuery(field string, values []string) *types.Query {
	q := &types.Query{}
	b := &types.BoolQuery{}

	m := make(map[string]types.TermsQueryField)
	m[field] = values

	query := types.Query{
		Terms: &types.TermsQuery{TermsQuery: m},
	}

	f := []types.Query{query}

	b.Filter = f
	q.Bool = b
	return q

}

func MustQuery(q1 *types.Query, q2 *types.Query) *types.Query {
	q := &types.Query{}
	b := &types.BoolQuery{}

	m := []types.Query{*q1, *q2}
	b.Must = m

	q.Bool = b

	return q

}

func ShouldQuery(q1 *types.Query, q2 *types.Query) *types.Query {
	q := &types.Query{}
	b := &types.BoolQuery{}

	s := []types.Query{*q1, *q2}
	b.Should = s

	q.Bool = b

	return q

}
