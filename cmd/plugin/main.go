package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/llumus/x-plugin-wallet-passes/internal/api"
	"github.com/llumus/x-plugin-wallet-passes/internal/auth"
	"github.com/llumus/x-plugin-wallet-passes/internal/components"
	"github.com/llumus/x-plugin-wallet-passes/internal/google"
	"github.com/llumus/x-plugin-wallet-passes/internal/logger"
	"github.com/llumus/x-plugin-wallet-passes/internal/pass/apple"
	"github.com/llumus/x-plugin-wallet-passes/internal/payment/stripe"
	"github.com/llumus/x-plugin-wallet-passes/internal/storage/pg"
	"github.com/llumus/x-plugin-wallet-passes/internal/storage/redis"
	"github.com/llumus/x-plugin-wallet-passes/internal/storage/s3"
	"github.com/sirupsen/logrus"
)

const (
	shutdownDeadline = time.Second * 30
)

var log = logrus.New()

const baseManifest = `{
  "schema_version": "v1",
  "name_for_human": "Wallet Passes",
  "name_for_model": "wallet_passes",
  "description_for_human": "Create Wallet Passes for iOS device (iPhone, iPod, iWatch) with ease!",
  "description_for_model": "Create a iOS Wallet Pass (Passbook) and returns a link to visit and add it on your iOS device (iPhone, iPad, iWatch).",
  "auth": {
    "type": "oauth",
    "client_url": "${baseUrl}welcome",
    "scope": "",
    "authorization_url": "${baseUrl}auth/open-ai",
    "authorization_content_type": "application/json",
    "verification_tokens": {
	  "openai": "${openaiToken}"
    }
  },
  "api": {
    "type": "openapi",
    "url": "${baseUrl}static/openapi.yml"
  },
  "logo_url": "https://storage.googleapis.com/walletpasses/logo_square.png",
  "contact_email": "contact@walletpasses.xyz",
  "legal_info_url": "${baseUrl}legal"
}`

func serveManifest(baseUrl, openAIPluginKey string) func(w http.ResponseWriter, r *http.Request) {
	var readyManifest = strings.ReplaceAll(baseManifest, "${baseUrl}", baseUrl)
	readyManifest = strings.ReplaceAll(readyManifest, "${openaiToken}", openAIPluginKey)

	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(readyManifest))
		w.WriteHeader(http.StatusOK)
	}
}

func main() {
	var port = os.Getenv("PORT")
	var redisAddress = os.Getenv("REDIS_ADDRESS")
	var redisPassword = os.Getenv("REDIS_PASSWORD")
	var redisTLS = os.Getenv("REDIS_TLS")
	var openAIPluginKey = os.Getenv("OPENAI_PLUGIN_KEY")
	var googleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	var teamIdentifier = os.Getenv("TEAM_IDENTIFIER")
	var basePath = os.Getenv("BASE_PATH")
	var bucketName = os.Getenv("AWS_BUCKET_NAME")
	var stripeWebhookSecret = os.Getenv("WALLET_STRIPE_WEBHOOK_SECRET")
	var stripeSecret = os.Getenv("WALLET_STRIPE_SECRET")
	var stripePrice = os.Getenv("WALLET_STRIPE_PRICE")
	var stripeTax = os.Getenv("WALLET_STRIPE_TAX")
	var baseUrl = os.Getenv("BASE_URL")

	var db = pg.NewDB()
	var cache = redis.NewCache(redisAddress, redisPassword, redisTLS)
	var pass = apple.NewPassbook(baseUrl, teamIdentifier, basePath, filepath.Join(basePath, "tmp"))
	var fs = s3.NewS3Storage(bucketName)
	var payment = stripe.NewStripePayment(stripeWebhookSecret, stripeSecret, stripePrice, stripeTax, baseUrl)
	var googleAuth = google.NewTokenValidator(googleClientID, cache)
	var currentApi = api.NewAPI(baseUrl, pass, fs, db, payment, cache)
	var fsPublic = http.FileServer(http.Dir("./public"))

	log.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: true,
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestLogger(&logger.StructuredLogger{Logger: log}))
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.Handler(cors.Options{
		Debug:            false,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Get("/.well-known/ai-plugin.json", serveManifest(baseUrl, openAIPluginKey))
	r.Handle("/static/*", fsPublic)

	r.Post("/auth/open-ai", auth.OpenAIOAuth(cache))
	r.Post("/auth/google", auth.GoogleOAuth(googleAuth, db, cache))
	r.Get("/api/v1/pass/{id}", currentApi.GetPass)
	r.Post("/api/v1/webhook", currentApi.WebhookStripe)

	r.Group(func(r chi.Router) {
		r.Use(auth.CookieIntoHeader())
		r.Use(auth.Verifier(db, cache))
		r.Get("/api/v1/me", currentApi.Me)
		r.Get("/api/v1/history", currentApi.History)
		r.Get("/api/v1/checkout", currentApi.Checkout)

		handler := components.HandlerWithOptions(components.NewStrictHandler(currentApi, nil), components.ChiServerOptions{
			BaseRouter: r,
		})

		r.Handle("/api/v1/*", http.StripPrefix("/api/v1", handler))
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(basePath, "/public/index.html"))
	})

	httpSrv := http.Server{
		Addr:    net.JoinHostPort("0.0.0.0", port),
		Handler: r,
	}

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
		<-signals

		ctx, cancel := context.WithTimeout(context.Background(), shutdownDeadline)
		defer cancel()

		log.Info("shutting down server..")
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Errorf("failed to shutdown %s", err)
		}
	}()

	err := httpSrv.ListenAndServe()
	if err != nil {
		log.Fatalf("failed to start server %s", err)
	}
}
