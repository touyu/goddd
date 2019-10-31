package godddcore

import (
"fmt"
"goddd/strmangle"
"os"
"text/template"
)

var templateFunctions = template.FuncMap{
	"titleCase": strmangle.TitleCase,
	"camelCase": strmangle.CamelCase,
}

type templateSet struct {
	FileName string
	OutputDir string
}

var templateSets = []*templateSet{
	{
		FileName: "sample.go.tql",
		OutputDir: "application",
	},
}

func generateOutput(name string) error {
	templateSet := templateSets[0]
	templateFileName := templateSet.FileName
	templateFilePath := fmt.Sprintf("templates/%s", templateFileName)

	tql, err := template.New(templateFileName).Funcs(templateFunctions).ParseFiles(templateFilePath)
	if err != nil {
		return err
	}

	data := templateData{
		Name: name,
	}

	if err := os.Mkdir(templateSet.OutputDir, 0700); err != nil {
		return err
	}

	outputFilePath := fmt.Sprintf("%s/%s.go", templateSet.OutputDir, name)
	outputfile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}

	return tql.Execute(outputfile, data)
}
