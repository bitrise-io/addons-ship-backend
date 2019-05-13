#!/bin/bash
set -ex

export DB_CONN_STRING="user=$DB_USER dbname=$DB_NAME password=$DB_PWD host=$DB_HOST sslmode=disable"
go build -i -o db/goose db/*.go
