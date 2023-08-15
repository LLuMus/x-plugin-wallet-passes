CREATE TABLE customers (
   email character varying PRIMARY KEY,
   details JSONB NOT NULL,
   auth_code character varying NOT NULL,
   credit_tokens integer NOT NULL,
   updated_at integer NOT NULL,
   created_at integer NOT NULL
);

CREATE INDEX customers_auth_code_idx ON customers(auth_code);

CREATE TABLE wallet_passes (
    id UUID PRIMARY KEY,
    customer_email character varying REFERENCES customers(email),
    file_reference character varying NOT NULL,
    payload JSONB NOT NULL,
    last_redeemed_at integer,
    created_at integer NOT NULL
);

CREATE INDEX wallet_passes_customer_email_idx ON wallet_passes(customer_email);