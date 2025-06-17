# Makefile

# Variables
APP_NAME=coupon-app
DOCKER_COMPOSE=docker-compose
PORT=8080

# Build the Go binary and Docker image
build:
	docker build -t $(APP_NAME) .

# Start all services using docker-compose
up:
	$(DOCKER_COMPOSE) up -d --build

# Stop all services
down:
	$(DOCKER_COMPOSE) down

# View logs from the app
logs:
	$(DOCKER_COMPOSE) logs -f app

# Run tests inside the container (you can customize this)
unit-test:
	go test ./tests/unittest/... -v

integration-test:
	go test ./tests/integration/... -v

# Remove volumes and clean everything
clean:
	$(DOCKER_COMPOSE) down -v --remove-orphans

# Rebuild everything from scratch
rebuild: clean build up

# Check health of containers
ps:
	$(DOCKER_COMPOSE) ps

# Run go fmt on source files
fmt:
	go fmt ./...
