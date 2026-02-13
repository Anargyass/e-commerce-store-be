package main

import (
	"context"
	"fmt"
	"net/http"
	"github.com/golang-jwt/jwt/v5"
)

// Gunakan HandlerFunc agar lebih fleksibel saat dipasang di main.go
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers untuk semua request
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		// Handle CORS preflight request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 1. Cek header Authorization
		authHeader := r.Header.Get("Authorization")
		
		// Validasi apakah header ada dan dimulai dengan "Bearer "
		if authHeader == "" || len(authHeader) < 7 {
			http.Error(w, "Butuh login, Token tidak ada", http.StatusUnauthorized)
			return
		}

		// 2. Ambil token (Menghilangkan "Bearer ")
		tokenString := authHeader[7:]

		// 3. Validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan metode signing-nya HMAC (sesuai saat kita buat token)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Metode signing tidak valid")
			}
			// Gunakan jwtKey yang sama dengan yang ada di loginHandler
			return jwtKey, nil 
		})

		// 4. Jika token error atau tidak valid
		if err != nil || !token.Valid {
			http.Error(w, "Token tidak valid atau kadaluarsa", http.StatusUnauthorized)
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Bisa akses data user_id dari claims jika perlu
			userID := claims["user_id"]
			ctx := context.WithValue(r.Context(), "user_id", userID)
			r = r.WithContext(ctx)
		}

		// 5. Lanjut ke handler berikutnya (misal: productsHandler)
		next.ServeHTTP(w, r)
	}
}