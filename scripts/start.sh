#!/bin/sh
set -e

until pg_isready -h "$DB_HOST" -U "$DB_USER"; do
  >&2 echo "Postgres is unavailable - waiting..."
  sleep 2
done

echo "Running goose migrations..."
goose -dir /app/migrations postgres "$DATABASE_URL" up

echo "Starting server..."
exec /app/server
