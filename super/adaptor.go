package main

import (
	imports "github.com/divakarmanoj/go-scaffolding/imports"
)

func RequestToSuper(request *SuperRequest) *SuperModel {
	if request == nil {
		return nil
	}
	return &SuperModel{
		Name:    imports.NullStringPtr(request.Name),
		Age:     request.Age,
		Address: RequestToAddress(request.Address),
	}
}
func RequestToAddress(request *AddressRequest) *AddressModel {
	if request == nil {
		return nil
	}
	return &AddressModel{
		StreetName: request.StreetName,
		City:       request.City,
		State:      imports.NullStringPtr(request.State),
		Zip:        request.Zip,
	}
}
func ModelToSuper(model *SuperModel) *SuperResponse {
	if model == nil {
		return nil
	}
	return &SuperResponse{
		ID:        model.Model.ID,
		CreatedAt: model.Model.CreatedAt,
		UpdatedAt: model.Model.UpdatedAt,
		Name:      imports.NullStringToPtr(model.Name),
		Age:       model.Age,
		Address:   ModelToAddress(model.Address),
	}
}
func ModelToAddress(model *AddressModel) *AddressResponse {
	if model == nil {
		return nil
	}
	return &AddressResponse{
		ID:         model.Model.ID,
		CreatedAt:  model.Model.CreatedAt,
		UpdatedAt:  model.Model.UpdatedAt,
		StreetName: model.StreetName,
		City:       model.City,
		State:      imports.NullStringToPtr(model.State),
		Zip:        model.Zip,
	}
}
