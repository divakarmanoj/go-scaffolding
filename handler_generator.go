package main

import (
	"fmt"
	"go/format"
	"os"
)

func GenerateHandler(parsedStruct *Structure) {
	output := "package main\n\n"
	imports := []string{"encoding/json", "net/http", "strconv", "github.com/divakarmanoj/go-scaffolding/imports"}
	output += generateImports(imports)
	output += GenerateHandlerCreate(parsedStruct)
	output += "\n\n"
	output += GenerateHandlerRead(parsedStruct)
	output += "\n\n"
	output += GenerateHandlerUpdate(parsedStruct)
	output += "\n\n"
	output += GenerateHandlerDelete(parsedStruct)

	outputBytes, err := format.Source([]byte(output))
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}

	// write rawOutput to file
	err = os.WriteFile(toSnakeCase(parsedStruct.Name)+"/handler.go", outputBytes, 0644)
	if err != nil {
		fmt.Println("Error:" + err.Error())
		os.Exit(1)
	}
}

func GenerateHandlerCreate(parsedStruct *Structure) string {

	//output += fmt.Sprintf("func Create%s(w http.ResponseWriter, r *http.Request) {\n\tvar %s %sModel\n\tjson.NewDecoder(r.Body).Decode(&%s)\n\t%s.Create%s(&%s)\n\tjson.NewEncoder(w).Encode(%s)\n}", parsedStruct.Name, toCamelCase(parsedStruct.Name), parsedStruct.Name, toCamelCase(parsedStruct.Name), toSnakeCase(parsedStruct.Name), parsedStruct.Name, toCamelCase(parsedStruct.Name), toCamelCase(parsedStruct.Name))
	output := fmt.Sprintf("func Create%s(w http.ResponseWriter, r *http.Request) {\n", toCamelCase(parsedStruct.Name))
	output += fmt.Sprintf("\tvar %s %sRequest\n", toCamelCase(parsedStruct.Name), parsedStruct.Name)
	output += fmt.Sprintf("\terr := json.NewDecoder(r.Body).Decode(&%s)\n", toCamelCase(parsedStruct.Name))
	output += "\tif err != nil {\n"
	output += "\t\thttp.Error(w, err.Error(), http.StatusBadRequest)\n"
	output += "\t\treturn \n"
	output += "\t}\n"
	output += "\n"
	output += fmt.Sprintf("\tmodel := RequestTo%s(&%s)\n", toCamelCase(parsedStruct.Name), toCamelCase(parsedStruct.Name))
	output += "\tif err := db.Create(model).Error; err != nil {\n"
	output += "\t\thttp.Error(w, err.Error(), http.StatusBadRequest)\n"
	output += "\t\treturn\n"
	output += "\t}\n"
	output += "\n"
	output += fmt.Sprintf("\tjson.NewEncoder(w).Encode(%s)\n", toCamelCase(parsedStruct.Name))
	output += "}"
	return output

}

