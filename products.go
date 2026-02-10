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
	
	// next js cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")


	if r.Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case "GET":
		var products []Product
		DB.Find(&products) // SQL: SELECT * FROM products WHERE deleted_at IS NULL;
		json.NewEncoder(w).Encode(products)

	case "POST":
		var newProd Product
		json.NewDecoder(r.Body).Decode(&newProd)
		DB.Create(&newProd) // SQL: INSERT INTO products (name, price...) VALUES (...);
		json.NewEncoder(w).Encode(newProd)

	case "DELETE":
		id := r.URL.Query().Get("id")
		// GORM Soft Delete: Data tidak benar-benar hilang, tapi ditandai 'deleted'
		DB.Delete(&Product{}, id) 
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Berhasil dihapus dari Database")

	}

}

