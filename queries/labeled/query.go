package labeled

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"strings"

	"github.com/ettle/strcase"
)

var funcMap = template.FuncMap{
	"ToSnake": strcase.ToSnake,
	"ToKebab": strcase.ToKebab,
	"ToLower": strings.ToLower,
	"ToUpper": strings.ToUpper,
}

//go:embed templates/prMultiRepoByLabel.gql
var tpl string

func generateQuery(conf Config) (string, error) {
	loadedTemplate, err := template.New("prMultiRepoByLabel").Funcs(funcMap).Parse(string(tpl))
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
