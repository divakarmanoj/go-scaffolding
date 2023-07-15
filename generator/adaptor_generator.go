package generator

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/format"
	"os"
)

func GenerateAdaptor(s *Config, outputDir string) string {
	f := NewFile("main")
	GenerateRequestToModel(s, f)
	GenerateModelToResponse(s, f)
	output := fmt.Sprintf("%#v", f)
	outputBytes, err := format.Source([]byte(output))
	if err != nil {
		fmt.Println(output)
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	// write rawOutput to file
	err = os.WriteFile(outputDir+ToSnakeCase(s.Name)+"/adaptor.go", outputBytes, 0644)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	return output
}

func GenerateRequestToModel(s *Config, f *File) {
	f.Func().Id("RequestTo" + s.camelCase).Params(Id("request").Id("*" + s.camelCase + "Request")).Id("*" + s.camelCase + "Model").BlockFunc(func(g *Group) {
		g.If(Id("request").Op("==").Nil()).Block(
			Return(Nil()),
		)
		g.Return(Op("&").Id(s.camelCase + "Model").ValuesFunc(func(g *Group) {

			for _, attr := range s.Attributes {
				if !isValidType(attr.Type) {
					fmt.Println("Error: invalid type " + attr.Type)
					os.Exit(1)
				}
				if attr.Type != "struct" {
					if !attr.IsRequired {
						g.Id(attr.camelCase).Op(":").Qual("github.com/divakarmanoj/go-scaffolding/imports", "Null"+toTitleCase(attr.Type)+"Ptr").Call(Id("request").Dot(attr.camelCase))
					} else {
						g.Id(attr.camelCase).Op(":").Id("request").Dot(attr.camelCase)
					}
				} else {
					GenerateRequestToModel(&Config{Name: attr.Name, camelCase: attr.camelCase, Attributes: attr.Attributes}, f)
					g.Id(attr.camelCase).Op(":").Id("RequestTo" + attr.camelCase).Call(Id("request").Dot(attr.camelCase))
				}
			}
		}))
	})
}

func GenerateModelToResponse(s *Config, f *File) {
	f.Func().Id("ModelTo" + s.camelCase).Params(Id("model").Id("*" + s.camelCase + "Model")).Id("*" + s.camelCase + "Response").BlockFunc(func(g *Group) {
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
					g.Id(attr.camelCase).Op(":").Id("ModelTo" + attr.camelCase).Call(Id("model").Dot(attr.camelCase))
				}
			}
		}))
	})
}
