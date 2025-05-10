# Use official Golang image
FROM golang:1.23

# Set working directory
WORKDIR /coupon-system

# Download dependencies early
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o main ./cmd/app

# Set the entry point
CMD ["./main"]
