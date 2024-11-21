package repo

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
	templatePath := "queries/repo/templates/prMultiRepo.gql"

	tpl, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("cannot read template path (%v): %w", templatePath, err)
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

	fileOutputPath := "./byRepoQuery.gql"

	outputFile, err := os.Create(fileOutputPath)
	if err != nil {
		return "", fmt.Errorf("cannot create file (%v): %w", fileOutputPath, err)
	}

	// write string to file
	_, err = outputFile.WriteString(savedPRQuery)
	if err != nil {
		return "", fmt.Errorf("cannot write to file (%v): %w", fileOutputPath, err)
	}

	return savedPRQuery, nil
}
