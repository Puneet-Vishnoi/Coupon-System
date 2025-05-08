# Coupon System

The Coupon System is a Go-based application designed to manage coupons, validate them, and provide functionality to check applicable coupons for users. It integrates with PostgreSQL for data storage and Redis for caching. This project also includes unit tests, integration tests, and Docker setup for easy deployment.

## Features

- **Coupon Management**: Create, validate, and fetch applicable coupons.
- **Redis Caching**: Utilizes Redis for storing frequently accessed coupon data.
- **PostgreSQL**: Stores all coupon-related and usage data.
- **Dockerized**: The system is containerized with Docker Compose for easy setup.
- **Testing**: Unit and integration tests ensure code stability.

## Technologies Used

- **Go**: The primary language for building the application.
- **PostgreSQL**: Database for storing coupon data and usage.
- **Redis**: Caching layer for performance optimization.
- **Docker**: Used for containerizing the application and related services.
- **Go Modules**: For managing dependencies.

## Directory Structure

```
coupon-system/
├── cache/redis
|           └── providers/providers.go
│           └── database.go
├── db/postgre
|           └── providers/providers.go
│           └── database.go
├── cmd/app
│   └── main.go
├── handlers/
│   └── handlers.go
├── service/
│   └── coupon_service.go
├── repository/
│   └── coupon_repo.go
├── models/
│   └── coupon.go
│   └── request.go
|   └── response.go
├── routes/
│   └── router.go
├── tests
|     └── mockdb/connection.go
│     └── unittest/unit_test.go
│     └── integration/integretion_test.go
├── Dockerfile
├── go.mod
├── go.sum
└── README.md                     # This file.
```

## Getting Started

To get the application up and running locally, follow these steps.

### Prerequisites

- **Docker** and **Docker Compose** installed on your system.
- **Go** installed for local development and testing.

### 1. Clone the Repository

Clone the repository to your local machine:

```bash
git clone https://github.com/Puneet-Vishnoi/Coupon-System.git
cd coupon-system
```

### 2. Install Dependencies

Install the Go dependencies:

```bash
go mod tidy
```

### 3. Set Up Environment Variables

Create a `.env` file in the root directory and add the following environment variables:

```
MONGO_URI=mongodb://localhost:27017/coupon-system
PORT=8080
REDIS_ADDR=redis:6379
POSTGRES_USER="user name"
POSTGRES_PASSWORD="Your password"
POSTGRES_DB=coupon-system
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_DSN=postgres://postgres:User@postgres:5432/coupon-system?sslmode=disable
```

### 4. Run the Application with Docker

Run the application and its dependencies (PostgreSQL, Redis) using Docker Compose:

```bash
docker-compose up -d
```

### 5. Access the Application

Once the containers are running, the application will be accessible at http://localhost:8080.

To check the logs for each service, run:

```bash
docker-compose logs -f
```

## Available Routes

### Coupon Management Routes

#### Create Coupon

**Route**: `POST /api/coupons`

**Description**: Create a new coupon.

**Payload**:
```json
{
  "code": "NEWYEAR2025",
  "discount": 20,
  "validity": "2025-12-31"
}
```

#### Get Applicable Coupons

**Route**: `POST /api/coupons/applicable`

**Description**: Retrieve a list of applicable coupons for a user.

**Payload**:
```json
{
  "user_id": "1234"
}
```

#### Validate Coupon

**Route**: `POST /api/coupons/validate`

**Description**: Validate if a coupon is valid.

**Payload**:
```json
{
  "coupon_code": "NEWYEAR2025"
}
```

### Redis Caching Routes

#### Get Cache

**Route**: `GET /api/cache`

**Description**: Retrieve cache details from Redis.

### PostgreSQL Routes

#### Database Connection

The application will connect to PostgreSQL for all data storage operations, including managing coupons and their usage.

## Running Tests

To run unit tests and integration tests:

```bash
go test ./...
```

### Dockerized Tests

If you want to run tests within the Docker container:

1. Start the Docker containers:
```bash
docker-compose up -d
```

2. Access the coupon-app container:
```bash
docker exec -it coupon-app sh
```

3. Run the tests:
```bash
go test ./... -v
```

## Docker Compose Setup

This project uses Docker Compose to manage services, including:

- **coupon-app**: The main application.
- **postgres**: PostgreSQL database for storing coupon data.
- **redis**: Redis instance for caching.

To bring up the services, use:
```bash
docker-compose up -d
```

To shut them down:
```bash
docker-compose down
```

## Notes

- **Unit Tests**: Located in the tests folder.
- **Integration Tests**: Located in tests/integration_test.go.
- **Mocking**: Mock interfaces using GoMock or Testify for unit testing.
- **Redis**: Caching logic is handled by Redis in the cache folder.

## License

This project is licensed under the MIT License - see the LICENSE file for details.