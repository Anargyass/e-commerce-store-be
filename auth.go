package main

import (
	"encoding/json"
	"net/http"
	"golang.org/x/crypto/bcrypt"
)

// User adalah model untuk tabel users di database
type User struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Username string `json:"username" gorm:"unique;not null"`
    Email    string `json:"email" gorm:"unique;not null"`
    Password string `json:"-"` // Password disembunyikan dari JSON demi keamanan
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	// next js cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")


	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST"{
		var user User // buat nampung data user

		// 1. ambil data dari frontend
		json.NewDecoder(r.Body).Decode(&user)

		// 2. hash password
		hashedPassword,_:=bcrypt.GenerateFromPassword([]byte(user.Password), 10)
		user.Password=string(hashedPassword)

		// 3. simpan ke database, buat handle errornya
		result := DB.Create(&user)
		if result.Error != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return
		}
		// 4. kirim response sukses
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})

	}
}
