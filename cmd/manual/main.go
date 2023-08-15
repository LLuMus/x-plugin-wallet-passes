package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/llumus/x-plugin-wallet-passes/internal/api"
	"github.com/llumus/x-plugin-wallet-passes/internal/components"
	"github.com/llumus/x-plugin-wallet-passes/internal/pass/apple"
	"github.com/llumus/x-plugin-wallet-passes/internal/payment/stripe"
	"github.com/llumus/x-plugin-wallet-passes/internal/storage/pg"
	"github.com/llumus/x-plugin-wallet-passes/internal/storage/redis"
	"github.com/llumus/x-plugin-wallet-passes/internal/storage/s3"
)

func main() {
	var email = os.Getenv("WALLET_MANUAL_EMAIL")
	var filePath = os.Getenv("FILE_PATH")
	var redisAddress = os.Getenv("REDIS_ADDRESS")
	var redisPassword = os.Getenv("REDIS_PASSWORD")
	var redisTLS = os.Getenv("REDIS_TLS")
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
	var currentApi = api.NewAPI(baseUrl, pass, fs, db, payment, cache)

	var ctx = context.Background()

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		os.Exit(1)
	}

	var body = &components.CreatePassbookRequest{}
	err = json.Unmarshal(data, body)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		os.Exit(1)
	}

	var req = &components.CreatePassbookRequestObject{
		Body: body,
	}

	link, err := currentApi.CreatePassbookWithEmail(ctx, *req, email)
	if err != nil {
		fmt.Println("Error creating pass:", err)
		os.Exit(1)
	}

	fmt.Println("Pass created successfully at link: " + link)
}
