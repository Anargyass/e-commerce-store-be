package main

import (
    "fmt"
    "net/http"
)

func main() {
    // 1. Inisialisasi Database (dari database.go)
    InitDB()

    // 2. Daftarkan Route (Handler dari products.go)
    http.HandleFunc("/api/products", authMiddleware(productsHandler))
	http.HandleFunc("/api/register", registerHandler)
    http.HandleFunc("/api/login", loginHandler)

    // 3. Jalankan Server
    fmt.Println("Server running at http://localhost:8080/api/products")
	fmt.Println("Server running at http://localhost:8080/api/register")
    fmt.Println("Server running at http://localhost:8080/api/login")

    http.ListenAndServe(":8080", nil)
}
