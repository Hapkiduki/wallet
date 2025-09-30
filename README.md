# Go Wallet API 💳

A robust, production-ready Wallet REST API built in Go, demonstrating senior-level engineering practices including Clean Architecture, SOLID principles, and a complete containerized development environment.

## ✨ Features

- **User Management**: Create new users.
- **Wallet System**: Each user is automatically assigned a new, empty wallet upon creation.
- **Wallet Recharge**: Deposit funds into a wallet.
- **Funds Transfer**: Transfer funds between wallets with transactional integrity (i.e., funds are only transferred if the sender has a sufficient balance).

## 🏛️ Architecture Overview

This project is built using **Clean Architecture** principles to ensure a separation of concerns, maintainability, and testability. The code is organized into four distinct layers:

- **Domain**: Contains the core business models and interfaces (contracts). It has no external dependencies.
- **Usecase**: Contains the pure business logic that orchestrates the domain models to fulfill application-specific tasks.
- **Infrastructure**: Provides concrete implementations for the domain interfaces, such as the PostgreSQL repository, Redis cache, and transaction manager.
- **Handler**: The outermost layer, responsible for handling HTTP requests (using Fiber), parsing data, and calling the use cases.

The entire development environment is containerized with **Docker** and **Docker Compose**, featuring hot-reloading for a seamless development experience.

## 📁 Project Structure

```
.
├── cmd/api/            # Main application entry point
├── db/migration/       # SQL database migrations
├── docs/               # Auto-generated Swagger/OpenAPI docs
├── internal/
│   ├── config/         # Viper configuration loading
│   ├── domain/         # Core models and repository interfaces
│   ├── handler/        # Fiber HTTP handlers and DTOs
│   ├── infrastructure/ # GORM, Redis implementations
│   └── usecase/        # Business logic layer
├── .air.toml           # Air configuration for hot-reloading
├── .env                # Local environment variables (gitignored)
├── Dockerfile          # Production multi-stage Dockerfile
├── Dockerfile.dev      # Development Dockerfile with Air
├── Makefile            # Convenience commands (make up, make migrateup)
├── go.mod              # Go module definition
└── docker-compose.yml  # Docker services for local dev
```

## 🛠️ Tech Stack

