package main

import (
	"database/sql"
	"github.com/divakarmanoj/go-scaffolding/imports"
)

type ExampleModel struct {
	imports.Model
	Name      sql.NullString `json:"name"`
	Age       int16          `json:"age"`
	AddressID uint           `json:"address_id"`
	Address   *AddressModel  `json:"address"`
}

type AddressModel struct {
	imports.Model
	StreetName string         `json:"street_name"`
	City       string         `json:"city"`
	State      sql.NullString `json:"state"`
	Zip        int16          `json:"zip"`
}

func (model *ExampleModel) ToResponse() *ExampleResponse {
	if model == nil {
		return nil
	}
	return &ExampleResponse{ID: model.Model.ID, CreatedAt: model.Model.CreatedAt, UpdatedAt: model.Model.UpdatedAt, Name: imports.NullStringToPtr(model.Name), Age: model.Age, Address: model.Address.ToResponse()}
}
func (model *AddressModel) ToResponse() *AddressResponse {
	if model == nil {
		return nil
	}
	return &AddressResponse{ID: model.Model.ID, CreatedAt: model.Model.CreatedAt, UpdatedAt: model.Model.UpdatedAt, StreetName: model.StreetName, City: model.City, State: imports.NullStringToPtr(model.State), Zip: model.Zip}
}
