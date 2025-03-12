# Crypto Monitor

A real-time cryptocurrency monitoring application built with Go and Angular.

## Features

- Real-time cryptocurrency price monitoring
- User portfolio management
- Price alerts
- WebSocket-based live updates
- Secure authentication

## Prerequisites

- Docker and Docker Compose
- Go 1.19+
- Node.js 16+
- PostgreSQL 13+
- Redis 6+

## Getting Started

1. Clone the repository:
bash
git clone https://github.com/yourusername/crypto-monitor.git
cd crypto-monitor

2. Copy and configure environment variables:
cp .env.example .env
# Edit .env with your configurations

3. Start the application:
docker-compose up -d

4. Run databasemigrations:
docker-compose run --rm backend ./migrations up


## Development

1. Backend
cd backend
go mod download
go run cmd/server/main.go

2. Frontend
cd frontend
npm install
ng serve


## Testing

# Backend tests
cd backend
go test ./...

# Frontend tests
cd frontend
ng test


## Deployment
./scripts/deploy.sh


## Documentation
API documentation is available at /docs/api.yaml


## Monitoring
The application includes Prometheus metrics and Grafana dashboards for monitoring.