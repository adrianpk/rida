APP_NAME = rida

RIDA_API_KEY ?= demo-api-key
RIDA_OTTAWA_CLIENTS ?= 1
RIDA_MONTREAL_CLIENTS ?= 2
RIDA_HTTP_PORT ?= :8080

.PHONY: all build run run-race test lint format check install-hooks run-docker stop-docker test-docker

all: build

build:
	go build -o bin/rida .

run: build
	./bin/$(APP_NAME) \
		-api-key=$(RIDA_API_KEY) \
		-ottawa-clients=$(RIDA_OTTAWA_CLIENTS) \
		-montreal-clients=$(RIDA_MONTREAL_CLIENTS) \
		-http-port=$(RIDA_HTTP_PORT)

run-race:
	go run -race main.go \
		-api-key=$(RIDA_API_KEY) \
		-ottawa-clients=$(RIDA_OTTAWA_CLIENTS) \
		-montreal-clients=$(RIDA_MONTREAL_CLIENTS) \
		-http-port=$(RIDA_HTTP_PORT)

test:
	go test ./...

test-race:
	go test -race ./...

lint:
	golangci-lint run ./...

format:
	go fmt ./...

check: format lint

install-hooks:
	ln -sf scripts/hooks/pre-commit .git/hooks/pre-commit
	chmod +x scripts/hooks/pre-commit

run-docker:
	docker-compose -f deployment/docker-compose.yml up --build

stop-docker:
	docker-compose -f deployment/docker-compose.yml down

test-docker:
	docker-compose -f deployment/docker-compose.yml up -d --build
	sleep 3
	curl --fail http://localhost:8080/healthz || (echo "\n[ERROR] The app did not respond in Docker" && exit 1)
	docker-compose -f deployment/docker-compose.yml down
