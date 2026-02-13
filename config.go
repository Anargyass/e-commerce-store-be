package main

import (
	"fmt"
	"os"
	"unicode"

	"github.com/joho/godotenv"
)

var jwtKey []byte

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Peringatan: Tidak bisa memuat file .env, menggunakan sistem env")
	}

	secret := os.Getenv("JWT_SECRET")
	if err := validateJWTSecret(secret); err != nil {
		panic(err)
	}
	jwtKey = []byte(secret)
}

func validateJWTSecret(secret string) error {
	if len(secret) < 32 {
		return fmt.Errorf("JWT_SECRET minimal 32 karakter")
	}

	var hasUpper, hasLower, hasDigit, hasSymbol bool
	for _, r := range secret {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		default:
			hasSymbol = true
		}
	}

	if !hasUpper || !hasLower || !hasDigit || !hasSymbol {
		return fmt.Errorf("JWT_SECRET harus mengandung huruf besar, huruf kecil, angka, dan simbol")
	}

	return nil
}
