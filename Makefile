.PHONY: up down deps run build start-binary

up:
	docker compose up --build

down:
	docker compose down

deps:
	go mod download

run:
	go run cmd/main.go

build:
	go build -o bin/service cmd/main.go

start-binary:
	bin/service