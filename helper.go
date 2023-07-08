package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ParseStruct function accepts a string and tries to parse it into a Struct
func ParseStruct(s string) (*Structure, error) {
	var Data Structure

	err := json.Unmarshal([]byte(s), &Data)
	if err != nil {
		return nil, err
	}
	return &Data, nil
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

func generateImports(imports []string) string {
	if len(imports) == 0 {
		return ""
	}
	output := "import (\n"
	for _, v := range imports {
		if strings.Contains(v, "-") {
			packageSplit := strings.Split(v, "/")
			pakageName := packageSplit[len(packageSplit)-1]
			// replace - with _
			pakageName = strings.Replace(pakageName, "-", "_", -1)
			output += fmt.Sprintf("%s \"%s\"\n", pakageName, v)
		} else {
			output += fmt.Sprintf("\"%s\"\n", v)
		}
	}
	output += ")\n\n"
	return output
}
