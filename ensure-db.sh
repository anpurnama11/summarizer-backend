#!/bin/sh
set -e

# Ensure the database directory exists
mkdir -p db

# Create the database file if it doesn't exist
if [ ! -f db/database.sqlite ]; then
    echo "Creating new SQLite database file..."
    touch db/database.sqlite
    sqlite3 db/database.sqlite "PRAGMA journal_mode=WAL;"
fi

# Run migrations
echo "Running migrations..."
migrate -path db/migrations -database "sqlite3://db/database.sqlite" up

echo "Migrations completed successfully"