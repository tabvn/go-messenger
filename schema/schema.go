package schema

import (
	"github.com/graphql-go/graphql"
	"messenger/query"
	"messenger/mutation"
	"fmt"
	"context"
)

var Schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    query.Query,
		Mutation: mutation.Mutation,
	},
)

func ExecuteQuery(context context.Context, query string, operation string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
		OperationName: operation,
		Context:       context,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}
