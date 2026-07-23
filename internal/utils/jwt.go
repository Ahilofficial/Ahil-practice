package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"time"
	

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID uint, sessionID string) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"session_id":sessionID,
		"user_id": userID,
		"iat":     now.Unix(),
		"exp":     now.Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "supersecretkey"
	}

	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(userID uint, sessionID string) (string, error) {
	
	now := time.Now()
	claims := jwt.MapClaims{
		"session_id":sessionID,
		"user_id": userID,
		"iat":     now.Unix(),
		"exp":     now.Add(30 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = "supersecretrefreshkey"
	}

	return token.SignedString([]byte(secret))
}

func SignUpToken()string{ 
	byte:=make([]byte, 32)
	_,err:=rand.Read(byte)
	if err!=nil{
		fmt.Println("Cant able to convert bytes into random numbers")
	}
	return hex.EncodeToString(byte)
}

func ReseTToken() string{
	byte:=make([]byte, 32)
	_,err:=rand.Read(byte)
	if err!=nil{
		fmt.Println("Cant able to convert bytes into random numbers")
	}
	return hex.EncodeToString(byte)
}