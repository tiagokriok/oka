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

new_migration:
	migrate create -ext sql -dir migrations -seq $(name)

migrate_up:
	migrate -path migrations -database "mysql://${DB_URL}" -verbose up

migrate_down:
	migrate -path migrations -database "mysql://${DB_URL}" -verbose down

migrate_fix:
	migrate -path migrations -database "mysql://${DB_URL}" force $(version)
