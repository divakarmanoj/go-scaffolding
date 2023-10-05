package main

import imports "github.com/divakarmanoj/go-scaffolding/imports"

type ExampleResponse struct {
	ID        uint             `json:"id"`
	CreatedAt int64            `json:"created_at"`
	UpdatedAt int64            `json:"updated_at"`
	Name      *string          `json:"name,omitempty"`
	Age       int16            `json:"age"`
	Address   *AddressResponse `json:"address"`
}

type AddressResponse struct {
	ID         uint    `json:"id"`
	CreatedAt  int64   `json:"created_at"`
	UpdatedAt  int64   `json:"updated_at"`
	StreetName string  `json:"street_name"`
	City       string  `json:"city"`
	State      *string `json:"state,omitempty"`
	Zip        int16   `json:"zip"`
}

type ExampleRequest struct {
	Name    *string         `json:"name,omitempty"`
	Age     int16           `json:"age"`
	Address *AddressRequest `json:"address"`
}

type AddressRequest struct {
	StreetName string  `json:"street_name"`
	City       string  `json:"city"`
	State      *string `json:"state,omitempty"`
	Zip        int16   `json:"zip"`
}

func (request *ExampleRequest) ToModel() *ExampleModel {
	if request == nil {
		return nil
	}
	return &ExampleModel{Name: imports.NullStringPtr(request.Name), Age: request.Age, Address: request.Address.ToModel()}
}
func (request *AddressRequest) ToModel() *AddressModel {
	if request == nil {
		return nil
	}
	return &AddressModel{StreetName: request.StreetName, City: request.City, State: imports.NullStringPtr(request.State), Zip: request.Zip}
}
