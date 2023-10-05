package main

import (
	sqlite "gorm.io/driver/sqlite"
	gorm "gorm.io/gorm"
	"net/http"
	"os"
)

var db *gorm.DB

func main() {
	var err error

	db, err = gorm.Open(sqlite.Open("Example.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	defer func() {
		err = os.RemoveAll("Example.db")
		if err != nil {
			return
		}
	}()

	db.AutoMigrate(&ExampleModel{})
	db.AutoMigrate(&AddressModel{})

	http.HandleFunc("/example/read", ReadExample)
	http.HandleFunc("/example/create", CreateExample)
	http.HandleFunc("/example/update", UpdateExample)
	http.HandleFunc("/example/delete", DeleteExample)
	err = http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
}
