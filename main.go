package main

import (
	"fmt"
	"os"
)

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

	// Create directory if not exists with name of struct
	_ = os.Mkdir(toSnakeCase(parsedStruct.Name), 0755)
	GenerateRequestResponse(parsedStruct)
	_, model := GenerateModel(parsedStruct)
	GenerateAdaptor(parsedStruct)
	GenerateHandler(parsedStruct)
	GenerateMain(parsedStruct, model)
}
