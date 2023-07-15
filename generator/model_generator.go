package generator

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"go/format"
	"os"
)

func GenerateModel(s *Config, outputDir string) (string, []string) {
	f := jen.NewFile("main")
	f.ImportName("github.com/divakarmanoj/go-scaffolding/imports", "imports")
	modelNames := GormStruct(s, f)
	output := fmt.Sprintf("%#v", f)
	outputBytes, err := format.Source([]byte(output))
	if err != nil {
		fmt.Println(output)
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	// write rawOutput to file
	err = os.WriteFile(outputDir+ToSnakeCase(s.Name)+"/models.go", outputBytes, 0644)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	return output, modelNames
}

func GormStruct(s *Config, f *jen.File) []string {
	modelNames := []string{s.camelCase + "Model"}
	f.Type().Id(s.camelCase + "Model").StructFunc(func(g *jen.Group) {
		g.Qual("github.com/divakarmanoj/go-scaffolding/imports", "Model")
		for _, attr := range s.Attributes {
			if !isValidType(attr.Type) {
				fmt.Printf("Error: Invalid type %s\n", attr.Type)
				os.Exit(1)
			}
			if attr.Type != "struct" {
				if attr.IsRequired {
					g.Id(attr.camelCase).Id(attr.Type).Tag(map[string]string{"json": ToSnakeCase(attr.Name)})
				} else {
					g.Id(attr.camelCase).Qual("database/sql", "Null"+toTitleCase(attr.Type)).Tag(map[string]string{"json": ToSnakeCase(attr.Name)})
				}
			} else {
				f.Line()
				NestedModelNames := GormStruct(&Config{Name: attr.Name, camelCase: attr.camelCase, Attributes: attr.Attributes}, f)
				modelNames = append(modelNames, NestedModelNames...)
				g.Id(attr.camelCase + "ID").Id("uint").Tag(map[string]string{"json": ToSnakeCase(attr.Name) + "_id"})
				g.Id(attr.camelCase).Id("*" + attr.camelCase + "Model").Tag(map[string]string{"json": ToSnakeCase(attr.Name)})
			}
		}
	})
	f.Line()
	return modelNames
}
