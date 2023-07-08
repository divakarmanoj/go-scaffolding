package main

import (
	"encoding/json"
	imports "github.com/divakarmanoj/go-scaffolding/imports"
	"net/http"
	"strconv"
)

func CreateSuper(w http.ResponseWriter, r *http.Request) {
	var Super SuperRequest
	err := json.NewDecoder(r.Body).Decode(&Super)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := RequestToSuper(&Super)
	if err := db.Create(model).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(Super)
}

func ReadSuper(w http.ResponseWriter, r *http.Request) {
	page_numbers, ok := r.URL.Query()["page_number"]
	page_number := 1
	if ok && len(page_numbers[0]) > 1 {
		var err error
		page_number, err = strconv.Atoi(page_numbers[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	page_sizes, ok := r.URL.Query()["page_size"]
	page_size := 10
	if ok && len(page_sizes[0]) > 1 {
		var err error
		page_size, err = strconv.Atoi(page_sizes[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	var Super []SuperModel
	if err := db.Limit(page_size).Offset((page_number - 1) * page_size).Find(&Super).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data []*SuperResponse
	for i, _ := range Super {
		data = append(data, ModelToSuper(&Super[i]))
	}
	var output = imports.Response{
		Data:    data,
		Status:  "success",
		Message: "",
	}

	json.NewEncoder(w).Encode(output)
}

func UpdateSuper(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
		return
	}
	id := ids[0]

	var Super SuperRequest
	err := json.NewDecoder(r.Body).Decode(&Super)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := RequestToSuper(&Super)
	if err := db.Model(&model).Where("id = ?", id).Updates(model).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(Super)
}

func DeleteSuper(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		http.Error(w, "Url Param 'id' is missing", http.StatusBadRequest)
		return
	}
	id := ids[0]
	if err := db.Delete(&SuperModel{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(imports.Response{Message: "Super deleted successfully", Status: "success"})
}
