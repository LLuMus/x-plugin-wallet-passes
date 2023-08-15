package google

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/llumus/x-plugin-wallet-passes/internal/storage"
)

type JWTValidator interface {
	ValidateGoogleJWT(ctx context.Context, tokenString string) (*TokenClaims, error)
}

type TokenClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	FirstName     string `json:"given_name"`
	LastName      string `json:"family_name"`
	jwt.StandardClaims
}

type TokenValidator struct {
	clientId string
	cache    storage.Cache
}

const publicKeyCacheKey = "google_public_key"

func NewTokenValidator(clientId string, cache storage.Cache) *TokenValidator {
	return &TokenValidator{
		clientId: clientId,
		cache:    cache,
	}
}

func (v *TokenValidator) getPublicPem(ctx context.Context, keyId string) (string, error) {
	cacheKey, err := v.cache.Get(ctx, publicKeyCacheKey)
	if err == nil && cacheKey != "" {
		return cacheKey, nil
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/certs")
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

	key, ok := myResp[keyId]
	if !ok {
		return "", errors.New("key not found")
	}

	_ = v.cache.Set(ctx, publicKeyCacheKey, key, 24*time.Hour)

	return key, nil
}

func (v *TokenValidator) ValidateGoogleJWT(ctx context.Context, tokenString string) (*TokenClaims, error) {
	baseClaims := &TokenClaims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		baseClaims,
		func(token *jwt.Token) (interface{}, error) {
			pem, err := v.getPublicPem(ctx, fmt.Sprintf("%s", token.Header["kid"]))
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
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("invalid Google JWT")
	}

	if claims.Issuer != "accounts.google.com" && claims.Issuer != "https://accounts.google.com" {
		return nil, errors.New("iss is invalid")
	}

	if claims.Audience != v.clientId {
		return nil, errors.New("aud is invalid")
	}

	if claims.ExpiresAt < time.Now().UTC().Unix() {
		return nil, errors.New("JWT is expired")
	}

	return claims, nil
}
