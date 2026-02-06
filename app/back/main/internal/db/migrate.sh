#!/bin/sh

# Load env vars from env.env
set -a
. ../../env.env
set +a

# Run all migrations
migrate \
  -path ./migrations \
  -database "$DATABASE_URL" \
  up
