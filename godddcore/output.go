package godddcore

import (
	"os"
	"text/template"
)

func generateOutput(name string) error {
	tql, err := template.ParseFiles("templates/sample.go.tql")
	if err != nil {
		return err
	}

	data := templateData{
		Name: name,
	}

	return tql.Execute(os.Stdout, data)
}
