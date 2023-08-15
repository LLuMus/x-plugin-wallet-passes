FROM golang:1.21-bullseye

# Set destination for COPY
WORKDIR /app

ENV CGO_ENABLED=0

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY assets ./assets
COPY cmd ./cmd
COPY internal ./internal
COPY assets ./assets
COPY public/build ./public
COPY tmp ./tmp

RUN apt-get update
RUN apt-get upgrade -y
RUN DEBIAN_FRONTEND=noninteractive apt-get install -y build-essential libssl-dev

# Build
RUN GOOS=linux go build ./cmd/plugin/main.go

EXPOSE 80

# Run
CMD ["/app/main"]