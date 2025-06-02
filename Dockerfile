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

ENV STEAM_API_KEY="$STEAM_API_KEY"
ENV DATABASE_URL="$DATABASE_URL"
ENV PLATFORM="dev"
ENV PORT="8080"
ENV JWTSECRET="$JWT_SECRET"

CMD ["/app/start.sh"]
