package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/llumus/x-plugin-wallet-passes/internal/google"
	"github.com/llumus/x-plugin-wallet-passes/internal/storage"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

const defaultExpiration = 15 * 24 * 60 * time.Minute // 15 days

func GoogleOAuth(googleJWTValidator google.JWTValidator, db storage.Database, cache storage.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		type parameters struct {
			GoogleJWT string `json:"access_token"`
		}

		decoder := json.NewDecoder(r.Body)
		params := parameters{}
		err := decoder.Decode(&params)
		if err != nil || params.GoogleJWT == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		claims, err := googleJWTValidator.ValidateGoogleJWT(r.Context(), params.GoogleJWT)
		if err != nil || claims == nil || claims.Email == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		details, err := json.Marshal(claims)
		if err != nil {
			log.Errorf("[GoogleAuth] error marshalling claims: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var customer = &storage.Customer{
			Email:        claims.Email,
			Details:      string(details),
			CreditTokens: 20,
			AuthCode:     uuid.New().String(),
			UpdatedAt:    time.Now().Unix(),
			CreatedAt:    time.Now().Unix(),
		}

		logged, err := db.LoginCustomer(customer)
		if err != nil {
			log.Errorf("[GoogleAuth] error logging in customer: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = cache.Set(r.Context(), logged.AuthCode, logged.Email, defaultExpiration)
		if err != nil {
			log.Errorf("[GoogleAuth] error setting cache: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		result, err := json.Marshal(logged)
		if err != nil {
			log.Errorf("[GoogleAuth] error marshalling customer: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     cookieName,
			Value:    logged.AuthCode,
			Expires:  time.Now().Add(defaultExpiration),
			HttpOnly: false,
			SameSite: http.SameSiteStrictMode,
			Secure:   true,
		})

		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
		w.WriteHeader(http.StatusOK)
	}
}
