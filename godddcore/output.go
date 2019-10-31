package godddcore

import (
	"fmt"
	"goddd/strmangle"
	"log"
	"os"
	"strings"
	"text/template"
)

var templateFunctions = template.FuncMap{
	"titleCase": strmangle.TitleCase,
	"camelCase": strmangle.CamelCase,
}

type templateSet struct {
	FileName string
	OutputDir string
	OutputFileName string
}

var templateSets = []*templateSet{
	{
		FileName: "application.go.tql",
		OutputDir: "application",
	},
	{
		FileName: "service_interface.go.tql",
		OutputDir: "domain/service/creator",
		OutputFileName: "interface.go",
	},
}

func generateOutput(name string) error {
	for _, set := range templateSets {
		err := executeTemplate(set, name)
		if err != nil {
			return err
		}
	}
	return nil
}

func executeTemplate(templateSet *templateSet, name string) error {
	templateFileName := templateSet.FileName
	templateFilePath := fmt.Sprintf("templates/%s", templateFileName)

	tql, err := template.New(templateFileName).Funcs(templateFunctions).ParseFiles(templateFilePath)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(templateSet.OutputDir, 0700); err != nil {
		return err
	}

	outputFileName := name + ".go"
	if templateSet.OutputFileName != "" {
		outputFileName = templateSet.OutputFileName
	}

	outputFilePath := fmt.Sprintf("%s/%s", templateSet.OutputDir, outputFileName)
	outputfile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}

	data := templateData{
		Name: name,
		CurrentDir: getWorkingDirName(),
	}

	return tql.Execute(outputfile, data)
}

func getWorkingDirName() string {
	p, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	paths := strings.Split(p, "/")
	return paths[len(paths)-1]
}
