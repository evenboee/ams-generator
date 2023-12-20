# AMS

## Configure

Rename `example.env` to `.env`. Change values if necessary

## Start

1. Start docker with `docker compose --env-file=.env up -d --build`
2. Run migration with `go run cmd/migrate/main.go up`
3. Run generator with `go run cmd/main.go`
