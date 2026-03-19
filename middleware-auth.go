package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mosgizy5/auth-app/internal/auth"
)

// type authHandler func(http.ResponseWriter,http.Request,database.User)

type contextKey string
const UserContextKey contextKey = "user"

func (apiCfg *apiConfig) authMiddleWare(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
		token,err := auth.GetToken(r.Header)
		if err != nil {
			respondWithError(w,400,fmt.Sprintf("%v",err))
			return
		}

		claims, err := ValidateToken(token)
		if err != nil {
			respondWithError(w,400,fmt.Sprintf("Invalid or expired token: %v",err))
			return
		}

		ctx := context.WithValue(r.Context(),UserContextKey,claims)
		handler.ServeHTTP(w,r.WithContext(ctx))
	})
}