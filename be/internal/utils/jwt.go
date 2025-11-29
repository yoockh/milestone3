package utils

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// for generate and validate jwt token
func GenerateJwtToken(email, role string, id int) (string, error) {
	expired_jwt := os.Getenv("EXPIRED_JWT")
	expired_jwt_int, err := strconv.Atoi(expired_jwt)
	if err != nil {
		log.Println("error parsing expired_jwt to int")
		return "", err
	}
	expired := time.Now().Add(time.Hour * time.Duration(expired_jwt_int)).Unix()
	jwt_claim := jwt.NewWithClaims(jwt.SigningMethodHS256, 
	jwt.MapClaims{
		"id": id,
		"email": email,
		"role": role,
		"exp": expired,
	})
		
	secret_key := os.Getenv("SECRET_KEY")
	tokenString, err := jwt_claim.SignedString([]byte(secret_key))
	if err != nil {
		log.Println("error signed string jwt")
		return "", err
	}

	return tokenString, nil
}