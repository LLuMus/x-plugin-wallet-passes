package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/llumus/x-plugin-wallet-passes/internal/auth"
	"github.com/llumus/x-plugin-wallet-passes/internal/components"
	"github.com/llumus/x-plugin-wallet-passes/internal/pass"
	"github.com/llumus/x-plugin-wallet-passes/internal/payment"
	"github.com/llumus/x-plugin-wallet-passes/internal/storage"
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72"
)

type Api struct {
	baseUrl  string
	passbook pass.Passbook
	fs       storage.FileSystem
	db       storage.Database
	payment  payment.Payment
	cache    storage.Cache
}

var log = logrus.New()

func NewAPI(baseUrl string, passbook pass.Passbook, fs storage.FileSystem, db storage.Database, payment payment.Payment, cache storage.Cache) *Api {
	return &Api{
		baseUrl:  baseUrl,
		passbook: passbook,
		fs:       fs,
		db:       db,
		payment:  payment,
		cache:    cache,
	}
}

func (a *Api) GetPass(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	passbook, err := a.db.GetPassbookById(id.String())
	if err != nil {
		log.Errorf("[Redeem] db fetch error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if passbook.Id == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result, err := json.Marshal(passbook)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) Checkout(w http.ResponseWriter, r *http.Request) {
	email, err := auth.EmailFromContext(r.Context())
	if err != nil || email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	sessionId, err := a.payment.CreateSession(email, "en")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"sessionId": "` + sessionId + `"}`))
	w.WriteHeader(http.StatusOK)
}

func (a *Api) Me(w http.ResponseWriter, r *http.Request) {
	email, err := auth.EmailFromContext(r.Context())
	if err != nil || email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	customer, err := a.db.GetCustomerByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if customer.Email == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	result, err := json.Marshal(customer)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) History(w http.ResponseWriter, r *http.Request) {
	email, err := auth.EmailFromContext(r.Context())
	if err != nil || email == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	passes, err := a.db.GetPassbooksByCustomerEmail(email)
	if err != nil {
		log.Errorf("[History] db fetch error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if passes == nil {
		passes = []*storage.Pass{}
	}

	result, err := json.Marshal(passes)
	if err != nil {
		log.Errorf("[History] json marshal error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
	w.WriteHeader(http.StatusOK)
}

func (a *Api) CreatePassbook(ctx context.Context, request components.CreatePassbookRequestObject) (components.CreatePassbookResponseObject, error) {
	email, err := auth.EmailFromContext(ctx)
	if err != nil || email == "" {
		log.Errorf("[CreatePassbook] error: %s", err.Error())
		return nil, err
	}

	link, err := a.CreatePassbookWithEmail(ctx, request, email)
	if err != nil {
		log.Errorf("[CreatePassbook] error: %s", err.Error())
		return nil, err
	}

	return &components.CreatePassbook200JSONResponse{
		Link: &link,
	}, nil
}

func (a *Api) CreatePassbookWithEmail(ctx context.Context, request components.CreatePassbookRequestObject, email string) (string, error) {
	customer, err := a.db.GetCustomerByEmail(email)
	if err != nil {
		log.Errorf("[CreatePassbook] db error: %s", err.Error())
		return "", err
	}

	if customer.Email == "" {
		log.Errorf("[CreatePassbook] customer not found")
		return "", err
	}

	if customer.CreditTokens <= 0 {
		return a.baseUrl + "credit", nil
	}

	id, path, err := a.passbook.CreatePass(ctx, request)
	if err != nil {
		log.Errorf("[CreatePassbook] passbook error: %s", err.Error())
		return "", err
	}

	fileKey, err := a.fs.Save(filepath.Join(id, "WalletPass_"+id+".pkpass"), path, "application/vnd.apple.pkpass")
	if err != nil {
		log.Errorf("[CreatePassbook] fs error: %s", err.Error())
		return "", err
	}

	payloadAsBytes, err := json.Marshal(request)
	if err != nil {
		log.Errorf("[CreatePassbook] json marshal error: %s", err.Error())
		return "", err
	}

	err = a.db.InsertPass(&storage.Pass{
		Id:            id,
		CustomerEmail: customer.Email,
		FileReference: fileKey,
		Payload:       string(payloadAsBytes),
		CreatedAt:     time.Now().Unix(),
	})
	if err != nil {
		log.Errorf("[CreatePassbook] db error: %s", err.Error())
		return "", err
	}

	_, err = a.db.UseToken(customer.Email)
	if err != nil {
		log.Errorf("[CreatePassbook] db error using token - moving forward: %s", err.Error())
	}

	return a.baseUrl + "pass/" + id, nil
}

func (a *Api) WebhookStripe(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Errorf("[WebhookStripe] error reading body: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	event, err := a.payment.ConstructEvent(body, r.Header.Get("Stripe-Signature"))
	if err != nil {
		log.Errorf("[WebhookStripe] error constructing event: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if event.Type == "checkout.session.completed" {
		log.Infof("[WebhookStripe] event data: %s", string(event.Data.Raw))

		var session stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &session)
		if err != nil {
			log.Errorf("[WebhookStripe] error unmarshling json: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Infof("[WebhookStripe] session customer email: %s", session.CustomerEmail)

		customer, err := a.db.GetCustomerByEmail(session.CustomerEmail)
		if err != nil {
			log.Errorf("[WebhookStripe] error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if customer.Email == "" {
			http.Error(w, "customer not found", http.StatusNotFound)
			return
		}

		totalTokens, err := a.db.CreditCustomer(customer.Email, 1000)
		if err != nil {
			log.Errorf("[WebhookStripe] error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = a.cache.Set(r.Context(), session.ID, strconv.Itoa(totalTokens), -1)
		if err != nil {
			log.Errorf("[WebhookStripe] error: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}
