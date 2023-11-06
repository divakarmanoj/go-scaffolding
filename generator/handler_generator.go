package generator

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/format"
	"os"
)

func GenerateHandler(s *Config, outputDir string) {
	f := NewFile("main")
	f.ImportName("github.com/divakarmanoj/go-scaffolding/imports", "imports")
	f.ImportName("gorm.io/gorm/clause", "clause")
	create(s, f)
	f.Line()
	read(s, f)
	f.Line()
	update(s, f)
	f.Line()
	remove(s, f)
	output := fmt.Sprintf("%#v", f)
	outputBytes, err := format.Source([]byte(output))
	if err != nil {
		fmt.Println(output)
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	// write rawOutput to file
	err = os.WriteFile(outputDir+ToSnakeCase(s.Name)+"/handler.go", outputBytes, 0644)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
}

func remove(parsedStruct *Config, f *File) *Statement {
	return f.Func().Id("Delete"+parsedStruct.camelCase).Params(Id("w").Qual("net/http", "ResponseWriter"), Id("r").Op("*").Qual("net/http", "Request")).BlockFunc(func(mainBlock *Group) {
		mainBlock.List(Id("ids"), Id("ok")).Op(":=").Id("r").Dot("URL").Dot("Query").Call().Index(Lit("id"))
		mainBlock.If(Op("!").Id("ok").Op("||").Len(Id("ids").Index(Lit(0))).Op("<").Lit(1)).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Lit("id is required"), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Id("id").Op(":=").Id("ids").Index(Lit(0))
		mainBlock.If(Id("err").Op(":=").Id("db").Dot("Delete").Call(Op("&").Id(parsedStruct.camelCase+"Model{}"), Id("id")).Dot("Error").Id(";").Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Id("err").Dot("Error").Call(), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Line()
		mainBlock.Var().Id("output").Op("=").Qual("github.com/divakarmanoj/go-scaffolding/imports", "Response").Values(Dict{
			Id("Status"):  Id("\"success\""),
			Id("Message"): Lit(parsedStruct.Name + " deleted successfully"),
		})
		mainBlock.Qual("encoding/json", "NewEncoder").Call(Id("w")).Dot("Encode").Call(Id("output"))
	})
}

func update(parsedStruct *Config, f *File) *Statement {
	return f.Func().Id("Update"+parsedStruct.camelCase).Params(Id("w").Qual("net/http", "ResponseWriter"), Id("r").Op("*").Qual("net/http", "Request")).BlockFunc(func(mainBlock *Group) {
		mainBlock.List(Id("ids"), Id("ok")).Op(":=").Id("r").Dot("URL").Dot("Query").Call().Index(Lit("id"))
		mainBlock.If(Op("!").Id("ok").Op("||").Len(Id("ids").Index(Lit(0))).Op("<").Lit(1)).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Lit("id is required"), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Id("id").Op(":=").Id("ids").Index(Lit(0))
		mainBlock.Line()
		mainBlock.Var().Id(parsedStruct.camelCase).Id(parsedStruct.camelCase + "Request")
		mainBlock.If(Id("err").Op(":=").Qual("encoding/json", "NewDecoder").Call(Id("r").Dot("Body")).Dot("Decode").Call(Op("&").Id(parsedStruct.camelCase)).Op(";").Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Id("err").Dot("Error").Call(), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Line()
		mainBlock.Id("model").Op(":=").Id(parsedStruct.camelCase).Dot("ToModel").Call()
		mainBlock.If(Id("err").Op(":=").Id("db").Dot("Model").Call(Op("&").Id("model")).Dot("Where").Call(Id("\"id = ?\""), Id("id")).Dot("Updates").Call(Op("&").Id("model")).Dot("Error").Id(";").Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Id("err").Dot("Error").Call(), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Line()
		mainBlock.Var().Id("output").Op("=").Qual("github.com/divakarmanoj/go-scaffolding/imports", "Response").Values(Dict{
			Id("Status"):  Id("\"success\""),
			Id("Message"): Lit(parsedStruct.Name + " updated successfully"),
		})
		mainBlock.Qual("encoding/json", "NewEncoder").Call(Id("w")).Dot("Encode").Call(Id("output"))
	})
}

func read(parsedStruct *Config, f *File) *Statement {
	return f.Func().Id("Read"+parsedStruct.camelCase).Params(Id("w").Qual("net/http", "ResponseWriter"), Id("r").Op("*").Qual("net/http", "Request")).BlockFunc(func(mainBlock *Group) {
		mainBlock.Var().Id("err").Error()
		mainBlock.List(Id("cursors"), Id("ok")).Op(":=").Id("r").Dot("URL").Dot("Query").Call().Index(Lit("cursor"))
		mainBlock.Id("cursor").Op(":=").Lit(1)
		mainBlock.If(Id("ok").Op("&&").Len(Id("cursors").Index(Lit(0))).Op(">").Lit(1)).BlockFunc(func(ifBlock *Group) {
			ifBlock.List(Id("cursor"), Id("err")).Op("=").Qual("strconv", "Atoi").Call(Id("cursors").Index(Lit(0)))
			ifBlock.If(Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
				ifBlock.Qual("net/http", "Error").Call(Id("w"), Id("err").Dot("Error").Call(), Qual("net/http", "StatusBadRequest"))
				ifBlock.Return()
			})
		})
		mainBlock.Line()
		mainBlock.List(Id("pageSizes"), Id("ok")).Op(":=").Id("r").Dot("URL").Dot("Query").Call().Index(Lit("page_size"))
		mainBlock.Id("pageSize").Op(":=").Lit(10)
		mainBlock.If(Id("ok").Op("&&").Len(Id("pageSizes").Index(Lit(0))).Op(">").Lit(1)).BlockFunc(func(ifBlock *Group) {
			ifBlock.List(Id("pageSize"), Id("err")).Op("=").Qual("strconv", "Atoi").Call(Id("pageSizes").Index(Lit(0)))
			ifBlock.If(Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
				ifBlock.Qual("net/http", "Error").Call(Id("w"), Id("err").Dot("Error").Call(), Qual("net/http", "StatusBadRequest"))
				ifBlock.Return()
			})
		})
		mainBlock.Line()
		mainBlock.Var().Id(parsedStruct.camelCase).Index().Id(parsedStruct.camelCase + "Model")
		mainBlock.Id("err").Op("=").Id("db").Dot("Preload").Call(Qual("gorm.io/gorm/clause", "Associations")).Dot("Limit").Call(Id("pageSize").Op("+").Lit(1)).Dot("Where").Call(Id("\"id >= ?\""), Id("cursor")).Dot("Order").Call(Id("\"id asc\"")).Dot("Find").Call(Op("&").Id(parsedStruct.camelCase)).Dot("Error")
		mainBlock.If(Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Id("err").Dot("Error").Call(), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Line()
		mainBlock.Var().Id("data").Index().Op("*").Id(parsedStruct.camelCase + "Response")
		mainBlock.For(List(Id("_"), Id("i")).Op(":=").Range().Id(parsedStruct.camelCase)).BlockFunc(func(forBlock *Group) {
			forBlock.Id("data").Op("=").Append(Id("data"), Id("i").Dot("ToResponse").Call())
			forBlock.If(Len(Id("data")).Op("==").Id("pageSize")).BlockFunc(func(ifBlock *Group) {
				ifBlock.Break()
			})
		})
		mainBlock.Var().Id("output").Op("=").Qual("github.com/divakarmanoj/go-scaffolding/imports", "Response").Values(Dict{
			Id("Data"):    Id("data"),
			Id("Status"):  Id("\"success\""),
			Id("Message"): Lit(parsedStruct.Name + "s retrieved successfully"),
		})
		mainBlock.If(Len(Id(parsedStruct.camelCase)).Op(">").Id("pageSize")).BlockFunc(func(ifBlock *Group) {
			ifBlock.Id("output").Dot("Cursor").Op("=").Id(parsedStruct.camelCase).Index(Id("pageSize")).Dot("ID")
		})
		mainBlock.Qual("encoding/json", "NewEncoder").Call(Id("w")).Dot("Encode").Call(Id("output"))
	})
}

func create(parsedStruct *Config, f *File) *Statement {
	return f.Func().Id("Create"+parsedStruct.camelCase).Params(Id("w").Qual("net/http", "ResponseWriter"), Id("r").Op("*").Qual("net/http", "Request")).BlockFunc(func(mainBlock *Group) {
		mainBlock.Var().Id(parsedStruct.camelCase).Id(parsedStruct.camelCase + "Request")
		mainBlock.Id("err").Op(":=").Qual("encoding/json", "NewDecoder").Call(Id("r").Dot("Body")).Dot("Decode").Call(Op("&").Id(parsedStruct.camelCase))
		mainBlock.If(Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Id("err").Dot("Error").Call(), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Line()
		mainBlock.Id("model").Op(":=").Id(parsedStruct.camelCase).Dot("ToModel").Call()
		mainBlock.If(Id("err").Op(":=").Id("db").Dot("Create").Call(Id("model")).Dot("Error").Op(";").Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Id("err").Dot("Error").Call(), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Line()
		mainBlock.Var().Id("output").Op("=").Qual("github.com/divakarmanoj/go-scaffolding/imports", "Response").Values(Dict{
			Id("Data"):    Id("model").Dot("ToResponse").Call(),
			Id("Status"):  Id("\"success\""),
			Id("Message"): Lit(parsedStruct.Name + " created successfully"),
		})
		mainBlock.Qual("encoding/json", "NewEncoder").Call(Id("w")).Dot("Encode").Call(Id("output"))
	})
}
