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
		OutputDir: "domain/service/{{name}}",
		OutputFileName: "interface.go",
	},
}

func generateOutput(name string) error {
	currentDir := getWorkingDirName()

	fmt.Println("processing...\n")

	for _, set := range templateSets {
		err := executeTemplate(set, name, currentDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func executeTemplate(templateSet *templateSet, name string, currentDir string) error {
	templateFileName := templateSet.FileName
	templateFilePath := fmt.Sprintf("templates/%s", templateFileName)

	tql, err := template.New(templateFileName).Funcs(templateFunctions).ParseFiles(templateFilePath)
	if err != nil {
		return err
	}

	outputDir := strings.Replace(templateSet.OutputDir, "{{name}}", name, 1)

	if err := os.MkdirAll(outputDir, 0700); err != nil {
		return err
	}

	outputFileName := name + ".go"
	if templateSet.OutputFileName != "" {
		outputFileName = templateSet.OutputFileName
	}

	outputFilePath := fmt.Sprintf("%s/%s", outputDir, outputFileName)
	outputfile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}

	data := templateData{
		Name: name,
		CurrentDir: currentDir,
	}

	if err := tql.Execute(outputfile, data); err != nil {
		return err
	}

	fmt.Printf("generated: %s\n", outputFilePath)
	return nil
}

func getWorkingDirName() string {
	p, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	paths := strings.Split(p, "/")
	return paths[len(paths)-1]
}
