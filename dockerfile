# Use official Golang image
FROM golang:1.23

# Set working directory
WORKDIR /coupon-system

# Copy go mod files first to leverage layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy entire project
COPY . .

# Explicitly copy coupon.sql to the test path INSIDE the container
COPY db/postgres/coupon.sql tests/integration/db/postgres/coupon.sql
COPY db/postgres/coupon.sql tests/unittest/db/postgres/coupon.sql


# Build the binary
RUN go build -o main ./cmd/app

# Set the entry point
CMD ["./main"]
