#!/bin/sh
set -e

echo ">>> Running go mod tidy..."
go mod tidy

echo ">>> Starting Air (hot reload)..."
exec air -c .air.toml
