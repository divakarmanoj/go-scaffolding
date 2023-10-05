package main

import (
	"encoding/json"
	"github.com/divakarmanoj/go-scaffolding/imports"
	"gorm.io/gorm/clause"
	"net/http"
	"strconv"
)

func CreateExample(w http.ResponseWriter, r *http.Request) {
	var Example ExampleRequest
	err := json.NewDecoder(r.Body).Decode(&Example)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := Example.ToModel()
	if err := db.Create(model).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var output = imports.Response{
		Data:    model.ToResponse(),
		Message: "example created successfully",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(output)
}

func ReadExample(w http.ResponseWriter, r *http.Request) {
	var err error
	cursors, ok := r.URL.Query()["cursor"]
	cursor := 1
	if ok && len(cursors[0]) > 1 {
		cursor, err = strconv.Atoi(cursors[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	pageSizes, ok := r.URL.Query()["page_size"]
	pageSize := 10
	if ok && len(pageSizes[0]) > 1 {
		pageSize, err = strconv.Atoi(pageSizes[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var Example []ExampleModel
	err = db.Preload(clause.Associations).Limit(pageSize+1).Where("id >= ?", cursor).Order("id asc").Find(&Example).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data []*ExampleResponse
	for _, i := range Example {
		data = append(data, i.ToResponse())
		if len(data) == pageSize {
			break
		}
	}
	var output = imports.Response{
		Data:    data,
		Message: "examples retrieved successfully",
		Status:  "success",
	}
	if len(Example) > pageSize {
		output.Cursor = Example[pageSize].ID
	}
	json.NewEncoder(w).Encode(output)
}

func UpdateExample(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	id := ids[0]

	var Example ExampleRequest
	if err := json.NewDecoder(r.Body).Decode(&Example); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := Example.ToModel()
	if err := db.Model(&model).Where("id = ?", id).Updates(&model).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var output = imports.Response{
		Data:    model.ToResponse(),
		Message: "example updated successfully",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(output)
}

func DeleteExample(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	id := ids[0]
	if err := db.Delete(&ExampleModel{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var output = imports.Response{
		Message: "example deleted successfully",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(output)
}
