FROM golang:1.24.3 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends postgresql-client && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/server /app/server
COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY ./sql/schema /app/migrations
COPY ./scripts/start.sh /app/start.sh
RUN chmod +x /app/start.sh

ENV PORT=8080
ENV PLATFORM="dev"
ENV STEAM_API_KEY="test"
ENV DATABASE_URL="postgres://steam_lens:password@db:5432/steam_lens_db?sslmode=disable"
ENV JWTSECRET="test"

CMD ["/app/start.sh"]
