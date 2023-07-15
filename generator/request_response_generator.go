package generator

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"go/format"
	"os"
)

func GenerateRequestResponse(s *Config, outputDir string) {
	f := jen.NewFile("main")
	GenerateResponseStruct(s, f)
	GenerateRequestStruct(s, f)
	output := fmt.Sprintf("%#v", f)
	outputBytes, err := format.Source([]byte(output))
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	// write rawOutput to file
	err = os.WriteFile(outputDir+ToSnakeCase(s.Name)+"/requestResponse.go", outputBytes, 0644)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
}

func GenerateRequestStruct(s *Config, f *jen.File) {
	f.Type().Id(s.camelCase + "Request").StructFunc(func(g *jen.Group) {
		for _, attr := range s.Attributes {
			if !isValidType(attr.Type) {
				fmt.Printf("Error: Invalid type %s\n", attr.Type)
				os.Exit(1)
			}
			if attr.Type != "struct" {
				if attr.IsRequired {
					g.Id(attr.camelCase).Id(attr.Type).Tag(map[string]string{"json": ToSnakeCase(attr.Name)})
				} else {
					g.Id(attr.camelCase).Id("*" + attr.Type).Tag(map[string]string{"json": ToSnakeCase(attr.Name) + ",omitempty"})
				}
			} else {
				f.Line()
				GenerateRequestStruct(&Config{Name: attr.Name, camelCase: attr.camelCase, Attributes: attr.Attributes}, f)
				g.Id(attr.camelCase).Id("*" + attr.camelCase + "Request").Tag(map[string]string{"json": ToSnakeCase(attr.Name)})
			}
		}
	})
	f.Line()
}

func GenerateResponseStruct(s *Config, f *jen.File) {
	f.Type().Id(s.camelCase + "Response").StructFunc(func(g *jen.Group) {
		g.Id("ID").Id("uint").Tag(map[string]string{"json": "id"})
		g.Id("CreatedAt").Id("int64").Tag(map[string]string{"json": "created_at"})
		g.Id("UpdatedAt").Id("int64").Tag(map[string]string{"json": "updated_at"})
		for _, attr := range s.Attributes {
			if !isValidType(attr.Type) {
				fmt.Printf("Error: Invalid type %s\n", attr.Type)
				os.Exit(1)
			}
			if attr.Type != "struct" {
				if attr.IsRequired {
					g.Id(attr.camelCase).Id(attr.Type).Tag(map[string]string{"json": ToSnakeCase(attr.Name)})
				} else {
					g.Id(attr.camelCase).Id("*" + attr.Type).Tag(map[string]string{"json": ToSnakeCase(attr.Name) + ",omitempty"})
				}
			} else {
				f.Line()
				GenerateResponseStruct(&Config{Name: attr.Name, camelCase: attr.camelCase, Attributes: attr.Attributes}, f)
				g.Id(attr.camelCase).Id("*" + attr.camelCase + "Response").Tag(map[string]string{"json": ToSnakeCase(attr.Name)})
			}
		}
	})
	f.Line()
}
