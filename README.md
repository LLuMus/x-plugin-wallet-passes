# x-plugin-wallet-passes

A PoC of a Plugin for ChatGPT to generate iOS Wallet Passes, this PoC is running in production at https://walletpasses.xyz/ and you can find it on the ChatGPT Plugin Store.

This code base offers:
- A REST API to generate passes with an approved manifest for ChatGPT
- OAuth2 flow with Google Auth for user authentication
- Stripe integration for payments and acquisition of "tokens"
- A Webhook to receive Stripe events and update the user's tokens
- Generation of Apple Wallet Passes based on ChatGPT calls

## Requirements

```yaml
- OPENAI_PLUGIN_KEY=${OPENAI_PLUGIN_KEY}
- GOOGLE_CLIENT_ID=${GOOGLE_CLIENT_ID}
- TEAM_IDENTIFIER=${TEAM_IDENTIFIER}
- AWS_BUCKET_NAME=${AWS_BUCKET_NAME}
- AWS_REGION=${AWS_REGION}
- AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}
- AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}
- WALLET_STRIPE_WEBHOOK_SECRET=${WALLET_STRIPE_WEBHOOK_SECRET}
- WALLET_STRIPE_SECRET=${WALLET_STRIPE_SECRET}
- WALLET_STRIPE_PRICE=${WALLET_STRIPE_PRICE}
- WALLET_STRIPE_TAX=${WALLET_STRIPE_TAX}
```

Considering the following environment variables that we have to configure first, you can already see the list of services that we will have to prepare first:
- OpenAI Developer Account https://platform.openai.com/login?launch
- Google Developer Account https://console.developers.google.com/
  - This will be used for the Google Auth process
- AWS S3 https://aws.amazon.com/pm/serv-s3
- Stripe https://stripe.com/
- Apple Developer Account https://developer.apple.com/
  - After you have an account, you have to export and prepare a pass.p12 that will be used to sign the passes, you can follow this guide https://support.airship.com/hc/en-us/articles/213493683-How-to-make-an-Apple-Pass-Type-Certificate-for-Mobile-Wallet

## Run

```bash
$ docker-compose up
```

## Test

```bash
$ go test ./...
```

No tests written yet 👹
