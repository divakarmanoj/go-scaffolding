package main

import (
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"strings"
)

type Attributes struct {
	Name       string       `json:"name"`
	Type       string       `json:"type"`
	Attributes []Attributes `json:"attributes"`
	IsRequired bool         `json:"is_required"`
}

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
					"is_required": true
				},
				{
					"Name": "Zip",
					"Type": "int16",
					"is_required": true
				}
			],
			"is_required": true
		}
	]
}`

// ParseStruct function accepts a string and tries to parse it into a Struct
func ParseStruct(s string) (*Structure, error) {
	var Data Structure

	err := json.Unmarshal([]byte(s), &Data)
	if err != nil {
		return nil, err
	}
	return &Data, nil
}

// ConvertStruct function accepts a Structure and converts it into a string
func ConvertStruct(s *Structure, IsRequest bool, IsResponse bool) (string, error) {
	output := ""
	if IsRequest && IsResponse {
		return "", fmt.Errorf("cannot be both request and response")
	}
	if IsRequest {
		output += fmt.Sprintf("type %sRequest struct {\n", s.Name)
	} else if IsResponse {
		output += fmt.Sprintf("type %sResponse struct {\n\tResponse\n", s.Name)

	} else {
		output += fmt.Sprintf("type %s struct {\n", s.Name)
	}

	structs := []string{}
	for _, v := range s.Attributes {
		if !isValidType(v.Type) {
			return "", fmt.Errorf("invalid type %s", v.Type)
		}

		if v.Type != "struct" {
			if v.IsRequired {
				output += fmt.Sprintf("\t%s\t%s\t`json:\"%s\"`\n", toCamelCase(v.Name), v.Type, toSnakeCase(v.Name))
			} else {
				output += fmt.Sprintf("\t%s\t*%s\t`json:\"%s,omitempty\"`\n", toCamelCase(v.Name), v.Type, toSnakeCase(v.Name))
			}
		} else {
			nestedStructs, err := ConvertStruct(&Structure{Name: v.Name, Attributes: v.Attributes}, IsRequest, IsResponse)
			if err != nil {
				return "", err
			}
			if IsRequest {
				output += fmt.Sprintf("\t%s\t%sRequest\t`json:\"%s\"`\n", toCamelCase(v.Name), v.Name, toSnakeCase(v.Name))
			}
			if IsResponse {
				output += fmt.Sprintf("\t%s\t%sResponse\t`json:\"%s\"`\n", toCamelCase(v.Name), v.Name, toSnakeCase(v.Name))
			}
			structs = append(structs, nestedStructs)
		}
	}
	output += fmt.Sprintf("}")
	if len(structs) > 0 {
		output += "\n\n"
	}
	for _, v := range structs {
		output += v
	}
	if len(structs) == 0 {
		output += "\n\n"
	}
	return output, nil
}

// ConvertGormStruct function accepts a Structure and converts it into a string
func ConvertGormStruct(s *Structure) (string, error) {
	output := ""
	output += fmt.Sprintf("type %sModel struct {\n\tModel\n", s.Name)

	structs := []string{}
	for _, v := range s.Attributes {
		if !isValidType(v.Type) {
			return "", fmt.Errorf("invalid type %s", v.Type)
		}

		if v.Type != "struct" {
			if v.IsRequired {
				output += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", toCamelCase(v.Name), v.Type, toSnakeCase(v.Name))
			} else {
				output += fmt.Sprintf("\t%s *%s `json:\"%s,omitempty\"`\n", toCamelCase(v.Name), v.Type, toSnakeCase(v.Name))
			}
		} else {
			nestedStructs, err := ConvertGormStruct(&Structure{Name: v.Name, Attributes: v.Attributes})
			if err != nil {
				return "", err
			}
			output += fmt.Sprintf("\t%sModelID\tuint\t`json:\"%s_id\"`\n", toCamelCase(v.Name), toSnakeCase(v.Name))
			output += fmt.Sprintf("\t%s\t%sModel\t`json:\"%s\"`\n", toCamelCase(v.Name), v.Name, toSnakeCase(v.Name))
			structs = append(structs, nestedStructs)
		}
	}
	output += fmt.Sprintf("}")
	if len(structs) > 0 {
		output += "\n\n"
	}
	for _, v := range structs {
		output += v
	}
	if len(structs) == 0 {
		output += "\n\n"
	}
	return output, nil
}

func toCamelCase(name string) string {
	words := strings.Split(name, " ")
	result := ""

	for _, word := range words {
		result += strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
	}

	return result
}

func toSnakeCase(name string) string {
	name = strings.Replace(name, " ", "_", -1)
	name = strings.ToLower(name)

	return name
}

func isValidType(fieldType string) bool {
	validTypes := []string{
		"string",
		"int64",
		"int32",
		"int16",
		"float64",
		"bool",
		"time.Time",
		"struct",
	}
	for _, t := range validTypes {
		if fieldType == t {
			return true
		}
	}
	return false
}

func main() {
	parsedStruct, err := ParseStruct(input)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}

	rawOutput := GenerateRequestResponseModels(parsedStruct)
	rawOutput += "\n"
	gormModel, err := ConvertGormStruct(parsedStruct)
	rawOutput += gormModel

	outputBytes, err := format.Source([]byte(rawOutput))
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	fmt.Print(string(outputBytes))
}

func GenerateRequestResponseModels(parsedStruct *Structure) string {
	output := "package main\n\n"

	structs, err := ConvertStruct(parsedStruct, false, true)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	output += structs
	structs, err = ConvertStruct(parsedStruct, true, false)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	output += structs
	outputBytes, err := format.Source([]byte(output))
	if err != nil {

	}
	return string(outputBytes)
}

func generateImports(imports []string) string {
	if len(imports) == 0 {
		return ""
	}
	output := "import (\n"
	for _, v := range imports {
		output += fmt.Sprintf("\"%s\"\n", v)
	}
	output += ")\n\n"
	return output
}
