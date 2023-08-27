package auth_apis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/akshitbansal-1/async-testing/be/app"
	"github.com/golang-jwt/jwt"
)

const (
	GOOGLE_CLIENT_URL = "https://www.googleapis.com/oauth2/v1/certs"
	GOOGLE_ISSUER1 = "accounts.google.com"
	GOOGLE_ISSUER2 = "https://accounts.google.com"
)

type Service interface {
	LoginUser(tokenString string) (GoogleClaims, error)
}

type service struct {
	app app.App
}

func NewService(app app.App) Service {
	return &service{app}
}

func (s *service) LoginUser(tokenString string) (GoogleClaims, error) {
	return ValidateGoogleJWT(s, tokenString)
}

func getGooglePublicKey(keyID string) (string, error) {
	resp, err := http.Get(GOOGLE_CLIENT_URL)
	if err != nil {
		return "", err
	}
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	myResp := map[string]string{}
	err = json.Unmarshal(dat, &myResp)
	if err != nil {
		return "", err
	}
	key, ok := myResp[keyID]
	if !ok {
		return "", errors.New("key not found")
	}
	return key, nil
}

func ValidateGoogleJWT(s *service, tokenString string) (GoogleClaims, error) {
	cfg := s.app.GetConfig()
	claimsStruct := GoogleClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&claimsStruct,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := getGooglePublicKey(fmt.Sprintf("%s", token.Header["kid"]))
			if err != nil {
				return nil, err
			}
			key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(pem))
			if err != nil {
				return nil, err
			}
			return key, nil
		},
	)
	if err != nil {
		return GoogleClaims{}, err
	}

	claims, ok := token.Claims.(*GoogleClaims)
	if !ok {
		return GoogleClaims{}, errors.New("Invalid Google JWT")
	}

	if claims.Issuer != GOOGLE_ISSUER1 && claims.Issuer != GOOGLE_ISSUER2 {
		return GoogleClaims{}, errors.New("iss is invalid")
	}

	if claims.Audience != cfg.GOOGLE_SIGN_IN_CLIENT_ID {
		return GoogleClaims{}, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return GoogleClaims{}, errors.New("JWT is expired")
	}

	return *claims, nil
}
