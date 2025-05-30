#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd sqlc/schema
goose postgres $DATABASE_URL up
