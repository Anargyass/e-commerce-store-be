package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Name  string `json:"name"`
	Price int 	`json:"price"`
}

var DB *gorm.DB

func initDB() {
	// 1. Load file .env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Peringatan: Tidak bisa memuat file .env, menggunakan sistem env")
	}

	// 2. Ambil data dari .env
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// 3. Susun DSN
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port)

	// 4. Koneksi (Hapus 'var err error' karena err sudah ada di atas)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	DB.AutoMigrate(&Product{})
	fmt.Println("Database connected and migrated")
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
func main() {
	initDB()
	http.HandleFunc("/api/products", productsHandler)
	fmt.Println("Server running at http://localhost:8080/api/products")
	http.ListenAndServe(":8080", nil)
}
