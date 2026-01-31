# Go Auth Server Example with VeilMail

A net/http authentication server demonstrating email verification, password reset, two-factor authentication, and security notifications using the VeilMail Go SDK.

## Features

- **Email verification** on registration
- **Password reset** via email link
- **Two-factor authentication** with emailed codes
- **Security notifications** for password changes and 2FA toggles
- JWT-based session tokens

## Setup

```bash
cp .env.example .env
# Edit .env with your VeilMail API key and settings
go run .
```

## API Endpoints

| Method | Path                   | Description                  |
|--------|------------------------|------------------------------|
| POST   | /auth/register         | Register a new account       |
| GET    | /auth/verify-email     | Verify email via token       |
| POST   | /auth/login            | Log in (returns JWT or 2FA)  |
| POST   | /auth/verify-2fa       | Complete 2FA login           |
| POST   | /auth/forgot-password  | Request password reset email |
| POST   | /auth/reset-password   | Reset password with token    |
| GET    | /users/me              | Get current user (auth required) |
| POST   | /users/toggle-2fa      | Toggle 2FA (auth required)   |
