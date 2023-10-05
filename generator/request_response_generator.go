package generator

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/format"
	"os"
)

func GenerateRequestResponse(s *Config, outputDir string) {
	f := NewFile("main")
	GenerateResponseStruct(s, f)
	GenerateRequestStruct(s, f)
	GenerateRequestToModel(s, f)
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

func GenerateRequestStruct(s *Config, f *File) {
	f.Type().Id(s.camelCase + "Request").StructFunc(func(g *Group) {
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

func GenerateResponseStruct(s *Config, f *File) {
	f.Type().Id(s.camelCase + "Response").StructFunc(func(g *Group) {
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

func GenerateRequestToModel(s *Config, f *File) {
	f.Func().Params(Id("request").Id("*" + s.camelCase + "Request")).Id("ToModel").Params().Id("*" + s.camelCase + "Model").BlockFunc(func(g *Group) {
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
					g.Id(attr.camelCase).Op(":").Id("request").Dot(attr.camelCase).Dot("ToModel").Call()
				}
			}
		}))
	})
}
