package by_author

import (
	_ "embed"

	"appletini/queries/templating"
)

//go:embed templates/query.gql
var tpl string

func generateQuery(conf Config) (string, error) {
	return templating.Render("author_query", tpl, conf)
}
