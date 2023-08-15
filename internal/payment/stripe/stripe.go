package stripe

import (
	"github.com/sirupsen/logrus"
	"github.com/stripe/stripe-go/v72"
	stripeSession "github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/webhook"
)

type StripePayment struct {
	stripeWebhookSecret string
	stripeSecret        string
	stripePrice         string
	stripeTax           string
	baseUrl             string
}

var log = logrus.New()

func NewStripePayment(stripeWebhookSecret, stripeSecret, stripePrice, stripeTax, baseUrl string) *StripePayment {
	log.Debugf("[Stripe] Init with webhook: %s, secret: %s, price: %s, tax: %s, baseUrl: %s",
		stripeWebhookSecret, stripeSecret, stripePrice, stripeTax, baseUrl)

	stripe.Key = stripeSecret

	return &StripePayment{
		stripeWebhookSecret: stripeWebhookSecret,
		stripeSecret:        stripeSecret,
		stripePrice:         stripePrice,
		stripeTax:           stripeTax,
		baseUrl:             baseUrl,
	}
}

func (p *StripePayment) CreateSession(email, locale string) (string, error) {
	params := &stripe.CheckoutSessionParams{
		CustomerEmail: stripe.String(email),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Locale: stripe.String(locale),
		Mode:   stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(p.stripePrice),
				Quantity: stripe.Int64(1),
				TaxRates: stripe.StringSlice([]string{
					p.stripeTax,
				}),
			},
		},
		SuccessURL: stripe.String(p.baseUrl + "credit?success=true"),
		CancelURL:  stripe.String(p.baseUrl + "credit"),
	}

	session, err := stripeSession.New(params)
	if err != nil {
		return "", err
	}

	return session.ID, nil
}

func (p *StripePayment) ConstructEvent(payload []byte, header string) (stripe.Event, error) {
	return webhook.ConstructEvent(payload, header, p.stripeWebhookSecret)
}
