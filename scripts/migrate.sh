#!/bin/bash

# Migration script

set -e

DB_DRIVER=${DATABASE_DRIVER:-sqlite3}
DB_DSN=${DATABASE_DSN:-"file:roombooker.db?cache=shared&_fk=1"}

if [ "$DB_DRIVER" = "postgres" ]; then
    migrate -path ./migrations -database "postgres://$DB_DSN" "$@"
else
    migrate -path ./migrations -database "sqlite3://$DB_DSN" "$@"
fi
