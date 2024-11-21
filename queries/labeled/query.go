package labeled

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/ettle/strcase"
)

var funcMap = template.FuncMap{
	"ToSnake": strcase.ToSnake,
	"ToKebab": strcase.ToKebab,
	"ToLower": strings.ToLower,
	"ToUpper": strings.ToUpper,
}

func generateQuery(conf Config) (string, error) {
	templatePath := "queries/labeled/templates/prMultiRepoByLabel.gql"

	tpl, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("cannot read template file: %w", err)
	}

	loadedTemplate, err := template.New(templatePath).Funcs(funcMap).Parse(string(tpl))
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

	fileOutputPath := "queries/labeled/finalQuery.gql"

	outputFile, err := os.Create(fileOutputPath)
	if err != nil {
		return "", fmt.Errorf("cannot create file: %w", err)
	}

	// write string to file
	_, err = outputFile.WriteString(savedPRQuery)
	if err != nil {
		return "", fmt.Errorf("cannot write to file: %w", err)
	}

	return savedPRQuery, nil
}
