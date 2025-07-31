# Insider Case - Messaging Service

A robust Go-based messaging service that handles asynchronous message processing with scheduling capabilities. The service stores messages in a PostgreSQL database and periodically processes pending messages by sending them to external webhook endpoints.

## 🏗️ Architecture Overview

This service implements a message queue system with the following key components:

- **HTTP API Server**: RESTful API built with Fiber v3 for managing the messaging system
- **Message Repository**: PostgreSQL-based storage for message persistence
- **Scheduler**: Background processor that periodically sends pending messages
- **Cache Layer**: Redis-based caching for performance optimization
- **Message Sender**: HTTP client for delivering messages to external webhook endpoints
- **Validator**: Request validation using go-playground/validator

## 📋 Features

- **Asynchronous Message Processing**: Messages are stored and processed asynchronously
- **Configurable Scheduling**: Adjustable intervals for message processing
- **Message Status Tracking**: Track messages through pending, success, and failed states
- **Webhook Integration**: Send messages to external webhook endpoints
- **Pagination Support**: Retrieve sent messages with pagination
- **Graceful Shutdown**: Proper cleanup of database connections, cache, and background processes
- **Health Checks**: Docker health checks for PostgreSQL and Redis
- **Request Validation**: Comprehensive input validation and error handling

## 🛠️ Technology Stack

### Core Dependencies

- **Go 1.24**: Latest Go version for modern language features
- **PostgreSQL**: Primary database for message persistence
- **Redis**: Caching layer for performance optimization
- **Docker**: Containerization and service orchestration

## 🚀 Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository**
   ```bash
   git clone https://github.com/berkinyildiran/insider-case
   cd insider
   ```

2. **Start all services**
   ```bash
   docker-compose up -d
   ```

3. Insert mock data by running the following commands:
   ```bash
   docker exec -i insider-case-berkin-yildiran-postgres-1 psql -U postgres -d postgres < ./scripts/data.sql
   ```

The application will be available at `http://localhost:8080`

## 📁 Project Structure

```
insider/
├── cmd/
│   └── main.go                 # Application entry point
├── docs/
│   └── swagger.yaml            # Swagger Documentation
├── internal/
│   ├── cache/                  # Cache layer (Redis)
│   │   ├── cache.go
│   │   ├── config.go
│   │   └── redis/
│   │       └── redis.go
│   ├── config/                 # Configuration management
│   │   ├── config.go
│   │   └── config.yaml
│   ├── database/               # Database layer
│   │   ├── config.go
│   │   └── database.go
│   ├── message/                # Message domain
│   │   ├── dto.go              # DTOs
│   │   ├── handler.go          # HTTP handlers
│   │   ├── model.go            # Message data model
│   │   ├── repository.go       # Database operations
│   │   ├── response.go         # Responses
│   │   └── status.go           # Message status constants
│   ├── scheduler/              # Background job scheduler
│   │   ├── config.go
│   │   └── scheduler.go
│   ├── sender/                 # Message sending logic
│   │   ├── config.go
│   │   ├── payload.go
│   │   └── sender.go
│   ├── server/                 # HTTP server setup
│   │   ├── config.go
│   │   └── router.go
│   ├── transporter/            # HTTP client abstraction
│   │   ├── http/
│   │   │   └── http.go
│   │   └── transporter.go
│   └── validator/              # Request validation
│       └── validator.go
├── scripts/
│   └── data.sql                # Mock data for PostgreSQL
├── docker-compose.yml          # Docker services configuration
├── Dockerfile                  # Container build instructions
├── go.mod                      # Go module definition
└── go.sum                      # Go module checksums
```

## 🔌 API Endpoints

### Message Management

| Method | Endpoint | Description | Query Parameters |
|--------|----------|-------------|------------------|
| `GET` | `/messaging/sent` | Retrieve sent messages with pagination | `limit`, `offset` |
| `POST` | `/messaging/start` | Start the message scheduler | - |
| `POST` | `/messaging/stop` | Stop the message scheduler | - |

#### Example Requests

**Get Sent Messages**
```bash
  curl "http://localhost:8080/messaging/sent?limit=10&offset=0"
```

**Start Scheduler**
```bash
  curl -X POST "http://localhost:8080/messaging/start"
```

