# Stage 1: Build the binary
FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /coupon-system

# Copy go.mod and go.sum first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application binary
RUN go build -o main ./cmd/app

# Stage 2: Run the binary in a minimal image
FROM alpine:latest

# Set working directory
WORKDIR /coupon-system

# Copy the built binary from the builder stage
COPY --from=builder /coupon-system/main .

# Set entry point to run the binary
CMD ["./main"]
