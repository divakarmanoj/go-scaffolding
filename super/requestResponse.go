package main

type SuperResponse struct {
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

type SuperRequest struct {
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
