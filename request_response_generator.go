package main

import (
	"fmt"
	"go/format"
	"os"
)

// ConvertStruct function accepts a Structure and converts it into a string
func ConvertStruct(s *Structure, IsRequest bool, IsResponse bool) (string, error) {
	output := ""
	if IsRequest && IsResponse {
		return "", fmt.Errorf("cannot be both request and response")
	}
	if IsRequest {
		output += fmt.Sprintf("type %sRequest struct {\n", s.Name)
	} else if IsResponse {
		output += fmt.Sprintf("type %sResponse struct {\n\timports.Response\n", s.Name)

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
				output += fmt.Sprintf("\t%s\t*%sRequest\t`json:\"%s\"`\n", toCamelCase(v.Name), v.Name, toSnakeCase(v.Name))
			}
			if IsResponse {
				output += fmt.Sprintf("\t%s\t*%sResponse\t`json:\"%s\"`\n", toCamelCase(v.Name), v.Name, toSnakeCase(v.Name))
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

func GenerateRequestResponse(parsedStruct *Structure) string {
	output := "package main\n\n"
	output += generateImports([]string{"github.com/divakarmanoj/go-scaffolding/imports"})

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
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	// write rawOutput to file
	err = os.WriteFile(toSnakeCase(parsedStruct.Name)+"/requestResponse.go", outputBytes, 0644)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	return string(outputBytes)
}
