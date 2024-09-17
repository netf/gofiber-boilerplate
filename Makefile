.PHONY: build run migrate-up migrate-down test clean swag docker-build docker-run docker-up docker-down

build:
	go build -o bin/golang-fiber-boilerplate cmd/main.go

run: build
	./bin/golang-fiber-boilerplate

migrate-up:
	migrate -database $(DATABASE_URL) -path migrations up

migrate-down:
	migrate -database $(DATABASE_URL) -path migrations down

swag:
	swag init -g cmd/main.go --output docs --parseDependency --parseInternal

test:
	go test ./...

clean:
	rm -rf bin/

docker-build:
	docker build -t golang-fiber-boilerplate -f docker/Dockerfile .

docker-run: docker-build
	docker run -p 8080:8080 golang-fiber-boilerplate

docker-up:
	docker-compose -f docker/docker-compose.yml up --build

docker-down:
	docker-compose -f docker/docker-compose.yml down
