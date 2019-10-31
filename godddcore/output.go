package godddcore

import (
	"goddd/strmangle"
	"os"
	"text/template"
)

var templateFunctions = template.FuncMap{
	"titleCase": strmangle.TitleCase,
	"camelCase": strmangle.CamelCase,
}

func generateOutput(name string) error {
	templateFile := "sample.go.tql"
	templateFilePath := "templates/sample.go.tql"

	tql, err := template.New(templateFile).Funcs(templateFunctions).ParseFiles(templateFilePath)
	if err != nil {
		return err
	}

	data := templateData{
		Name: name,
	}

	return tql.Execute(os.Stdout, data)
}
