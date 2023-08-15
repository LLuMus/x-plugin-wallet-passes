package payment

import "github.com/stripe/stripe-go/v72"

type Payment interface {
	CreateSession(email, locale string) (string, error)
	ConstructEvent(payload []byte, header string) (stripe.Event, error)
}
