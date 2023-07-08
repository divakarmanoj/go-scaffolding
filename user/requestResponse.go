package main

import (
	imports "github.com/divakarmanoj/go-scaffolding/imports"
)

type UserResponse struct {
	imports.Response
	Name    *string          `json:"name,omitempty"`
	Age     int16            `json:"age"`
	Address *AddressResponse `json:"address"`
}

type AddressResponse struct {
	imports.Response
	StreetName string  `json:"street_name"`
	City       string  `json:"city"`
	State      *string `json:"state,omitempty"`
	Zip        int16   `json:"zip"`
}

type UserRequest struct {
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
