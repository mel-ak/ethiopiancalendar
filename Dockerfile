# Build stage
FROM golang:1.25.1-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o ethiopian-calendar ./api

# Final stage
FROM alpine:latest

# Install CA certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/ethiopian-calendar .
# Copy static files
COPY public/ ./public/

# Expose port 8080
EXPOSE 8080

# Command to run the application
CMD ["./ethiopian-calendar"]