package main

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/format"
	"os"
)

func GenerateHandler(s *Structure) {
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
	err = os.WriteFile(toSnakeCase(s.Name)+"/handler.go", outputBytes, 0644)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
}

func remove(parsedStruct *Structure, f *File) *Statement {
	return f.Func().Id("Delete"+parsedStruct.Name).Params(Id("w").Qual("net/http", "ResponseWriter"), Id("r").Op("*").Qual("net/http", "Request")).BlockFunc(func(mainBlock *Group) {
		mainBlock.List(Id("ids"), Id("ok")).Op(":=").Id("r").Dot("URL").Dot("Query").Call().Index(Lit("id"))
		mainBlock.If(Op("!").Id("ok").Op("||").Len(Id("ids").Index(Lit(0))).Op("<").Lit(1)).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Lit("id is required"), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Id("id").Op(":=").Id("ids").Index(Lit(0))
		mainBlock.If(Id("err").Op(":=").Id("db").Dot("Delete").Call(Op("&").Id(toCamelCase(parsedStruct.Name)+"Model{}"), Id("id")).Dot("Error").Id(";").Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
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

func update(parsedStruct *Structure, f *File) *Statement {
	return f.Func().Id("Update"+parsedStruct.Name).Params(Id("w").Qual("net/http", "ResponseWriter"), Id("r").Op("*").Qual("net/http", "Request")).BlockFunc(func(mainBlock *Group) {
		mainBlock.List(Id("ids"), Id("ok")).Op(":=").Id("r").Dot("URL").Dot("Query").Call().Index(Lit("id"))
		mainBlock.If(Op("!").Id("ok").Op("||").Len(Id("ids").Index(Lit(0))).Op("<").Lit(1)).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Lit("id is required"), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Id("id").Op(":=").Id("ids").Index(Lit(0))
		mainBlock.Line()
		mainBlock.Var().Id(toCamelCase(parsedStruct.Name)).Id(toCamelCase(parsedStruct.Name) + "Request")
		mainBlock.If(Id("err").Op(":=").Qual("encoding/json", "NewDecoder").Call(Id("r").Dot("Body")).Dot("Decode").Call(Op("&").Id(toCamelCase(parsedStruct.Name))).Op(";").Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Id("err").Dot("Error").Call(), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Line()
		mainBlock.Id("model").Op(":=").Id("RequestTo" + toCamelCase(parsedStruct.Name)).Call(Op("&").Id(toCamelCase(parsedStruct.Name)))
		mainBlock.If(Id("err").Op(":=").Id("db").Dot("Model").Call(Op("&").Id("model")).Dot("Where").Call(Id("\"id = ?\""), Id("id")).Dot("Updates").Call(Op("&").Id("model")).Dot("Error").Id(";").Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Id("err").Dot("Error").Call(), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Line()
		mainBlock.Var().Id("output").Op("=").Qual("github.com/divakarmanoj/go-scaffolding/imports", "Response").Values(Dict{
			Id("Data"):    Id("ModelTo" + toCamelCase(parsedStruct.Name)).Call(Id("model")),
			Id("Status"):  Id("\"success\""),
			Id("Message"): Lit(parsedStruct.Name + " updated successfully"),
		})
		mainBlock.Qual("encoding/json", "NewEncoder").Call(Id("w")).Dot("Encode").Call(Id("output"))
	})
}

func read(parsedStruct *Structure, f *File) *Statement {
	return f.Func().Id("Read"+parsedStruct.Name).Params(Id("w").Qual("net/http", "ResponseWriter"), Id("r").Op("*").Qual("net/http", "Request")).BlockFunc(func(mainBlock *Group) {
		mainBlock.Var().Id("err").Error()
		mainBlock.List(Id("pageNumbers"), Id("ok")).Op(":=").Id("r").Dot("URL").Dot("Query").Call().Index(Lit("page_number"))
		mainBlock.Id("pageNumber").Op(":=").Lit(1)
		mainBlock.If(Id("ok").Op("&&").Len(Id("pageNumbers").Index(Lit(0))).Op(">").Lit(1)).BlockFunc(func(ifBlock *Group) {
			ifBlock.List(Id("pageNumber"), Id("err")).Op("=").Qual("strconv", "Atoi").Call(Id("pageNumbers").Index(Lit(0)))
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
		mainBlock.Var().Id(parsedStruct.Name).Index().Id(parsedStruct.Name + "Model")
		mainBlock.Id("err").Op("=").Id("db").Dot("Limit").Call(Id("pageSize")).Dot("Offset").Call(Id("pageSize").Op("*").Parens(Id("pageNumber").Op("-").Lit(1))).Dot("Preload").Call(Qual("gorm.io/gorm/clause", "Associations")).Dot("Find").Call(Op("&").Id(parsedStruct.Name)).Dot("Error")
		mainBlock.If(Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Id("err").Dot("Error").Call(), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Line()
		mainBlock.Var().Id("data").Index().Op("*").Id(toCamelCase(parsedStruct.Name) + "Response")
		mainBlock.For(List(Id("_"), Id("i")).Op(":=").Range().Id(toCamelCase(parsedStruct.Name))).BlockFunc(func(forBlock *Group) {
			forBlock.Id("data").Op("=").Append(Id("data"), Id("ModelTo"+toCamelCase(parsedStruct.Name)).Call(Op("&").Id("i")))
		})
		mainBlock.Var().Id("output").Op("=").Qual("github.com/divakarmanoj/go-scaffolding/imports", "Response").Values(Dict{
			Id("Data"):    Id("data"),
			Id("Status"):  Id("\"success\""),
			Id("Message"): Lit(parsedStruct.Name + "s retrieved successfully"),
		})
		mainBlock.Qual("encoding/json", "NewEncoder").Call(Id("w")).Dot("Encode").Call(Id("output"))
	})
}

func create(parsedStruct *Structure, f *File) *Statement {
	return f.Func().Id("Create"+parsedStruct.Name).Params(Id("w").Qual("net/http", "ResponseWriter"), Id("r").Op("*").Qual("net/http", "Request")).BlockFunc(func(mainBlock *Group) {
		mainBlock.Var().Id(toCamelCase(parsedStruct.Name)).Id(parsedStruct.Name + "Request")
		//mainBlock.Id("err := json.NewDecoder(r.Body).Decode(&" + toCamelCase(parsedStruct.Name) + ")")
		mainBlock.Id("err").Op(":=").Qual("encoding/json", "NewDecoder").Call(Id("r").Dot("Body")).Dot("Decode").Call(Op("&").Id(toCamelCase(parsedStruct.Name)))
		mainBlock.If(Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Id("err").Dot("Error").Call(), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Line()
		mainBlock.Id("model").Op(":=").Id("RequestTo" + toCamelCase(parsedStruct.Name)).Call(Op("&").Id(toCamelCase(parsedStruct.Name)))
		mainBlock.If(Id("err").Op(":=").Id("db").Dot("Create").Call(Id("model")).Dot("Error").Op(";").Id("err").Op("!=").Nil()).BlockFunc(func(ifBlock *Group) {
			ifBlock.Qual("net/http", "Error").Call(Id("w"), Id("err").Dot("Error").Call(), Qual("net/http", "StatusBadRequest"))
			ifBlock.Return()
		})
		mainBlock.Line()
		mainBlock.Var().Id("output").Op("=").Qual("github.com/divakarmanoj/go-scaffolding/imports", "Response").Values(Dict{
			Id("Data"):    Id("ModelTo" + toCamelCase(parsedStruct.Name)).Call(Id("model")),
			Id("Status"):  Id("\"success\""),
			Id("Message"): Lit(parsedStruct.Name + " created successfully"),
		})
		mainBlock.Qual("encoding/json", "NewEncoder").Call(Id("w")).Dot("Encode").Call(Id("output"))
	})
}
