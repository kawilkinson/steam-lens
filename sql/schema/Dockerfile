FROM golang:1.24.3

RUN apt-get update && \
    apt-get install -y --no-install-recommends postgresql-client && \
    rm -rf /var/lib/apt/lists/*

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

WORKDIR /migrations
COPY ./ /migrations

ENTRYPOINT ["/bin/sh", "-c"]
CMD ["goose -dir ./ postgres "$DATABASE_URL" up"]
