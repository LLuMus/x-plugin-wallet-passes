package auth

import (
	"encoding/json"
	"github.com/llumus/x-plugin-wallet-passes/internal/storage"
	"io"
	"net/http"
)

type OpenAIResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
}

type OpenAIRequest struct {
	Code         string `json:"code"`
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func OpenAIOAuth(cache storage.Cache) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		bodyRaw, _ := io.ReadAll(r.Body)
		log.Infof("[OpenAIOAuth] request body: %v", string(bodyRaw))

		var req = &OpenAIRequest{}
		err := json.Unmarshal(bodyRaw, req)
		if err != nil {
			log.Errorf("[OpenAIOAuth] error decoding request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		email, err := cache.Get(r.Context(), req.Code)
		if err != nil || email == "" {
			log.Errorf("[OpenAIOAuth] error getting email from cache: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		j, err := json.Marshal(&OpenAIResponse{
			AccessToken: req.Code,
			TokenType:   "bearer",
		})
		if err != nil {
			log.Errorf("[OpenAIOAuth] error marshalling response: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(j)
		w.WriteHeader(http.StatusOK)
	}
}
