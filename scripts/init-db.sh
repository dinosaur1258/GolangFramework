#!/bin/sh
set -e

echo "Waiting for PostgreSQL to be ready..."
sleep 5

echo "Running migrations..."
migrate -path /app/db/migrations -database "postgresql://postgres:postgres@postgres:5432/golang_framework?sslmode=disable" up

echo "Migrations completed!"