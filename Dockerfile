# Use the official Golang image as the base image
FROM golang:1.24-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code from the current directory to the working directory inside the container
COPY . .

# Build the Go app
RUN go build -o main cmd/main.go

# Run tests
# RUN go test ./...

# Use a minimal Alpine Linux image for the final image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy necessary files (e.g., schema.sql, data_csv)
COPY database/schema.sql ./database/schema.sql
COPY data_csv/SWIFT_CODES.csv ./data_csv/SWIFT_CODES.csv

# Expose the port the app listens on
EXPOSE 8080

# Set environment variables
# ENV DB_USER=postgres
# ENV DB_PASSWORD=root
# ENV DB_HOST=localhost
# ENV DB_PORT=5433
# ENV DB_NAME=demodb

# Command to run the executable
CMD ["./main"]
