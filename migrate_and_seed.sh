#!/bin/sh
set -e

echo "Running migrations..."
go run migrator

echo "Running seed..."
psql "$DATABASE_URL" -f /migrations/seed.sql

echo "Done."
