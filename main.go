package main

import (
	"fmt"
	"go/format"
	"os"
	"strings"
)

type Structure struct {
	Name       string       `json:"name"`
	Attributes []Attributes `json:"attributes"`
}

var input = `{
	"name": "User",
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
	GenerateModel(parsedStruct)
	GenerateAdaptor(parsedStruct)
}

func GenerateAdaptor(s *Structure) string {
	output := "package main\n\n"
	requestToModelOutput, shouldImportImports := GenerateRequestToModel(s)
	if shouldImportImports {
		output += generateImports([]string{"github.com/divakarmanoj/go-scaffolding/imports"})
	}
	output += requestToModelOutput
	modelToResponse, _ := GenerateModelToResponse(s)
	output += modelToResponse

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

func GenerateModelToResponse(s *Structure) (string, bool) {
	output := ""
	output += fmt.Sprintf("func ModelTo%s(model *%sModel) *%sResponse {\n", s.Name, s.Name, s.Name)
	output += fmt.Sprintf("\treturn &%sResponse{\n", s.Name)
	nestedFunctions := []string{}
	shouldImportImports := false
	for _, v := range s.Attributes {
		if !isValidType(v.Type) {
			fmt.Println("Error: invalid type " + v.Type)
			os.Exit(1)
		}

		if v.Type != "struct" {
			if !v.IsRequired {
				shouldImportImports = true
				output += fmt.Sprintf("\t\t\t%s: imports.Null%sToPtr(model.%s),\n", toCamelCase(v.Name), strings.Title(v.Type), toCamelCase(v.Name))
			} else {
				output += fmt.Sprintf("\t\t%s: model.%s,\n", toCamelCase(v.Name), toCamelCase(v.Name))
			}
		} else {
			nestedFunction, nestedShouldImportImports := GenerateModelToResponse(&Structure{Name: v.Name, Attributes: v.Attributes})
			if nestedShouldImportImports {
				shouldImportImports = true
			}
			nestedFunctions = append(nestedFunctions, nestedFunction)
			output += fmt.Sprintf("\t\t%s: %s(model.%s),\n", toCamelCase(v.Name), "ModelTo"+v.Name, toCamelCase(v.Name))
		}
	}
	output += "\t}\n"
	output += "}\n"

	for _, v := range nestedFunctions {
		output += v
	}

	return output, shouldImportImports
}

func GenerateRequestToModel(s *Structure) (string, bool) {
	output := ""
	output += fmt.Sprintf("func RequestTo%s(request *%sRequest) *%sModel {\n", s.Name, s.Name, s.Name)
	output += fmt.Sprintf("\treturn &%sModel{\n", s.Name)
	nestedFunctions := []string{}
	shouldImportImports := false
	for _, v := range s.Attributes {
		if !isValidType(v.Type) {
			fmt.Println("Error: invalid type " + v.Type)
			os.Exit(1)
		}

		if v.Type != "struct" {
			if !v.IsRequired {
				shouldImportImports = true
				output += fmt.Sprintf("\t\t\t%s: imports.Null%sPtr(request.%s),\n", toCamelCase(v.Name), strings.Title(v.Type), toCamelCase(v.Name))
			} else {
				output += fmt.Sprintf("\t\t%s: request.%s,\n", toCamelCase(v.Name), toCamelCase(v.Name))
			}
		} else {
			nestedFunction, nestedShouldImportImports := GenerateRequestToModel(&Structure{Name: v.Name, Attributes: v.Attributes})
			if nestedShouldImportImports {
				shouldImportImports = true
			}
			nestedFunctions = append(nestedFunctions, nestedFunction)
			output += fmt.Sprintf("\t\t%s: RequestTo%s(request.%s),\n", toCamelCase(v.Name), v.Name, toCamelCase(v.Name))
		}
	}
	output += fmt.Sprintf("\t}\n")
	output += fmt.Sprintf("}\n")

	for _, v := range nestedFunctions {
		output += v
	}
	return output, shouldImportImports
}
