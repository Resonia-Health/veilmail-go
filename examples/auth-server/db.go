package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"sync"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// In-memory store for simplicity — swap for a real DB in production.

type User struct {
	ID               int       `json:"id"`
	Email            string    `json:"email"`
	Name             string    `json:"name"`
	HashedPassword   string    `json:"-"`
	EmailVerified    bool      `json:"email_verified"`
	TwoFactorEnabled bool      `json:"two_factor_enabled"`
	CreatedAt        time.Time `json:"created_at"`
}

type Token struct {
	Email     string
	Token     string
	ExpiresAt time.Time
}

type TwoFactorCode struct {
	Email     string
	Code      string
	ExpiresAt time.Time
}

var (
	mu                  sync.Mutex
	users               = map[string]*User{}  // email -> user
	usersByID           = map[int]*User{}
	nextID              = 1
	verificationTokens  = map[string]*Token{} // token -> Token
	passwordResetTokens = map[string]*Token{}
	twoFactorCodes      = map[string]*TwoFactorCode{} // email -> code
)

func initDB() {
	// Nothing to initialize for in-memory store
}

func createUser(email, name, password string) (*User, error) {
	mu.Lock()
	defer mu.Unlock()

	if _, exists := users[email]; exists {
		return nil, fmt.Errorf("email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	u := &User{
		ID:             nextID,
		Email:          email,
		Name:           name,
		HashedPassword: string(hash),
		CreatedAt:      time.Now(),
	}
	nextID++
	users[email] = u
	usersByID[u.ID] = u
	return u, nil
}

func findUserByEmail(email string) *User {
	mu.Lock()
	defer mu.Unlock()
	return users[email]
}

func findUserByID(id int) *User {
	mu.Lock()
	defer mu.Unlock()
	return usersByID[id]
}

func checkPassword(user *User, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password)) == nil
}

func generateToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func generate2FACode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

func storeVerificationToken(email string) string {
	mu.Lock()
	defer mu.Unlock()
	t := generateToken()
	verificationTokens[t] = &Token{Email: email, Token: t, ExpiresAt: time.Now().Add(time.Hour)}
	return t
}

func storePasswordResetToken(email string) string {
	mu.Lock()
	defer mu.Unlock()
	t := generateToken()
	passwordResetTokens[t] = &Token{Email: email, Token: t, ExpiresAt: time.Now().Add(time.Hour)}
	return t
}

func store2FACode(email string) string {
	mu.Lock()
	defer mu.Unlock()
	code := generate2FACode()
	twoFactorCodes[email] = &TwoFactorCode{Email: email, Code: code, ExpiresAt: time.Now().Add(5 * time.Minute)}
	return code
}
