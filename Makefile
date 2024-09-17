.PHONY: build run migrate-up migrate-down test clean swag

build:
	go build -o bin/golang-fiber-boilerplate cmd/main.go

run: build
	./bin/golang-fiber-boilerplate

migrate-up:
	migrate -database $(DATABASE_URL) -path migrations up

migrate-down:
	migrate -database $(DATABASE_URL) -path migrations down

swag:
	swag init -g cmd/main.go --output docs

test:
	go test ./...

clean:
	rm -rf bin/
