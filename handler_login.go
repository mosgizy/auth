package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (apiCfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct{
		Email string `json:"email"`
		Password string `json:"password"`
	}

	param := parameters{}

	json.NewDecoder(r.Body).Decode(&param)

	if err := IsValidEmail(param.Email); err != nil {
		respondWithError(w,400,fmt.Sprintf("email address is not a valid type of email: %v",err))
		return
	}

	user, err := apiCfg.DB.GetUserByEmail(r.Context(), param.Email)
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("Invalid email or password: %v",err))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(param.Password))
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("Invalid email or password: %v",err))
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

	respondWithJSON(w,200,userReturn)
}