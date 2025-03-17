include .env
GOLINT := $(GOPATH)/bin/golint

migrate-down:
	migrate -database '$(POSTGRES_CONN_LOCAL)?sslmode=disable' -path migrations down
.PHONY: migrate-down

migrate-up:
	migrate -database '$(POSTGRES_CONN_LOCAL)?sslmode=disable' -path migrations up
.PHONY: migrate-up

up:
	docker compose up
.PHONY: up

fmt:
	go fmt ./...
.PHONY: fmt

install-lint:
	go install golang.org/x/lint/golint@latest
.PHONY: install-lint

lint:
	$(GOLINT) ./...
.PHONY: lint

test:
	go test ./...
.PHONY: test

integration:
	go test ./... -tags=integration
.PHONY: integration

cover:
	go tool cover -html=cover.out -o cover.html
.PHONY: cover

