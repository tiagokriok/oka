include .env
export

build:
	go build -o tmp/main ./cmd/server/...

dev: docker_up
	air

start: build
	@./tmp/main

install:
	@go mod tidy

docker_up:
	docker compose up -d

docker_down:
	docker compose down
