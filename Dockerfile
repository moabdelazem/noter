# Selection of the base image
FROM golang:1.24-alpine3.21 AS builder

# Set the working directory
WORKDIR /code

# Copy the go mod and sum files
COPY go.mod go.sum ./

# Download the dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o bin/api cmd/api/main.go

# Stage 2: Run the application in non-shell mode
# Using the distroless image
FROM gcr.io/distroless/static-debian12:nonroot

# Copy the binary from the builder stage
COPY --from=builder /code/bin/api /bin/api

# Using the user nonroot
USER nonroot:nonroot

# Expose the port
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["/bin/api"]





