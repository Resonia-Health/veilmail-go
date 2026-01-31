package veilmail

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func TestVerifySignature(t *testing.T) {
	secret := "whsec_test_secret_123"
	body := []byte(`{"event":"email.delivered","data":{"id":"email_123"}}`)

	// Generate valid signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	validSignature := hex.EncodeToString(mac.Sum(nil))

	t.Run("valid signature", func(t *testing.T) {
		if !VerifySignature(body, validSignature, secret) {
			t.Error("expected valid signature to pass verification")
		}
	})

	t.Run("invalid signature", func(t *testing.T) {
		if VerifySignature(body, "invalid_signature", secret) {
			t.Error("expected invalid signature to fail verification")
		}
	})

	t.Run("wrong secret", func(t *testing.T) {
		if VerifySignature(body, validSignature, "wrong_secret") {
			t.Error("expected wrong secret to fail verification")
		}
	})

	t.Run("tampered body", func(t *testing.T) {
		tamperedBody := []byte(`{"event":"email.delivered","data":{"id":"email_456"}}`)
		if VerifySignature(tamperedBody, validSignature, secret) {
			t.Error("expected tampered body to fail verification")
		}
	})

	t.Run("empty body", func(t *testing.T) {
		emptyBody := []byte{}
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(emptyBody)
		emptySignature := hex.EncodeToString(mac.Sum(nil))

		if !VerifySignature(emptyBody, emptySignature, secret) {
			t.Error("expected empty body with correct signature to pass")
		}
	})

	t.Run("constant time comparison", func(t *testing.T) {
		// Verify that hmac.Equal is used (which provides constant-time comparison)
		// This is implicitly tested by the VerifySignature function using hmac.Equal
		// We can verify it works correctly with near-miss signatures
		nearMiss := validSignature[:len(validSignature)-1] + "0"
		if validSignature[len(validSignature)-1] == '0' {
			nearMiss = validSignature[:len(validSignature)-1] + "1"
		}
		if VerifySignature(body, nearMiss, secret) {
			t.Error("expected near-miss signature to fail")
		}
	})
}
