package main

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"go/format"
	"os"
)

func GenerateModel(s *Structure) (string, []string) {
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
	err = os.WriteFile(toSnakeCase(s.Name)+"/models.go", outputBytes, 0644)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	return output, modelNames
}

func GormStruct(s *Structure, f *jen.File) []string {
	modelNames := []string{s.Name + "Model"}
	f.Type().Id(s.Name + "Model").StructFunc(func(g *jen.Group) {
		g.Qual("github.com/divakarmanoj/go-scaffolding/imports", "Model")
		for _, attr := range s.Attributes {
			if !isValidType(attr.Type) {
				fmt.Printf("Error: Invalid type %s\n", attr.Type)
				os.Exit(1)
			}
			if attr.Type != "struct" {
				if attr.IsRequired {
					g.Id(toCamelCase(attr.Name)).Id(attr.Type).Tag(map[string]string{"json": toSnakeCase(attr.Name)})
				} else {
					g.Id(toCamelCase(attr.Name)).Qual("database/sql", "Null"+toTitleCase(attr.Type)).Tag(map[string]string{"json": toSnakeCase(attr.Name)})
				}
			} else {
				f.Line()
				NestedModelNames := GormStruct(&Structure{Name: attr.Name, Attributes: attr.Attributes}, f)
				modelNames = append(modelNames, NestedModelNames...)
				g.Id(toCamelCase(attr.Name) + "ID").Id("uint").Tag(map[string]string{"json": toSnakeCase(attr.Name) + "_id"})
				g.Id(toCamelCase(attr.Name)).Id("*" + attr.Name + "Model").Tag(map[string]string{"json": toSnakeCase(attr.Name)})
			}
		}
	})
	f.Line()
	return modelNames
}
