package storage

import (
	"context"
	"io"
	"time"
)

type Pass struct {
	Id            string `json:"id"`
	CustomerEmail string `json:"customer_email"`
	FileReference string `json:"file_reference"`
	Payload       string `json:"payload"`
	CreatedAt     int64  `json:"created_at"`
}

type Customer struct {
	Email        string `json:"email"`
	Details      string `json:"details"`
	CreditTokens int64  `json:"credit_tokens"`
	AuthCode     string `json:"auth_code"`
	UpdatedAt    int64  `json:"updated_at"`
	CreatedAt    int64  `json:"created_at"`
}

type FileSystem interface {
	Save(key, filePath, contentType string) (string, error)
	SaveFile(key string, buffer io.Reader, contentType string, contentLength int64) (string, error)
	GenerateSignedUrl(key string, expires time.Duration) (string, error)
}

type Database interface {
	InsertPass(*Pass) error
	GetPassbookById(string) (*Pass, error)
	GetPassbooksByCustomerEmail(string) ([]*Pass, error)
	GetCustomerByEmail(string) (*Customer, error)
	GetCustomerByAuthCode(string) (*Customer, error)
	LoginCustomer(*Customer) (*Customer, error)
	UseToken(email string) (int, error)
	CreditCustomer(string, int) (int, error)
}

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
}
