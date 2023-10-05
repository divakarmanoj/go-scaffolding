package generator

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/format"
	"os"
)

func GenerateModel(s *Config, outputDir string) (string, []string) {
	f := NewFile("main")
	f.ImportName("github.com/divakarmanoj/go-scaffolding/imports", "imports")
	modelNames := GormStruct(s, f)
	GenerateModelToResponse(s, f)
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

func GormStruct(s *Config, f *File) []string {
	modelNames := []string{s.camelCase + "Model"}
	f.Type().Id(s.camelCase + "Model").StructFunc(func(g *Group) {
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

func GenerateModelToResponse(s *Config, f *File) {
	f.Func().Params(
		Id("model").Id("*" + s.camelCase + "Model"),
	).Id("ToResponse").Params().Id("*" + s.camelCase + "Response").BlockFunc(func(g *Group) {
		g.If(Id("model").Op("==").Nil()).Block(
			Return(Nil()),
		)
		g.Return(Op("&").Id(s.camelCase + "Response").ValuesFunc(func(g *Group) {
			g.Id("ID").Op(":").Id("model").Dot("Model").Dot("ID")
			g.Id("CreatedAt").Op(":").Id("model").Dot("Model").Dot("CreatedAt")
			g.Id("UpdatedAt").Op(":").Id("model").Dot("Model").Dot("UpdatedAt")
			for _, attr := range s.Attributes {
				if !isValidType(attr.Type) {
					fmt.Println("Error: invalid type " + attr.Type)
					os.Exit(1)
				}
				if attr.Type != "struct" {
					if !attr.IsRequired {
						g.Id(attr.camelCase).Op(":").Qual("github.com/divakarmanoj/go-scaffolding/imports", "Null"+toTitleCase(attr.Type)+"ToPtr").Call(Id("model").Dot(attr.camelCase))
					} else {
						g.Id(attr.camelCase).Op(":").Id("model").Dot(attr.camelCase)
					}
				} else {
					GenerateModelToResponse(&Config{Name: attr.Name, camelCase: attr.camelCase, Attributes: attr.Attributes}, f)
					g.Id(attr.camelCase).Op(":").Id("model").Dot(attr.camelCase).Dot("ToResponse").Call()
				}
			}
		}))
	})
}
