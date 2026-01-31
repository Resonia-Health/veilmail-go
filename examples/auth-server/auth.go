package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const userKey contextKey = "user"

func createJWT(user *User) string {
	secret := os.Getenv("SECRET_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   fmt.Sprintf("%d", user.ID),
		"email": user.Email,
		"exp":   time.Now().Add(30 * time.Minute).Unix(),
	})
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func withAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			writeError(w, 401, "Missing authorization header")
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		secret := os.Getenv("SECRET_KEY")

		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			writeError(w, 401, "Invalid token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			writeError(w, 401, "Invalid token claims")
			return
		}

		sub, _ := claims.GetSubject()
		id, _ := strconv.Atoi(sub)
		user := findUserByID(id)
		if user == nil {
			writeError(w, 401, "User not found")
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		next(w, r.WithContext(ctx))
	}
}
