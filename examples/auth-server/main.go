package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	loadEnv()
	initDB()

	mux := http.NewServeMux()

	// Auth routes
	mux.HandleFunc("POST /auth/register", handleRegister)
	mux.HandleFunc("GET /auth/verify-email", handleVerifyEmail)
	mux.HandleFunc("POST /auth/login", handleLogin)
	mux.HandleFunc("POST /auth/verify-2fa", handleVerify2FA)
	mux.HandleFunc("POST /auth/forgot-password", handleForgotPassword)
	mux.HandleFunc("POST /auth/reset-password", handleResetPassword)

	// Protected routes
	mux.HandleFunc("GET /users/me", withAuth(handleGetMe))
	mux.HandleFunc("POST /users/toggle-2fa", withAuth(handleToggle2FA))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{
			"message": "VeilMail + Go Auth Example",
			"docs":    "See README.md for API endpoints",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
