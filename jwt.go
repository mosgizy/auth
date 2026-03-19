package main

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ID string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

var jwtStr = os.Getenv("jwt_secret")
var jwtSecret = []byte(jwtStr)

func GenerateToken(id,email string) (string,error){
	claims := Claims{
		ID: id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
      Issuer:    "asterisk",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)

	return token.SignedString(jwtSecret)
}

func ValidateToken(tokenStr string) (*Claims, error){
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func (token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok{
			return nil, errors.New("Unexpected signing method")
		}

		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}