- **Language**: [Go](https://golang.org/) 1.25.1+
- **Web Framework**: [Fiber](https://gofiber.io/)
- **ORM**: [GORM](https://gorm.io/)
- **Database**: [PostgreSQL](https://www.postgresql.org/)
- **Caching**: [Redis](https://redis.io/)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Containerization**: [Docker](https://www.docker.com/) & Docker Compose
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **API Documentation**: [Swagger (swag)](https://github.com/swaggo/swag)
- **Hot-Reloading**: [Air](https://github.com/cosmtrek/air)
- **Logging**: `log/slog` (Structured JSON logging)
- **Error Tracking**: [Sentry](https://sentry.io/) (Real-time error monitoring)
- **Graceful Shutdown**: Signal handling with connection cleanup

## 🚀 Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (v1.25.+)
- [Docker](https://docs.docker.com/get-docker/) & Docker Compose
- [Homebrew](https://brew.sh/) (on macOS, to install migrate)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate)
  ```bash
  brew install golang-migrate
  ```

### Installation & Running

1.  **Clone the repository:**
    ```bash
    git clone <your-repo-url>
    cd <your-repo-name>
    ```

2.  **Create your environment file:**
    Copy the example file and fill in your details. For local development, the defaults should work.
    ```bash
    cp .env.example .env
    ```

3.  **Start the environment:**
    This command will start the PostgreSQL and Redis containers in the background.
    ```bash
    make up
    ```

4.  **Run database migrations:**
    Wait a moment for the database to be ready, then apply the schema.
    ```bash
    make migrateup
    ```
    Your application will now be running with hot-reloading at `http://localhost:8080`.

## 📖 API Documentation

Once the application is running, the interactive Swagger/OpenAPI documentation is available at:

[http://localhost:8080/swagger](http://localhost:8080/swagger)

## ⚙️ Configuration

The application uses **Viper** for configuration management, automatically loading environment variables. All configuration options are documented in the `.env.example` file.

| Variable        | Description                               | Default                       | Required |
| --------------- | ----------------------------------------- | ----------------------------- | -------- |
| `SERVER_PORT`   | Port for the API server to listen on     | `8080`                        | No       |
| `DB_SOURCE`     | PostgreSQL connection string              | (See `.env.example`)          | Yes      |
| `REDIS_ADDR`    | Redis connection address                  | `localhost:6379`              | Yes      |
| `SENTRY_DSN`    | DSN for Sentry error reporting           | `""`                          | No       |
| `GO_ENV`        | Environment (development/production)      | `development`                 | No       |

### Configuration Loading
The application uses a structured configuration approach:
- Environment variables are automatically loaded via Viper
- Configuration is validated at startup
- Missing required variables cause graceful startup failure

## 🔍 Observability & Monitoring

### Structured Logging
The application uses Go's native `log/slog` package for structured JSON logging:
- **Structured format**: All logs are in JSON format for easy parsing
- **Contextual information**: Each log entry includes relevant context (user_id, request_id, etc.)
- **Log levels**: DEBUG, INFO, WARN, ERROR levels supported
- **Performance**: High-performance logging with minimal overhead

Example log output:
```json
{"time":"2024-01-01T10:00:00Z","level":"INFO","msg":"Starting server","port":"8080"}
{"time":"2024-01-01T10:00:01Z","level":"ERROR","msg":"Database connection failed","error":"connection refused"}
```

### Error Tracking with Sentry
- **Real-time error monitoring**: Automatic error capture and reporting
- **Error grouping**: Similar errors are grouped for easier analysis  
- **Performance monitoring**: Track API response times and database queries
- **Release tracking**: Monitor errors across different deployments
- **Alerting**: Get notified immediately when errors occur

### Graceful Shutdown
The application implements graceful shutdown to ensure data integrity:
- **Signal handling**: Listens for SIGINT and SIGTERM signals
- **Connection cleanup**: Properly closes database and Redis connections
- **Request completion**: Allows in-flight requests to complete (15s timeout)
- **Resource cleanup**: Ensures all resources are properly released

## ⚡ Caching Strategy

### Redis-based Caching
The application implements a **Decorator Pattern** for caching:
- **User data caching**: Frequently accessed user data is cached in Redis
- **Cache-aside pattern**: Data is loaded from cache first, then from database if not found
- **Automatic invalidation**: Cache is updated when data changes
- **Performance boost**: Reduces database load and improves response times

### Cache Implementation
```go
// Cached repository wraps the original repository
type CachedUserRepository struct {
    cache    CacheRepository
    original UserRepository
}

// Transparent caching - no changes needed in business logic
userRepo := cache.NewCachedUserRepository(cacheRepo, postgresUserRepo)
```

## 🧪 Running Tests

To run all unit tests:
```bash
go test -v ./...
```

To run all benchmark tests:
```bash
go test -bench=. ./...
```

## 🛠️ Development Commands

The project includes a comprehensive Makefile for common development tasks:

```bash
# Start all services (PostgreSQL, Redis)
make up

# Stop all services  
make down

# Apply database migrations
make migrateup

# Rollback database migrations
make migratedown

# Generate/update API documentation
make docs
```

## 🏗️ Architecture Patterns

### Clean Architecture
The project follows Clean Architecture principles with clear separation of concerns:

- **Domain Layer**: Core business models and interfaces (no external dependencies)
- **Use Case Layer**: Business logic that orchestrates domain models
- **Infrastructure Layer**: External concerns (database, cache, HTTP)
- **Handler Layer**: HTTP request/response handling and data transformation

### Design Patterns Used
- **Repository Pattern**: Abstract data access logic
- **Decorator Pattern**: Add caching behavior transparently  
- **Dependency Injection**: Loose coupling between layers
- **Factory Pattern**: Create configured instances
- **Transaction Script**: Handle complex business operations atomically

### SOLID Principles
- **Single Responsibility**: Each component has one reason to change
- **Open/Closed**: Open for extension, closed for modification
- **Liskov Substitution**: Interfaces can be substituted seamlessly
- **Interface Segregation**: Small, focused interfaces
- **Dependency Inversion**: Depend on abstractions, not concretions

## 🚀 Production Considerations

### Performance
- **Connection pooling**: Database connections are pooled for efficiency
- **Caching**: Redis caching reduces database load
- **Structured logging**: Minimal performance overhead
- **Graceful shutdown**: Prevents data corruption during deployments

### Security
- **Input validation**: All API inputs are validated
- **SQL injection prevention**: GORM provides built-in protection
- **Error handling**: Sensitive information is not exposed in error messages
- **Environment-based configuration**: Secrets are managed via environment variables

### Monitoring
- **Health checks**: Built-in endpoints for monitoring service health
- **Error tracking**: Sentry integration for real-time error monitoring
- **Structured logs**: Easy to parse and analyze
- **Metrics**: Ready for integration with Prometheus/Grafana