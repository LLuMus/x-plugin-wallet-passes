package pg

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
	"github.com/llumus/x-plugin-wallet-passes/internal/storage"
)

type DB struct {
	db *sql.DB
}

func NewDB() *DB {
	var (
		err     error
		connURI = os.Getenv("DB_URI")
	)

	c, err := sql.Open("postgres", connURI)
	if err != nil {
		panic(err)
	}

	return &DB{
		db: c,
	}
}

func (d *DB) InsertPass(pass *storage.Pass) error {
	var sqlStr = `INSERT INTO wallet_passes (id, customer_email, payload, file_reference, created_at) 
		VALUES ($1, $2, $3, $4, $5)`

	_, err := d.db.Exec(sqlStr, pass.Id, pass.CustomerEmail, pass.Payload, pass.FileReference, pass.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) GetPassbookById(id string) (*storage.Pass, error) {
	sqlStr := `SELECT id, customer_email, file_reference, payload, created_at FROM wallet_passes WHERE id = $1`

	rows, err := d.db.Query(sqlStr, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var n = &storage.Pass{}
	for rows.Next() {
		err := rows.Scan(
			&n.Id,
			&n.CustomerEmail,
			&n.FileReference,
			&n.Payload,
			&n.CreatedAt)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return n, nil
}

func (d *DB) GetPassbooksByCustomerEmail(email string) ([]*storage.Pass, error) {
	sqlStr := `SELECT id, customer_email, payload, created_at FROM wallet_passes WHERE customer_email = $1 ORDER BY created_at DESC`

	rows, err := d.db.Query(sqlStr, email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var passes []*storage.Pass
	for rows.Next() {
		var n = &storage.Pass{}
		err := rows.Scan(
			&n.Id,
			&n.CustomerEmail,
			&n.Payload,
			&n.CreatedAt)
		if err != nil {
			return nil, err
		}
		passes = append(passes, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return passes, nil
}

func (d *DB) GetCustomerByAuthCode(authToken string) (*storage.Customer, error) {
	var sqlStr = `SELECT email, details, credit_tokens, auth_code, updated_at, created_at FROM customers WHERE auth_code = $1`

	rows, err := d.db.Query(sqlStr, authToken)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var n = &storage.Customer{}
	for rows.Next() {
		err := rows.Scan(
			&n.Email,
			&n.Details,
			&n.CreditTokens,
			&n.AuthCode,
			&n.UpdatedAt,
			&n.CreatedAt)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return n, nil
}

func (d *DB) GetCustomerByEmail(email string) (*storage.Customer, error) {
	var sqlStr = `SELECT email, details, credit_tokens, auth_code, updated_at, created_at FROM customers WHERE email = $1`

	rows, err := d.db.Query(sqlStr, email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var n = &storage.Customer{}
	for rows.Next() {
		err := rows.Scan(
			&n.Email,
			&n.Details,
			&n.CreditTokens,
			&n.AuthCode,
			&n.UpdatedAt,
			&n.CreatedAt)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return n, nil
}

func (d *DB) LoginCustomer(customer *storage.Customer) (*storage.Customer, error) {
	var sqlStr = `INSERT INTO customers (email, details, credit_tokens, auth_code, updated_at, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (email) DO UPDATE SET details = $2, updated_at = $5 
		RETURNING email, details, credit_tokens, auth_code, updated_at, created_at`

	rows, err := d.db.Query(sqlStr, customer.Email, customer.Details, customer.CreditTokens, customer.AuthCode, customer.UpdatedAt, customer.CreatedAt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var n = &storage.Customer{}
	for rows.Next() {
		err := rows.Scan(
			&n.Email,
			&n.Details,
			&n.CreditTokens,
			&n.AuthCode,
			&n.UpdatedAt,
			&n.CreatedAt)
		if err != nil {
			return nil, err
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return n, nil
}

func (d *DB) UseToken(email string) (int, error) {
	sqlStr := `UPDATE customers SET credit_tokens = credit_tokens - 1 WHERE email = $1 RETURNING credit_tokens`

	rows, err := d.db.Query(sqlStr, email)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var totalCredit = 0
	for rows.Next() {
		err := rows.Scan(&totalCredit)
		if err != nil {
			return 0, err
		}
	}

	return totalCredit, nil
}

func (d *DB) CreditCustomer(email string, tokens int) (int, error) {
	sqlStr := `UPDATE customers SET credit_tokens = credit_tokens + $1 WHERE email = $2 RETURNING credit_tokens`

	rows, err := d.db.Query(sqlStr, tokens, email)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	var totalCredit = 0
	for rows.Next() {
		err := rows.Scan(&totalCredit)
		if err != nil {
			return 0, err
		}
	}

	return totalCredit, nil
}
