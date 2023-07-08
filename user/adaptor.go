package main

import (
	imports "github.com/divakarmanoj/go-scaffolding/imports"
)

func RequestToUser(request *UserRequest) *UserModel {
	return &UserModel{
		Name:    imports.NullStringPtr(request.Name),
		Age:     request.Age,
		Address: RequestToAddress(request.Address),
	}
}
func RequestToAddress(request *AddressRequest) *AddressModel {
	return &AddressModel{
		StreetName: request.StreetName,
		City:       request.City,
		State:      imports.NullStringPtr(request.State),
		Zip:        request.Zip,
	}
}
func ModelToUser(model *UserModel) *UserResponse {
	return &UserResponse{
		Name:    imports.NullStringToPtr(model.Name),
		Age:     model.Age,
		Address: ModelToAddress(model.Address),
	}
}
func ModelToAddress(model *AddressModel) *AddressResponse {
	return &AddressResponse{
		StreetName: model.StreetName,
		City:       model.City,
		State:      imports.NullStringToPtr(model.State),
		Zip:        model.Zip,
	}
}
