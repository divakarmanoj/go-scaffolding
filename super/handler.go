package main

import (
	"encoding/json"
	"github.com/divakarmanoj/go-scaffolding/imports"
	"gorm.io/gorm/clause"
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

	var output = imports.Response{
		Data:    ModelToSuper(model),
		Message: "Super created successfully",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(output)
}

func ReadSuper(w http.ResponseWriter, r *http.Request) {
	var err error
	pageNumbers, ok := r.URL.Query()["page_number"]
	pageNumber := 1
	if ok && len(pageNumbers[0]) > 1 {
		pageNumber, err = strconv.Atoi(pageNumbers[0])
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

	var Super []SuperModel
	err = db.Limit(pageSize).Offset(pageSize * (pageNumber - 1)).Preload(clause.Associations).Find(&Super).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data []*SuperResponse
	for _, i := range Super {
		data = append(data, ModelToSuper(&i))
	}
	var output = imports.Response{
		Data:    data,
		Message: "Supers retrieved successfully",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(output)
}

func UpdateSuper(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	id := ids[0]

	var Super SuperRequest
	if err := json.NewDecoder(r.Body).Decode(&Super); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	model := RequestToSuper(&Super)
	if err := db.Model(&model).Where("id = ?", id).Updates(&model).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var output = imports.Response{
		Data:    ModelToSuper(model),
		Message: "Super updated successfully",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(output)
}

func DeleteSuper(w http.ResponseWriter, r *http.Request) {
	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	id := ids[0]
	if err := db.Delete(&SuperModel{}, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var output = imports.Response{
		Message: "Super deleted successfully",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(output)
}
