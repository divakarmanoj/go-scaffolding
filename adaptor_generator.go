package main

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/format"
	"os"
)

func GenerateAdaptor(s *Structure) string {
	f := NewFile("main")
	GenerateRequestToModel(s, f)
	GenerateModelToResponse(s, f)
	fmt.Printf("%#v", f)
	output := fmt.Sprintf("%#v", f)
	outputBytes, err := format.Source([]byte(output))
	if err != nil {
		fmt.Println(output)
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	// write rawOutput to file
	err = os.WriteFile(toSnakeCase(s.Name)+"/adaptor.go", outputBytes, 0644)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	return output
}

func GenerateRequestToModel(s *Structure, f *File) {
	f.Func().Id("RequestTo" + toCamelCase(s.Name)).Params(Id("request").Id("*" + s.Name + "Request")).Id("*" + toCamelCase(s.Name) + "Model").BlockFunc(func(g *Group) {
		g.If(Id("request").Op("==").Nil()).Block(
			Return(Nil()),
		)
		g.Return(Op("&").Id(toCamelCase(s.Name) + "Model").ValuesFunc(func(g *Group) {

			for _, attr := range s.Attributes {
				if !isValidType(attr.Type) {
					fmt.Println("Error: invalid type " + attr.Type)
					os.Exit(1)
				}
				if attr.Type != "struct" {
					if !attr.IsRequired {
						g.Id(toCamelCase(attr.Name)).Op(":").Qual("github.com/divakarmanoj/go-scaffolding/imports", "Null"+toTitleCase(attr.Type)+"Ptr").Call(Id("request").Dot(toCamelCase(attr.Name)))
					} else {
						g.Id(toCamelCase(attr.Name)).Op(":").Id("request").Dot(toCamelCase(attr.Name))
					}
				} else {
					GenerateRequestToModel(&Structure{Name: attr.Name, Attributes: attr.Attributes}, f)
					g.Id(toCamelCase(attr.Name)).Op(":").Id("RequestTo" + toCamelCase(attr.Name)).Call(Id("request").Dot(toCamelCase(attr.Name)))
				}
			}
		}))
	})
}

func GenerateModelToResponse(s *Structure, f *File) {
	f.Func().Id("ModelTo" + toCamelCase(s.Name)).Params(Id("model").Id("*" + s.Name + "Model")).Id("*" + toCamelCase(s.Name) + "Response").BlockFunc(func(g *Group) {
		g.If(Id("model").Op("==").Nil()).Block(
			Return(Nil()),
		)
		g.Return(Op("&").Id(toCamelCase(s.Name) + "Response").ValuesFunc(func(g *Group) {
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
						g.Id(toCamelCase(attr.Name)).Op(":").Qual("github.com/divakarmanoj/go-scaffolding/imports", "Null"+toTitleCase(attr.Type)+"ToPtr").Call(Id("model").Dot(toCamelCase(attr.Name)))
					} else {
						g.Id(toCamelCase(attr.Name)).Op(":").Id("model").Dot(toCamelCase(attr.Name))
					}
				} else {
					GenerateModelToResponse(&Structure{Name: attr.Name, Attributes: attr.Attributes}, f)
					g.Id(toCamelCase(attr.Name)).Op(":").Id("ModelTo" + toCamelCase(attr.Name)).Call(Id("model").Dot(toCamelCase(attr.Name)))
				}
			}
		}))
	})
}
