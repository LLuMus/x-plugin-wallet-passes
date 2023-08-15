.PHONY: test run generate

default: run

generate:
	oapi-codegen -package components --config=./codegen.config.yml ./public/openapi.yml > ./internal/components/components.gen.go

test:
	go test ./... -coverprofile=coverage.out

run:
	PORT=80 go run ./cmd/plugin/main.go
