package auth

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/jwtauth"
	"github.com/llumus/x-plugin-wallet-passes/internal/storage"
)

const ctxKey = "auth_user_email"

func Verifier(db storage.Database, cache storage.Cache) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authCode := jwtauth.TokenFromHeader(r)
			if authCode == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			log.Infof("[Verifier] authCode: %v", authCode)

			email, err := cache.Get(r.Context(), authCode)
			if err != nil || email == "" {
				log.Errorf("[Verifier] error getting email from cache: %v", err)

				c, err := db.GetCustomerByAuthCode(authCode)
				if err != nil {
					log.Errorf("[Verifier] error getting customer by auth code: %v", err)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				if c.Email == "" {
					log.Errorf("[Verifier] customer email is empty")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				err = cache.Set(r.Context(), authCode, c.Email, defaultExpiration)
				if err != nil {
					log.Errorf("[Verifier] error setting email in cache: %v", err)
				}

				email = c.Email
			}

			ctx := context.WithValue(r.Context(), ctxKey, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func EmailFromContext(ctx context.Context) (string, error) {
	email, ok := ctx.Value(ctxKey).(string)
	if !ok {
		return "", fmt.Errorf("could not get email from context")
	}

	return email, nil
}
