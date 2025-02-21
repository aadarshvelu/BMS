# Books Management System (BMS)

- **GO Version: 1.23.4**
- **Go Package: github.com/aadarshvelu/bms**
- **API Documentation: [github.com/aadarshvelu/bms](http://13.203.8.233/docs/index.html#/books)**

## Tech Stack
- Go
- Kafka
- Redis
- PostgreSQL
- Gin (Web Framework)
- Gorm (ORM)

# Development Setup

## Using Air for Live Reload

- Air provides live reloading during development. The configuration varies by operating system.
- Initialize .env with your local credentials

#### 1. Install Air:

```go install github.com/cosmtrek/air@latest```

#### 2. Configure `.air.toml` based on your OS:

.air.toml already initialized for windows, for other users, change the bin in .air.toml:7

##### For Windows:

```bin = "tmp\\main.exe"```

##### For Mac/Linux:

```bin = "tmp/main"```

#### 3. Start Air Server:

```air```

## Using Docker

- Kakfa, Redis, Postgres are already declared in dockerfile
- Initialize .env.docker with your docker credentials

#### 1. Install Docker & Docker-compose:

##### For Windows:

```docker compose --env-file .env.docker up```

##### For Mac/Linux:

```docker-compose --env-file .env.docker up -d --build```

# Project Structure:


### Directory Structure Explanation

- **app/**: Core application code
  - `cache/`: Redis caching implementation
  - `events/`: Event handling (Kafka)
  - `handler/`: HTTP route handlers
  - `helpers/`: Utility functions and validators
  - `models/`: Data models
  - `repo/`: Database repositories
  - `services/`: Business logic implementation

- **cmd/**: Application entry points
  - `main.go`: Main application bootstrap

- **docs/**: API documentation
  - Swagger/OpenAPI specifications

- **config/**: Configuration management
  - `config.go`: Environment and app configuration
  - `database.go`: Database configuration

- **pkg/log**: Shared packages
  - `logger.go`: Logging utilities

# Logs & Events:

- logs are captured and stored in **api_logs.txt** in root
- Kafka Events preview:

![alt text](image.png)