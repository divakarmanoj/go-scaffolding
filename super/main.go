package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"net/http"
	"os"
)

func main() {
	var err error
	db, err = gorm.Open(sqlite.Open("Super.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	defer func() {
		err = os.RemoveAll("Super.db")
		if err != nil {
			return
		}
	}()
	db.AutoMigrate(&SuperModel{})
	db.AutoMigrate(&AddressModel{})

	http.HandleFunc("/super/read", ReadSuper)
	http.HandleFunc("/super/create", CreateSuper)
	http.HandleFunc("/super/update", UpdateSuper)
	http.HandleFunc("/super/delete", DeleteSuper)
	err = http.ListenAndServe(":3333", nil)
	if err != nil {
		panic(err)
	}
}
