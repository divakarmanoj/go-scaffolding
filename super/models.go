package main

import (
	"database/sql"
	"github.com/divakarmanoj/go-scaffolding/imports"
)

type SuperModel struct {
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
