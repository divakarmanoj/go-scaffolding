package main

type Attributes struct {
	Name       string       `json:"name"`
	Type       string       `json:"type"`
	Attributes []Attributes `json:"attributes"`
	IsRequired bool         `json:"is_required"`
}
