.PHONY: build run migrate-up migrate-down test clean swag docker-build docker-run docker-up docker-down setup-pre-commit generate-auth-keys

build: swag
	go build -o ./bin/app ./cmd/main.go

run: build
	./bin/app start

migrate-up:
	go run cmd/main.go migrate up

migrate-down:
	go run cmd/main.go migrate down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

swag:
	which swag || go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g ./cmd/main.go -o ./docs --parseDependency --parseInternal

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

setup-pre-commit:
	pip install pre-commit
	pre-commit install

generate-auth-keys:
	@echo "Generating ECDSA private key..."
	@openssl ecparam -name prime256v1 -genkey -noout -out auth_private_key.pem
	@echo "Private key generated and saved as auth_private_key.pem"
	@echo "Generating salt..."
	@openssl rand -base64 32 > auth_salt.txt
	@echo "Salt generated and saved as auth_salt.txt"
	@echo "Please store these securely and never commit them to version control."
	@echo "You can use these values in your .env file or as environment variables:"
	@echo "AUTH_PRIVATE_KEY=\$$(cat auth_private_key.pem | base64 -w 0)"
	@echo "AUTH_SALT=\$$(cat auth_salt.txt)"
