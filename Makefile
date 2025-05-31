BINARY_NAME=edu-helper
SRC_EDUHELPER=./cmd/eduhelper
SRC_MIGRATOR=./cmd/migrator
MIGRATIONS_PATH=./migrations
CONFIG_PATH=./config/

CONFIG_FILE?=local.yaml

user?=root
password?=
host?=localhost
port?=3306
db_name?=test
table?=schema_migrations

DB_URL=mysql://$(user):$(password)@tcp($(host):$(port))/$(db_name)?multiStatements=true
MIGRATE=go run github.com/golang-migrate/migrate/v4/cmd/migrate@latest

.PHONY: all build run test lint tidy clean migrate-up migrate-down migrate

all: build

build:
	go build -o bin/$(BINARY_NAME) $(SRC_EDUHELPER)

run:
	go run $(SRC_EDUHELPER) -config='$(CONFIG_PATH)/$(CONFIG_FILE)'

test:
	go test ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

clean:
	rm -f bin/$(BINARY_NAME)

migrate-up:
	go run $(SRC_MIGRATOR) \
		--migrations-path=$(MIGRATIONS_PATH) \
		--migrations-table=$(table) \
		--db-user=$(user) \
		--db-password=$(password) \
		--db-host=$(host) \
		--db-port=$(port) \
		--db-name=$(db_name)

migrate-down:
	go run $(SRC_MIGRATOR) \
		--migrations-path=$(MIGRATIONS_PATH) \
		--migrations-table=$(table) \
		--db-user=$(user) \
		--db-password=$(password) \
		--db-host=$(host) \
		--db-port=$(port) \
		--db-name=$(db_name)
		--down