**Stop Scheduler**
```bash
  curl -X POST "http://localhost:8080/messaging/stop"
```

## 💾 Database Schema

### Messages Table

| Column | Type | Description |
|--------|------|-------------|
| `id` | UUID | Primary key (auto-generated) |
| `content` | VARCHAR(100) | Message content |
| `recipient_phone_number` | VARCHAR(20) | Recipient's phone number |
| `sending_status` | SMALLINT | Message status (0=pending, 1=success, 2=failed) |
| `created_at` | TIMESTAMP | Creation timestamp |
| `updated_at` | TIMESTAMP | Last update timestamp |

**Indexes:**
- Composite index on `(sending_status, created_at)` for efficient pending message queries

## 🔄 Message Processing Flow

1. **Message Creation**: Messages are stored in the database with `pending` status
2. **Scheduler Processing**: Background scheduler fetches pending messages in batches
3. **Message Sending**: Each message is sent to the configured webhook endpoint
4. **Status Update**: Message status is updated to `success` or `failed` based on response
5. **Caching**: Successfully sent messages are cached to prevent duplicate processing

## 🛑 Graceful Shutdown Implementation

The application implements a comprehensive graceful shutdown mechanism to ensure data integrity and proper resource cleanup when the service is terminated.

### Signal Handling

The application listens for system signals to initiate graceful shutdown:

**Supported Signals:**
- `SIGINT` (Ctrl+C): Interactive interrupt signal
- `SIGTERM`: Termination request signal (commonly used by Docker and orchestrators)

### Shutdown Sequence

When a shutdown signal is received, the application follows a specific sequence to ensure clean termination:

#### 1. HTTP Server Shutdown

- Stops accepting new HTTP requests
- Completes processing of in-flight requests
- Closes all active connections gracefully

#### 2. Background Scheduler Termination

- Stops the background message processing scheduler
- Prevents new message batch processing from starting
- Allows current batch processing to complete

#### 3. Cache Connection Cleanup

- Properly closes Redis client connections
- Ensures cached data is persisted
- Prevents connection leaks

#### 4. Database Connection Cleanup

- Closes PostgreSQL connection pool
- Ensures all transactions are completed
- Prevents database connection leaks

### Benefits of This Implementation

- **Data Integrity**: No message processing is interrupted mid-transaction
- **Resource Management**: All connections and resources are properly released
- **Container Compatibility**: Works seamlessly with Docker and Kubernetes
- **Monitoring**: Detailed logging for each shutdown step
- **Error Handling**: Continues shutdown process even if individual steps fail
- **Timeout Protection**: Each component handles its own shutdown timeouts

This implementation ensures that the service can be safely terminated in any environment without data loss or resource leaks.

## 🏗️ Cache Layer Design Patterns

The cache layer implements several well-established design patterns to ensure maintainability, testability, and extensibility.

### 1. Interface Segregation Principle (ISP)

- **Minimal Dependencies**: Clients only depend on methods they actually use
- **Easy Testing**: Simple interface makes mocking straightforward
- **Clear Contract**: Interface clearly defines cache capabilities
- **Future-Proof**: Easy to extend without breaking existing code

### 2. Strategy Pattern

- **Runtime Flexibility**: Switch cache implementations without code changes
- **Environment Adaptation**: Different cache strategies for dev/staging/prod
- **Performance Optimization**: Choose optimal cache based on requirements
- **Vendor Independence**: Not locked into specific cache technology

### 3. Dependency Injection Pattern

- **Testability**: Easy to inject mock dependencies for testing
- **Loose Coupling**: Redis implementation doesn't create its own dependencies
- **Configuration Flexibility**: Runtime configuration injection
- **Lifecycle Management**: External control over dependency lifecycles

### 4. Context Pattern

- **Timeout Control**: Operations can be cancelled or timeout
- **Request Tracing**: Context can carry request-scoped values
- **Cancellation Propagation**: Upstream cancellations are respected
- **Resource Management**: Prevents hung operations

## 🐳 Docker Services

The `docker-compose.yaml` defines three services:

- **app**: The main Go application
- **postgres**: PostgreSQL database with health checks
- **redis**: Redis cache with persistence