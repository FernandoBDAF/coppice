# Standardized Service Folder Structure

## Overview

This document defines the standardized folder structure that all services must follow. The structure is designed to support clean architecture principles and maintain consistency across all services.

## Root Structure

```
service-name/
├── cmd/                    # Application entry points
│   └── main.go            # Main application entry point
├── internal/              # Private application code
│   ├── api/              # API layer (REST/gRPC handlers)
│   │   ├── handlers/     # HTTP/gRPC handlers
│   │   ├── middleware/   # HTTP/gRPC middleware
│   │   └── routes/       # Route definitions
│   ├── domain/           # Business logic and domain models
│   │   ├── models/       # Domain models
│   │   ├── services/     # Business logic services
│   │   └── interfaces/   # Interface definitions
│   ├── infrastructure/   # External service implementations
│   │   ├── database/     # Database implementations
│   │   ├── cache/        # Cache implementations
│   │   └── messaging/    # Message queue implementations
│   ├── config/           # Configuration management
│   │   ├── config.go     # Configuration structs
│   │   └── loader.go     # Configuration loaders
│   ├── pkg/              # Internal shared packages
│   │   ├── logger/       # Logging utilities
│   │   ├── metrics/      # Metrics collection
│   │   └── utils/        # Common utilities
│   └── server/           # Server setup and configuration
│       ├── http/         # HTTP server setup
│       └── grpc/         # gRPC server setup
├── pkg/                   # Public shared packages
│   └── client/           # Client libraries
├── api/                   # API definitions
│   ├── proto/            # Protocol buffer definitions
│   └── openapi/          # OpenAPI specifications
├── deployments/          # Deployment configurations
│   ├── docker/          # Docker-related files
│   └── k8s/             # Kubernetes manifests
├── docs/                 # Service documentation
│   ├── api/             # API documentation
│   └── architecture/     # Architecture documentation
├── scripts/             # Build and utility scripts
├── test/                # Test files
│   ├── integration/     # Integration tests
│   └── e2e/            # End-to-end tests
├── Dockerfile           # Service Dockerfile
├── docker-compose.yml   # Local development setup
├── go.mod              # Go module definition
├── go.sum              # Go module checksums
├── README.md           # Service documentation
└── .gitignore         # Git ignore rules
```

## Directory Purposes

### cmd/

- Contains the main application entry points
- Each executable should have its own subdirectory
- Main.go should be minimal, only handling initialization and startup

### internal/

- Contains all private application code
- Not importable by other services
- Organized by clean architecture layers

#### internal/api/

- API layer implementation
- Handlers, middleware, and route definitions
- Protocol-specific code (HTTP/gRPC)

#### internal/domain/

- Core business logic
- Domain models and interfaces
- Business rules and validations

#### internal/infrastructure/

- External service implementations
- Database, cache, and messaging implementations
- Adapters for external services

#### internal/config/

- Configuration management
- Environment variable handling
- Configuration validation

#### internal/pkg/

- Internal shared utilities
- Logging, metrics, and common functions
- Not exposed to other services

#### internal/server/

- Server setup and configuration
- Protocol-specific server implementations

### pkg/

- Public packages that can be imported by other services
- Client libraries and shared utilities
- Must maintain backward compatibility

### api/

- API definitions and specifications
- Protocol buffer definitions
- OpenAPI specifications

### deployments/

- Deployment configurations
- Docker and Kubernetes manifests
- Environment-specific configurations

### docs/

- Service documentation
- API documentation
- Architecture documentation

### scripts/

- Build and utility scripts
- Development tools
- CI/CD scripts

### test/

- Test files and test utilities
- Integration and end-to-end tests
- Test fixtures and mocks

## Naming Conventions

### Files

- Use lowercase with underscores for file names
- Use .go extension for Go files
- Use \_test.go suffix for test files
- Use .proto extension for protocol buffer files
- Use .yaml or .yml for YAML files

### Packages

- Use lowercase for package names
- Use singular form for package names
- Use descriptive names that indicate purpose

### Types and Interfaces

- Use PascalCase for type and interface names
- Use descriptive names that indicate purpose
- Use 'er' suffix for interfaces (e.g., Reader, Writer)

## Import Paths

- Use absolute imports from the module root
- Avoid relative imports
- Use consistent import ordering

## Documentation

- Each package should have a package comment
- Exported types and functions should have documentation
- Include examples for public APIs
- Keep documentation up to date

## Testing

- Unit tests should be in the same package
- Integration tests should be in the test directory
- Use table-driven tests where appropriate
- Include benchmarks for performance-critical code
