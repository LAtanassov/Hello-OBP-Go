#!/bin/sh -e

export DIRECT_URL="https://domain"
export USERNAME="username"
export PASSWORD="password"
export CONSUMER_KEY="consumer_key"


go test -run .IT. -tags=integration  -cover -race ./...
go vet ./...