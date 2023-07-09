package main

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	"go/format"
	"os"
)

func GenerateMain(parsedStruct *Structure, modelnames []string) {
	f := NewFile("main")

	f.Var().Id("db").Op("*").Qual("gorm.io/gorm", "DB")

	// Generate the main function
	f.Func().Id("main").Params().BlockFunc(func(mainBlock *Group) {
		mainBlock.Var().Id("err").Error()
		mainBlock.Line()

		// Open the database connection
		mainBlock.List(Id("db"), Id("err")).Op("=").Qual("gorm.io/gorm", "Open").Call(
			Qual("gorm.io/driver/sqlite", "Open").Call(
				Lit(toCamelCase(parsedStruct.Name)+".db")),
			Op("&").Qual("gorm.io/gorm", "Config").Values(),
		)
		mainBlock.If(Id("err").Op("!=").Nil()).Block(
			Panic(Lit("failed to connect database")),
		)
		mainBlock.Defer().Func().Params().Block(
			Id("err").Op("=").Qual("os", "RemoveAll").Call(
				Lit(toCamelCase(parsedStruct.Name)+".db"),
			),
			If(Id("err").Op("!=").Nil()).Block(
				Return(),
			),
		).Call()
		mainBlock.Line()

		// Auto-migrate models
		for _, model := range modelnames {
			mainBlock.Id("db").Dot("AutoMigrate").Call(Op("&").Id(model + "{}"))
		}
		mainBlock.Line()

		// Handle HTTP routes
		mainBlock.Qual("net/http", "HandleFunc").Call(
			Lit("/"+toSnakeCase(parsedStruct.Name)+"/read"),
			Id("Read"+toCamelCase(parsedStruct.Name)),
		)
		mainBlock.Qual("net/http", "HandleFunc").Call(
			Lit("/"+toSnakeCase(parsedStruct.Name)+"/create"),
			Id("Create"+toCamelCase(parsedStruct.Name)),
		)
		mainBlock.Qual("net/http", "HandleFunc").Call(
			Lit("/"+toSnakeCase(parsedStruct.Name)+"/update"),
			Id("Update"+toCamelCase(parsedStruct.Name)),
		)
		mainBlock.Qual("net/http", "HandleFunc").Call(
			Lit("/"+toSnakeCase(parsedStruct.Name)+"/delete"),
			Id("Delete"+toCamelCase(parsedStruct.Name)),
		)
		mainBlock.Err().Op("=").Qual("net/http", "ListenAndServe").Call(
			Lit(":3333"),
			Nil(),
		)
		mainBlock.If(Err().Op("!=").Nil()).Block(
			Panic(Err()),
		)
	})
	output := fmt.Sprintf("%#v", f)
	outputBytes, err := format.Source([]byte(output))
	if err != nil {
		fmt.Println(output)
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	// write rawOutput to file
	err = os.WriteFile(toSnakeCase(parsedStruct.Name)+"/main.go", outputBytes, 0644)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
}
