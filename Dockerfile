# Build stage
FROM golang:1.24.0-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod ./

# Download dependencies (this will create go.sum)
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Build the binary
RUN go build -o insider-case ./cmd/main.go

# Stage 2: Minimal runtime image
FROM golang:1.24.0-alpine

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/insider-case .

# Copy config
COPY --from=builder /app/internal/config/config.yaml /app/internal/config/config.yaml

CMD ["./insider-case"]