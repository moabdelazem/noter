include .env

# Build the application and save it to bin/api
build:
	echo "Building the application..."
	@CGO_ENABLED=0 go build -o bin/api cmd/api/main.go
	echo "Build completed successfully Check bin/api"

# Run the application
run: build
	echo "Running the application..."
	@bin/api

# Watch the application
watch:
	echo "Watching the application..."
	@air

# Clean the application
clean:
	@rm -f bin/api
	echo "Clean completed successfully"

# Run the tests
test:
	echo "Running the tests..."
	@go test -v ./...
	echo "Tests completed successfully"

# Create a new migration file
migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# Run migrations up
migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" up

# Run migrations down
migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)" down

# Start the database
db-start:
	docker-compose up -d noter_db

# Stop the database
db-stop:
	docker-compose down

# Docker Build
docker-build:
	docker build -t moabdelazem/noter .

# Docker Run
docker-run:
	docker run -p 8080:8080 moabdelazem/noter

# Docker Push
docker-push:
	docker push moabdelazem/noter

# Docker Pull
docker-pull:
	docker pull moabdelazem/noter

# Environment variables for database
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_NAME ?= noter
DB_SSLMODE ?= disable
