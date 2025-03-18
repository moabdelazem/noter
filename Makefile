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

