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
		Message: "super created successfully",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(output)
}

func ReadSuper(w http.ResponseWriter, r *http.Request) {
	var err error

	// get cursor based pagination

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

	var Super []SuperModel
	err = db.Preload(clause.Associations).Limit(pageSize+1).Where("id >= ?", cursor).Order("id asc").Find(&Super).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var data []*SuperResponse
	for _, i := range Super {
		data = append(data, ModelToSuper(&i))
		if len(data) == pageSize {
			break
		}
	}
	var output = imports.Response{
		Data:    data,
		Message: "supers retrieved successfully",
		Status:  "success",
	}
	if len(Super) > pageSize {
		output.Cursor = Super[pageSize].ID
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
		Message: "super updated successfully",
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
		Message: "super deleted successfully",
		Status:  "success",
	}
	json.NewEncoder(w).Encode(output)
}
