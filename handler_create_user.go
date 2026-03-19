package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mosgizy5/auth-app/internal/database"
	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type User struct{
		FullName string `json:"full_name"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	user_data := User{}
	
	err := json.NewDecoder(r.Body).Decode(&user_data)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("Error parsing JSON: %v", err))	
		return	
	}

	if len(user_data.Password) < 8 {
		respondWithError(w,400,"Password must be at least 8 characters")
		return
	}

	if err := IsValidEmail(user_data.Email); err != nil {
		respondWithError(w,400,fmt.Sprintf("Invalid email address: %v", err))	
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user_data.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("internal error: %v",err))
		return
	}

	user, err := apiCfg.DB.CreateUser(r.Context(),database.CreateUserParams{
		FullName: user_data.FullName,
		Email: strings.ToLower(user_data.Email),
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		respondWithError(w,400, fmt.Sprintf("Error creating user: %v", err))
		return
	}

	type UserT struct{
		ID string `json:"id"`
		Email string `json:"email"`
		FullName string `json:"full_name"`
		Token string `json:"token"`
	}

	token, err := GenerateToken(user.ID.String(), user.Email)
	if err != nil {
		respondWithError(w,400, fmt.Sprintf("Failed to generate token: %v", err))
		return
	}

	userReturn := UserT{
		ID: user.ID.String(),
		Email: user.Email,
		FullName: user.FullName,
		Token:token,
	}
	
	respondWithJSON(w,201,userReturn)
}