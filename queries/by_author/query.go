package by_author

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
)

//go:embed templates/query.gql
var tpl string

func generateQuery(conf Config) (string, error) {
	loadedTemplate, err := template.New("author_query").Parse(string(tpl))
	if err != nil {
		return "", fmt.Errorf("cannot load template funcmap: %w", err)
	}

	prQuery := new(bytes.Buffer)

	err = loadedTemplate.Execute(prQuery, conf)
	if err != nil {
		return "", fmt.Errorf("cannot create query: %w", err)
	}

	// query created from config
	savedPRQuery := fmt.Sprint(prQuery)

	return savedPRQuery, nil
}
