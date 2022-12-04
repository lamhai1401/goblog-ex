SHELL := /bin/bash

swagger:
	swagger generate spec -o ./swagger.json
	swagger serve ./swagger.json

mongo:
	docker rm -f mongodb || true
	docker run -d --name mongodb -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=password -p 27017:27017 mongo:4.4.3

redis:
	docker run -d --name redis -p 6379:6379 redis:6.0

run_server:
	JWT_SECRET=eUbP9shywUygMx7u \
	X_API_KEY=eUbP9shywUygMx7u \
	MONGO_URI="mongodb://admin:password@localhost:27017/test?authSource=admin" \
	MONGO_DATABASE=test \
	go run .

bench:
	ab -n 2000 -c 100 -g with-cache.data http://localhost:8080/recipes

test:
	go test -v -cover ./...

test-circle:
	go test -v -race -coverprofile=c.out -cover $(go list ./... | circleci tests split --split-by=timings)
	go tool cover -html=c.out -o coverage.html

test-local:
	go test -v -race -coverprofile=c.out -cover ./...
	go tool cover -html=c.out -o coverage.html

graph-init:
	go get github.com/99designs/gqlgen
	go run github.com/99designs/gqlgen init

graph-modify:
	go get github.com/99designs/gqlgen
	go run github.com/99designs/gqlgen generate .

.PHONY: test