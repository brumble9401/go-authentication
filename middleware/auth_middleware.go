package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gocql/gocql"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
    UserID   gocql.UUID `json:"user_id"`
    // Username string     `json:"username"`
    jwt.RegisteredClaims
}

var jwtKey = []byte("5bpehDpA1N0Hj1o+4piTXnRiJVosa9ND7n3QhBZR/cw=")

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Authorization header is required", http.StatusUnauthorized)
            return
        }

        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims := &Claims{}

        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            return jwtKey, nil
        })
        log.Debug("Token: ", token)
        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), "userID", claims.UserID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}