func GenerateHandlerRead(parsedStruct *Structure) string {
	output := fmt.Sprintf("func Read%s(w http.ResponseWriter, r *http.Request) {\n", toCamelCase(parsedStruct.Name))
	output += "\t page_numbers, ok := r.URL.Query()[\"page_number\"]\n"
	output += "\t page_number := 1\n"
	output += "\t if ok && len(page_numbers[0]) > 1 {\n"
	output += "\t\t var err error\n"
	output += "\t\t page_number, err = strconv.Atoi(page_numbers[0])\n"
	output += "\t\t if err != nil {\n"
	output += "\t\t\t http.Error(w, err.Error(), http.StatusBadRequest)\n"
	output += "\t\t\t return\n"
	output += "\t\t }\n"
	output += "\t }\n"
	output += "\n"
	output += "\t page_sizes, ok := r.URL.Query()[\"page_size\"]\n"
	output += "\t page_size := 10\n"
	output += "\t if ok && len(page_sizes[0]) > 1 {\n"
	output += "\t\t var err error\n"
	output += "\t\t page_size, err = strconv.Atoi(page_sizes[0])\n"
	output += "\t\t if err != nil {\n"
	output += "\t\t\t http.Error(w, err.Error(), http.StatusBadRequest)\n"
	output += "\t\t\t return\n"
	output += "\t\t }\n"
	output += "\t }\n"
	output += "\n"
	output += fmt.Sprintf("\tvar %s []%sModel\n", toCamelCase(parsedStruct.Name), parsedStruct.Name)
	output += fmt.Sprintf("\t if err := db.Limit(page_size).Offset((page_number - 1) * page_size).Find(&%s).Error; err != nil {\n", toCamelCase(parsedStruct.Name))
	output += "\t\thttp.Error(w, err.Error(), http.StatusBadRequest)\n"
	output += "\t\treturn\n"
	output += "\t}\n"
	output += "\n"
	output += fmt.Sprintf("var data []*%sResponse\n", toCamelCase(parsedStruct.Name))
	output += fmt.Sprintf("\tfor i, _ := range %s {\n", toCamelCase(parsedStruct.Name))
	output += fmt.Sprintf("\t\tdata = append(data, ModelTo%s(&%s[i]))\n", toCamelCase(parsedStruct.Name), toCamelCase(parsedStruct.Name))
	output += "\t}\n"
	output += "var output = imports.Response{\n"
	output += "\t\tData: data,\n"
	output += "\t\tStatus: \"success\",\n"
	output += "\t\tMessage: \"\",\n"
	output += "\t\t}\n"
	output += "\n"
	output += fmt.Sprintf("\tjson.NewEncoder(w).Encode(output)\n")
	output += "}"

	return output
}

func GenerateHandlerUpdate(parsedStruct *Structure) string {
	output := fmt.Sprintf("func Update%s(w http.ResponseWriter, r *http.Request) {\n", toCamelCase(parsedStruct.Name))
	output += "\t ids, ok := r.URL.Query()[\"id\"]\n"
	output += "\t if !ok || len(ids[0]) < 1 {\n"
	output += "\t\thttp.Error(w, \"Url Param 'id' is missing\", http.StatusBadRequest)\n"
	output += "\t\treturn\n"
	output += "\t}\n"
	output += "\t id := ids[0]\n"
	output += "\n"
	output += fmt.Sprintf("\tvar %s %sRequest\n", toCamelCase(parsedStruct.Name), parsedStruct.Name)
	output += fmt.Sprintf("\terr := json.NewDecoder(r.Body).Decode(&%s)\n", toCamelCase(parsedStruct.Name))
	output += "\tif err != nil {\n"
	output += "\t\thttp.Error(w, err.Error(), http.StatusBadRequest)\n"
	output += "\t\treturn \n"
	output += "\t}\n"
	output += "\n"
	output += fmt.Sprintf("\tmodel := RequestTo%s(&%s)\n", toCamelCase(parsedStruct.Name), toCamelCase(parsedStruct.Name))
	output += "\tif err := db.Model(&model).Where(\"id = ?\", id).Updates(model).Error; err != nil {\n"
	output += "\t\thttp.Error(w, err.Error(), http.StatusBadRequest)\n"
	output += "\t\treturn\n"
	output += "\t}\n"
	output += "\n"
	output += fmt.Sprintf("\tjson.NewEncoder(w).Encode(%s)\n", toCamelCase(parsedStruct.Name))
	output += "}"
	return output
}

func GenerateHandlerDelete(parsedStruct *Structure) string {
	output := fmt.Sprintf("func Delete%s(w http.ResponseWriter, r *http.Request) {\n", toCamelCase(parsedStruct.Name))
	output += "\tids, ok := r.URL.Query()[\"id\"]\n"
	output += "\tif !ok || len(ids[0]) < 1 {\n"
	output += "\t\thttp.Error(w, \"Url Param 'id' is missing\", http.StatusBadRequest)\n"
	output += "\t\treturn\n"
	output += "\t}\n"
	output += "\tid := ids[0]\n"
	output += fmt.Sprintf("\t if err:= db.Delete(&%sModel{}, id).Error; err != nil {\n", toCamelCase(parsedStruct.Name))
	output += "\t\thttp.Error(w, err.Error(), http.StatusBadRequest)\n"
	output += "\t\treturn\n"
	output += "\t}\n"
	output += fmt.Sprintf("\tjson.NewEncoder(w).Encode(imports.Response{Message: \"%s deleted successfully\", Status: \"success\"})", toCamelCase(parsedStruct.Name))
	output += "}"
	return output
}
