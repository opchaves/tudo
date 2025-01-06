# Tudo

A Note taking web application written in Go and TypeScript

# Requirements

- Go 1.23+
- Docker
- PostgreSQL client

# Development

- Clone the repo

```sh
git clone git@github.com:opchaves/tudo.git
```

- Create the .env file

```sh
cp .env.example .env
```

- Start postgres

```sh
docker compose up -d --build
```

- Run the migrations

```sh
make migrate-up
```

- Run the server

```sh
go run main.go
```
