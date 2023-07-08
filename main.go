package main

import (
	"fmt"
	"go/format"
	"os"
)

type Structure struct {
	Name       string       `json:"name"`
	Attributes []Attributes `json:"attributes"`
}

var input = `{
	"name": "Super",
	"Attributes": [{
			"name": "Name",
			"type": "string",
			"is_required": false
		},
		{
			"name": "Age",
			"type": "int16",
			"is_required": true
		},
		{
			"name": "Address",
			"type": "struct",
			"attributes": [{
					"Name": "Street Name",
					"Type": "string",
					"is_required": true
				},
				{
					"Name": "City",
					"Type": "string",
					"is_required": true
				},
				{
					"Name": "State",
					"Type": "string",
					"is_required": false
				},
				{
					"Name": "Zip",
					"Type": "int16",
					"is_required": true
				}
			]
		}
	]
}`

func main() {
	parsedStruct, err := ParseStruct(input)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}

	// Create directory
	err = os.Mkdir(toSnakeCase(parsedStruct.Name), 0755)

	GenerateRequestResponse(parsedStruct)
	_, modelnames := GenerateModel(parsedStruct)
	GenerateAdaptor(parsedStruct)
	GenerateHandler(parsedStruct)
	GenerateMain(parsedStruct, modelnames)
}

func GenerateMain(parsedStruct *Structure, modelnames []string) {
	output := "package main\n\n"
	output += "import (\n"
	output += "\t\"github.com/divakarmanoj/go-scaffolding/imports\"\n"
	output += "\t\"github.com/gorilla/mux\"\n"
	output += ")\n\n"
	output += "func main() {\n"
	output += "\tr := mux.NewRouter()\n"
	output += "\timports.InitDB()\n"

	output = "package main\n\n" +
		"import (\n" +
		"\t\"gorm.io/driver/sqlite\"\n" +
		"\t\"gorm.io/gorm\"\n" +
		"\t\"net/http\"\n" +
		"\t\"os\"\n" +
		")\n" +
		"\n" +
		"func main() {\n" +
		"\tvar err error\n" +
		fmt.Sprintf("\tdb, err = gorm.Open(sqlite.Open(\"%s.db\"), &gorm.Config{})\n", toCamelCase(parsedStruct.Name)) +
		"\tif err != nil {\n" +
		"\t\tpanic(\"failed to connect database\")\n" +
		"\t}\n" +
		"\tdefer func() {\n" +
		fmt.Sprintf("\t\terr = os.RemoveAll(\"%s.db\")\n", toCamelCase(parsedStruct.Name)) +
		"\t\tif err != nil {\n" +
		"\t\t\treturn\n\t\t}\n" +
		"\t}()\n"
	for _, model := range modelnames {
		output += fmt.Sprintf("\tdb.AutoMigrate(&%s{})\n", model)
	}
	output += "\n" +
		fmt.Sprintf("\thttp.HandleFunc(\"/%s/read\", Read%s)\n", toSnakeCase(parsedStruct.Name), toCamelCase(parsedStruct.Name)) +
		fmt.Sprintf("\thttp.HandleFunc(\"/%s/create\", Create%s)\n", toSnakeCase(parsedStruct.Name), toCamelCase(parsedStruct.Name)) +
		fmt.Sprintf("\thttp.HandleFunc(\"/%s/update\", Update%s)\n", toSnakeCase(parsedStruct.Name), toCamelCase(parsedStruct.Name)) +
		fmt.Sprintf("\thttp.HandleFunc(\"/%s/delete\", Delete%s)\n", toSnakeCase(parsedStruct.Name), toCamelCase(parsedStruct.Name)) +
		"\terr = http.ListenAndServe(\":3333\", nil)\n" +
		"\tif err != nil {\n" +
		"\t\tpanic(err)\n" +
		"\t}" +
		"\n" +
		"}\n"

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
