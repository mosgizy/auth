package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mosgizy5/auth-app/internal/database"
	"golang.org/x/crypto/bcrypt"
)


	type messages struct{
		Message string `json:"message"`
	}

func (apiCfg *apiConfig) handlerForgetPassword(w http.ResponseWriter, r *http.Request){
	type parameters struct{
		Email string `json:"email"`
	}

	message := messages{
		Message: "If your email is registered, you will recieve a link.",
	}

	params := parameters{}

	json.NewDecoder(r.Body).Decode(&params)

	user, err := apiCfg.DB.GetUserByEmail(r.Context(),params.Email)
	if err != nil {
		respondWithJSON(w,200,message)
		return
	}

	resetToken := uuid.New().String()
	expiresAt := time.Now().Add(15 * time.Minute)

	err = apiCfg.DB.SetResetToken(r.Context(), database.SetResetTokenParams{
		Email: params.Email,
		ResetToken: sql.NullString{
			String: resetToken,
			Valid: true,
		},
		ResetTokenExpires: sql.NullTime{
			Time:expiresAt,
			Valid:true,
		},
	})

	resetLink := fmt.Sprintf("https://localhost:8080/password-reset?token=%s",resetToken)

	// TODO: Send email using Resend / SendGrid / Zoho
  // For now, just log it (replace with real email service)

	fmt.Printf("Reset link for %s: %s\n", user.Email, resetLink)

	respondWithJSON(w,200,message)
}

func (apiCfg *apiConfig) handlerResetPassword(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Token string `json:"token"`
		Password string `json:"password"`
	}

	params := parameters{}

	json.NewDecoder(r.Body).Decode(&params)

	user, err := apiCfg.DB.GetUserByResetToken(r.Context(),sql.NullString{
		String: params.Token,
		Valid: true,
	})
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("Invalid or expired rest token: %s",err))
		return
	}

	hashedPassword,_ := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)

	err = apiCfg.DB.UpdatePassword(r.Context(),database.UpdatePasswordParams{
		ID: user.ID,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		respondWithError(w,400,fmt.Sprintf("failed to update password: %v",err))
		return
	}

	message := messages{
		Message: "Password reset successful. You can now login.",
	}

	respondWithJSON(w,200,message)
}