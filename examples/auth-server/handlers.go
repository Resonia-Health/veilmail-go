package main

import (
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func handleRegister(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "Invalid request body")
		return
	}

	user, err := createUser(body.Email, body.Name, body.Password)
	if err != nil {
		writeError(w, 400, err.Error())
		return
	}

	token := storeVerificationToken(user.Email)
	sendVerificationEmail(user.Email, token)

	writeJSON(w, 201, map[string]string{"message": "Verification email sent"})
}

func handleVerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	mu.Lock()
	record, exists := verificationTokens[token]
	if exists {
		delete(verificationTokens, token)
	}
	mu.Unlock()

	if !exists {
		writeError(w, 400, "Invalid token")
		return
	}
	if time.Now().After(record.ExpiresAt) {
		writeError(w, 400, "Token expired")
		return
	}

	user := findUserByEmail(record.Email)
	if user == nil {
		writeError(w, 404, "User not found")
		return
	}

	mu.Lock()
	user.EmailVerified = true
	mu.Unlock()

	sendWelcomeEmail(user.Email, user.Name)
	writeJSON(w, 200, map[string]string{"message": "Email verified"})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "Invalid request body")
		return
	}

	user := findUserByEmail(body.Email)
	if user == nil || !checkPassword(user, body.Password) {
		writeError(w, 401, "Invalid credentials")
		return
	}

	if !user.EmailVerified {
		token := storeVerificationToken(user.Email)
		sendVerificationEmail(user.Email, token)
		writeError(w, 403, "Email not verified. Verification email resent.")
		return
	}

	if user.TwoFactorEnabled {
		code := store2FACode(user.Email)
		sendTwoFactorCode(user.Email, code)
		writeJSON(w, 200, map[string]any{"two_factor_required": true})
		return
	}

	token := createJWT(user)
	writeJSON(w, 200, map[string]string{"access_token": token, "token_type": "bearer"})
}

func handleVerify2FA(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Code     string `json:"code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "Invalid request body")
		return
	}

	user := findUserByEmail(body.Email)
	if user == nil || !checkPassword(user, body.Password) {
		writeError(w, 401, "Invalid credentials")
		return
	}

	mu.Lock()
	record, exists := twoFactorCodes[body.Email]
	if exists {
		delete(twoFactorCodes, body.Email)
	}
	mu.Unlock()

	if !exists || record.Code != body.Code {
		writeError(w, 400, "Invalid code")
		return
	}
	if time.Now().After(record.ExpiresAt) {
		writeError(w, 400, "Code expired")
		return
	}

	token := createJWT(user)
	writeJSON(w, 200, map[string]string{"access_token": token, "token_type": "bearer"})
}

func handleForgotPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "Invalid request body")
		return
	}

	user := findUserByEmail(body.Email)
	if user != nil {
		token := storePasswordResetToken(body.Email)
		sendPasswordResetEmail(body.Email, token)
	}

	writeJSON(w, 200, map[string]string{"message": "If an account exists, a reset email has been sent"})
}

func handleResetPassword(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, 400, "Invalid request body")
		return
	}

	mu.Lock()
	record, exists := passwordResetTokens[body.Token]
	if exists {
		delete(passwordResetTokens, body.Token)
	}
	mu.Unlock()

	if !exists {
		writeError(w, 400, "Invalid token")
		return
	}
	if time.Now().After(record.ExpiresAt) {
		writeError(w, 400, "Token expired")
		return
	}

	user := findUserByEmail(record.Email)
	if user == nil {
		writeError(w, 404, "User not found")
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 12)
	mu.Lock()
	user.HashedPassword = string(hash)
	mu.Unlock()

	sendPasswordChangedEmail(user.Email)
	writeJSON(w, 200, map[string]string{"message": "Password updated"})
}

func handleGetMe(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey).(*User)
	writeJSON(w, 200, user)
}

func handleToggle2FA(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userKey).(*User)

	mu.Lock()
	user.TwoFactorEnabled = !user.TwoFactorEnabled
	enabled := user.TwoFactorEnabled
	mu.Unlock()

	send2FAToggledEmail(user.Email, enabled)

	status := "disabled"
	if enabled {
		status = "enabled"
	}
	writeJSON(w, 200, map[string]any{
		"two_factor_enabled": enabled,
		"message":            "2FA " + status,
	})
}
