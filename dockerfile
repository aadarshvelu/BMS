# Golang image
FROM golang:1.21 AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the application
RUN CGO_ENABLED=0 go build -o main ./cmd/main.go

# Use a minimal base image for production
FROM alpine:latest

WORKDIR /root/

# Install CA certificates for HTTPS connections
RUN apk --no-cache add ca-certificates

# Copy the compiled binary
COPY --from=builder /app/main .
COPY --from=builder /app/.env .env

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./main"]