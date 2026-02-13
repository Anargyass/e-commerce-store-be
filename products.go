package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Product struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Price int 	`json:"price"`
}


func productsHandler(w http.ResponseWriter, r *http.Request) {
	
	// Set Content-Type untuk response
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		var products []Product
		DB.Find(&products)
		json.NewEncoder(w).Encode(products)

	case "POST":
		var newProd Product
		json.NewDecoder(r.Body).Decode(&newProd)
		DB.Create(&newProd)
		json.NewEncoder(w).Encode(newProd)

	case "DELETE":
		id := r.URL.Query().Get("id")
		DB.Delete(&Product{}, id) 
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Berhasil dihapus dari Database")
	}
}

