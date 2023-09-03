package auth_apis

import "github.com/golang-jwt/jwt"

type GoogleClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

type AuthRequest struct {
	IdToken string `json:"idToken"`
}

type ClientToken struct {
	Token    string `json:"token"`
	TokenExp string `json:"tokenExp"`
}
