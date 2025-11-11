package middleware

import (
	"context"
	"fmt"
	"net/http"
	"pickel-backend/utils"
)

type contextKey string

const userContextKey = contextKey("userClaims")

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		var tokenStr string
		fmt.Println("cookie: ", cookie)
		if err == nil {
			tokenStr = cookie.Value
		}

		if tokenStr == "" {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ValidateJWT(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized (Middleware): "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(r *http.Request) (*utils.Claims, bool) {
	claims, ok := r.Context().Value(userContextKey).(*utils.Claims)

	return claims, ok
}
