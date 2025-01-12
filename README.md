# A/B Testing Service

A Go-based service for splitting web traffic for A/B testing purposes. The service includes a REST API, proxy modules for traffic splitting, and comprehensive monitoring capabilities.

## Features

- Dynamic proxy configuration for A/B testing
- REST API for managing proxies and viewing statistics
- Cookie-based user session persistence
- Traffic splitting based on configurable weights
- Prometheus metrics for monitoring
- Kafka integration for statistics processing
- Redis caching for proxy configurations
- PostgreSQL for persistent storage

## Requirements

- Docker and Docker Compose
- Go 1.23 or later (for local development)

## Quick Start

1. Clone the repository
2. Start the services using Docker Compose:
   ```bash
   docker-compose up -d
   ```

3. The service will be available at:
   - REST API: http://localhost:8080
   - Prometheus metrics: http://localhost:9090

## API Endpoints

- `GET /api/proxies` - List all proxies
- `POST /api/proxies` - Create a new proxy
- `GET /api/proxies/:id` - Get proxy details
- `DELETE /api/proxies/:id` - Delete a proxy
- `GET /api/proxies/:id/stats` - Get proxy statistics
- `PUT /api/proxies/:id/targets` - Update proxy targets

## Configuration

The service can be configured using the `config/config.yaml` file. Key configuration options include:

- Server port and host
- Database connection details
- Redis connection details
- Kafka configuration
- Prometheus settings

## Development

To run the service locally:

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Run the service:
   ```bash
   go run main.go
   ```

## Monitoring

The service exposes Prometheus metrics at `/metrics` endpoint, including:

- Request counts per target
- Request latencies
- Error rates

## License

MIT
