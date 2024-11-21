package gitter

import (
	"context"
	"fmt"

	"github.com/machinebox/graphql"
)

type GraphQLClient struct {
	Url   string
	Token string
}

func AuthorizedGraphQLQuery[ResponseType interface{}](client GraphQLClient, query string, response *ResponseType) error {
	req := graphql.NewRequest(query)

	gql := graphql.NewClient(client.Url)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", client.Token))

	err := gql.Run(context.Background(), req, &response)
	if err != nil {
		return fmt.Errorf("error fetching PRs: %w", err)
	}

	return nil
}
