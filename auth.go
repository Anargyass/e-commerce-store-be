package main

import (
	"encoding/json" // "encoding/json" untuk parsing JSON
	"net/http" // "net/http" untuk menangani request dan response
	"golang.org/x/crypto/bcrypt" // "golang.org/x/crypto/bcrypt" untuk hashing password
	"time"// "time" untuk mengatur waktu kadaluarsa token

	"github.com/golang-jwt/jwt/v5"// "github.com/golang-jwt/jwt/v5" untuk membuat dan memverifikasi JWT
)

// User adalah model untuk tabel users di database
type User struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Username string `json:"username" gorm:"unique;not null"`
    Email    string `json:"email" gorm:"unique;not null"`
    Password string `json:"-"` // Password disembunyikan dari JSON demi keamanan
}


func loginHandler(w http.ResponseWriter, r *http.Request) {
	// next js cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
		var input struct {
			Username string `json:"email"`
			Password string `json:"password"`
		}
		json.NewDecoder(r.Body).Decode(&input) // error handling diabaikan untuk singkatnya

		// 1. Cari user berdasarkan username/email
        // buat rate limit agar tidak mudah di-brute force
        var loginAttempts int64
        DB.Model(&User{}).Where("email = ?", input.Username).Count(&loginAttempts)
        if loginAttempts > 5 {
            http.Error(w, "Too many login attempts", http.StatusTooManyRequests)
            return
        }

		var user User
		if err := DB.Where("email = ?", input.Username).First(&user).Error; err != nil {
			http.Error(w, "Email tidak ditemukan", http.StatusUnauthorized)
			return
		}

		// 2. Bandingkan password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			http.Error(w, "Password salah", http.StatusUnauthorized)
			return
		}

		// 3. Buat JWT Token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": user.ID,
			"exp":     time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
		})
		tokenString, _ := token.SignedString(jwtKey) // Ganti dengan secret key yang aman

		// 4. Kirim token ke frontend
		json.NewEncoder(w).Encode(map[string]string{
			"token": tokenString,
			"message": "Login successful",
		})
	}

		
	}



func registerHandler(w http.ResponseWriter, r *http.Request) {
	// next js cors
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")


	if r.Method == "OPTIONS" {
		return
	}

	if r.Method == "POST" {
        // Struct sementara untuk menangkap password asli sebelum di-hash
        // karena User struct kamu punya tag `json:"-"` pada password
        var input struct {
            Username string `json:"username"`
            Email    string `json:"email"`
            Password string `json:"password"`
        }

        // 1. Ambil data dari frontend
        if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
            http.Error(w, "Invalid input", http.StatusBadRequest)
            return
        }

        // 2. VALIDASI LOGIKA (Manual)
        if len(input.Username) < 4 {
            http.Error(w, "Username minimal 4 karakter", http.StatusBadRequest)
            return
        }
        if len(input.Password) < 6 {
            http.Error(w, "Password minimal 6 karakter", http.StatusBadRequest)
            return
        }

        // 3. VALIDASI DATABASE (Cek Duplikat)
        var existingUser User
        // Cek username
        if err := DB.Where("username = ?", input.Username).First(&existingUser).Error; err == nil {
            http.Error(w, "Username sudah digunakan", http.StatusConflict)
            return
        }
        // Cek email
        if err := DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
            http.Error(w, "Email sudah terdaftar", http.StatusConflict)
            return
        }

        // 4. HASH PASSWORD
        hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)

        // 5. SIMPAN KE DATABASE
        newUser := User{
            Username: input.Username,
            Email:    input.Email,
            Password: string(hashedPassword),
        }

        if err := DB.Create(&newUser).Error; err != nil {
            http.Error(w, "Gagal simpan data", http.StatusInternalServerError)
            return
        }

        w.WriteHeader(http.StatusCreated)
        json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
    }
}
