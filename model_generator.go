package main

import (
	"fmt"
	"go/format"
	"os"
	"strings"
)

func GenerateModel(s *Structure) (string, []string) {
	output := "package main\n\n"
	imports := []string{"github.com/divakarmanoj/go-scaffolding/imports", "gorm.io/gorm"}
	modelOutput, importSQL, modelNames := ConvertGormStruct(s)

	if importSQL {
		imports = append(imports, "database/sql")
	}
	output += generateImports(imports)

	output += "\n\nvar db *gorm.DB\n\n"
	output += modelOutput
	outputBytes, err := format.Source([]byte(output))
	if err != nil {
		fmt.Println(output)
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	// write rawOutput to file
	err = os.WriteFile(toSnakeCase(s.Name)+"/models.go", outputBytes, 0644)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
	return output, modelNames
}

// ConvertGormStruct function accepts a Structure and converts it into a string
func ConvertGormStruct(s *Structure) (string, bool, []string) {
	output := ""
	output += fmt.Sprintf("type %sModel struct {\n\timports.Model\n", s.Name)
	importSQL := false
	structs := []string{}
	modelNames := []string{s.Name + "Model"}
	for _, v := range s.Attributes {
		if !isValidType(v.Type) {
			fmt.Println("Error: invalid type %s", v.Type)
			os.Exit(1)
		}

		if v.Type != "struct" {
			if v.IsRequired {
				output += fmt.Sprintf("\t%s %s `json:\"%s\"`\n", toCamelCase(v.Name), v.Type, toSnakeCase(v.Name))
			} else {
				importSQL = true
				output += fmt.Sprintf("\t%s sql.Null%s `json:\"%s\"`\n", toCamelCase(v.Name), strings.Title(v.Type), toSnakeCase(v.Name))
			}
		} else {
			nestedStructs, nestedImportSQL, NestedModelNames := ConvertGormStruct(&Structure{Name: v.Name, Attributes: v.Attributes})
			modelNames = append(modelNames, NestedModelNames...)
			if nestedImportSQL {
				importSQL = true
			}
			output += fmt.Sprintf("\t%sID\tuint\t`json:\"%s_id\"`\n", toCamelCase(v.Name), toSnakeCase(v.Name))
			output += fmt.Sprintf("\t%s\t*%sModel\t`json:\"%s\"`\n", toCamelCase(v.Name), v.Name, toSnakeCase(v.Name))
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
	return output, importSQL, modelNames
}
