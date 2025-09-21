# Hackathon User Service

[![CI: Unit Tests](https://github.com/FIAP-SOAT-G20/hackathon-user-lambda/actions/workflows/ci-unit-test.yml/badge.svg)](https://github.com/FIAP-SOAT-G20/hackathon-user-lambda/actions/workflows/ci-unit-test.yml)
[![CI: Integration Tests](https://github.com/FIAP-SOAT-G20/hackathon-user-lambda/actions/workflows/ci-integration-test.yml/badge.svg)](https://github.com/FIAP-SOAT-G20/hackathon-user-lambda/actions/workflows/ci-integration-test.yml)
[![CI: Security Scan](https://github.com/FIAP-SOAT-G20/hackathon-user-lambda/actions/workflows/ci-govulncheck.yml/badge.svg)](https://github.com/FIAP-SOAT-G20/hackathon-user-lambda/actions/workflows/ci-govulncheck.yml)
[![CI: Go Lint](https://github.com/FIAP-SOAT-G20/hackathon-user-lambda/actions/workflows/ci-golangci-lint.yml/badge.svg)](https://github.com/FIAP-SOAT-G20/hackathon-user-lambda/actions/workflows/ci-golangci-lint.yml)

A modern, production-ready microservice for user authentication and management built with Go, AWS Lambda, and DynamoDB.
This service implements clean architecture principles with hexagonal architecture pattern, providing JWT-based
authentication with comprehensive security features.

## 🚀 Features

- **JWT Authentication**: Secure token-based authentication with configurable expiration
- **User Management**: Complete user lifecycle (register, login, profile retrieval)
- **Clean Architecture**: Hexagonal architecture with separated concerns and dependency injection
- **Production Ready**: Comprehensive test coverage, CI/CD pipelines, and Docker support
- **AWS Native**: Built for AWS Lambda with DynamoDB integration
- **Security First**: Password hashing, input validation, and secure JWT handling
- **Observable**: Structured logging with request tracing
- **Scalable**: Serverless architecture with automatic scaling

## 📋 API Endpoints

### Authentication & User Management

| Method | Endpoint               | Description                         | Auth Required |
|--------|------------------------|-------------------------------------|---------------|
| `POST` | `/prod/users/register` | Register a new user                 | ❌             |
| `POST` | `/prod/users/login`    | Authenticate user and get JWT token | ❌             |
| `GET`  | `/prod/users/me`       | Get current user profile            | ✅             |

### POST /prod/users/register

Register a new user with email validation and password hashing.

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response (201 Created):**
```json
{
  "user_id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Error Responses:**

- `400 Bad Request`: Invalid input or validation errors
- `409 Conflict`: Email already registered

### POST /prod/users/login

Authenticate user credentials and receive a JWT token.

**Request:**
```json
{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses:**

- `400 Bad Request`: Invalid input
- `401 Unauthorized`: Invalid credentials

### GET /prod/users/me

Retrieve current user profile information.

**Headers:**

```
Authorization: Bearer <jwt-token>
```

**Response (200 OK):**
```json
{
  "user_id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Error Responses:**

- `401 Unauthorized`: Missing or invalid token
- `404 Not Found`: User not found

## 🏗️ Architecture

### Clean Architecture Layers

```
┌─────────────────────────────────────────────────────────────┐
│                    cmd/api (Entry Point)                    │
├─────────────────────────────────────────────────────────────┤
│  internal/adapter/          │  internal/adapter/           │
│  controller/                 │  presenter/                  │
│  (HTTP Handlers)            │  (Response Formatting)       │
├─────────────────────────────────────────────────────────────┤
│                internal/core/usecase/                       │
│                (Business Logic)                             │
├─────────────────────────────────────────────────────────────┤
│  internal/core/domain/      │  internal/core/dto/          │
│  (Entities)                 │  (Data Transfer Objects)      │
├─────────────────────────────────────────────────────────────┤
│                internal/core/port/                          │
│               (Interface Definitions)                       │
├─────────────────────────────────────────────────────────────┤
│  internal/infrastructure/                                   │
│  (External Dependencies: Database, Auth, Config)           │
└─────────────────────────────────────────────────────────────┘
```

### Project Structure

```
├── cmd/api/                         # Application entry point
│   └── main.go                      # Lambda handler setup
├── internal/
│   ├── adapter/                     # External interface adapters
│   │   ├── controller/              # HTTP request handlers
│   │   │   ├── user_controller.go
│   │   │   └── user_controller_test.go
│   │   └── presenter/               # Response formatting
│   │       └── json_presenter.go
│   ├── core/                        # Business logic layer
│   │   ├── domain/                  # Domain entities
│   │   │   └── user.go
│   │   ├── dto/                     # Data transfer objects
│   │   │   └── user_dto.go
│   │   ├── port/                    # Interface definitions
│   │   │   ├── mocks/              # Generated test mocks
│   │   │   ├── jwt_signer_port.go
│   │   │   ├── presenter_port.go
│   │   │   ├── user_controller_port.go
│   │   │   ├── user_repository_port.go
│   │   │   └── user_usecase_port.go
│   │   └── usecase/                 # Business use cases
│   │       ├── user_usecase.go
│   │       ├── user_usecase_test.go
│   │       └── user_usecase_suite_test.go
│   └── infrastructure/              # Infrastructure layer
│       ├── auth/                    # JWT implementation
│       │   └── jwt.go
│       ├── config/                  # Configuration management
│       │   └── config.go
│       ├── datasource/              # Database implementation
│       │   └── dynamodb_user_repository.go
│       └── logger/                  # Logging utilities
│           ├── logger.go
│           └── pretty_handler.go
├── .github/workflows/               # CI/CD pipelines
├── dist/                           # Build artifacts (generated)
├── Dockerfile                      # Multi-stage container build
├── Makefile                        # Build automation
├── go.mod                          # Go module definition
└── *.json                         # AWS resource templates
```

## 🛠️ Technology Stack

- **Language**: Go 1.25
- **Runtime**: AWS Lambda (Custom Runtime)
- **Database**: Amazon DynamoDB
- **Authentication**: JWT with HMAC-SHA256
- **Security**: bcrypt password hashing
- **Testing**: Go testing with testify and gomock
- **CI/CD**: GitHub Actions
- **Containerization**: Docker with multi-stage builds

## 📦 Dependencies

### Core Dependencies

- `github.com/aws/aws-lambda-go` - AWS Lambda runtime
- `github.com/aws/aws-sdk-go-v2` - AWS SDK v2
- `github.com/golang-jwt/jwt/v5` - JWT implementation
- `golang.org/x/crypto` - Cryptographic functions

### Development Dependencies

- `github.com/stretchr/testify` - Testing framework
- `go.uber.org/mock` - Mock generation
- `github.com/fatih/color` - Colored output for development

## ⚙️ Configuration

### Environment Variables

| Variable           | Description                 | Example               | Required |
|--------------------|-----------------------------|-----------------------|----------|
| `USERS_TABLE_NAME` | DynamoDB users table name   | `hackathon-users`     | ✅        |
| `IDS_TABLE_NAME`   | DynamoDB ID sequence table  | `hackathon-ids`       | ✅        |
| `AWS_REGION`       | AWS region                  | `us-east-1`           | ✅        |
| `JWT_SECRET`       | HMAC secret for JWT signing | `your-256-bit-secret` | ✅        |
| `JWT_EXPIRATION`   | Token expiration duration   | `24h`                 | ✅        |

### Local Development (.env)

```env
JWT_SECRET=your-secure-256-bit-secret-here
AWS_REGION=us-east-1
USERS_TABLE_NAME=hackathon-users-local
IDS_TABLE_NAME=hackathon-ids-local
JWT_EXPIRATION=24h
```

## 🚀 Quick Start

### Prerequisites

- Go 1.25 or later
- AWS CLI configured
- Make (for build automation)
- Docker (optional, for containerized builds)

### Development Workflow

1. **Clone and setup:**
   ```bash
   git clone https://github.com/FIAP-SOAT-G20/hackathon-user-lambda.git
   cd hackathon-user-lambda
   cp .env.example .env  # Configure your environment
   ```

2. **Install dependencies:**
   ```bash
   go mod download
   ```

3. **Generate mocks for testing:**
   ```bash
   make mock
   ```

4. **Run tests:**
   ```bash
   make test           # Run all tests
   make coverage       # Run tests with coverage report
   ```

5. **Build for deployment:**
   ```bash
   make build          # Build Linux binary
   make package        # Create deployment package
   ```

### Available Make Commands

| Command         | Description                       |
|-----------------|-----------------------------------|
| `make build`    | Build Lambda binary for Linux     |
| `make package`  | Create ZIP deployment package     |
| `make test`     | Run all tests with race detection |
| `make coverage` | Generate test coverage report     |
| `make mock`     | Generate mocks for interfaces     |
| `make fmt`      | Format Go code                    |
| `make clean`    | Clean build artifacts             |

## 🚀 Deployment

### AWS Lambda Deployment

1. **Build the deployment package:**
   ```bash
   make package
   ```

2. **Deploy using AWS CLI:**
   ```bash
   aws lambda create-function \
     --function-name hackathon-user-service \
     --runtime provided.al2023 \
     --role arn:aws:iam::ACCOUNT:role/lambda-execution-role \
     --handler bootstrap \
     --zip-file fileb://dist/function.zip \
     --environment Variables='{
       "USERS_TABLE_NAME":"hackathon-users",
       "IDS_TABLE_NAME":"hackathon-ids",
       "AWS_REGION":"us-east-1",
       "JWT_SECRET":"your-secret",
       "JWT_EXPIRATION":"24h"
     }'
   ```

3. **Update existing function:**
   ```bash
   aws lambda update-function-code \
     --function-name hackathon-user-service \
     --zip-file fileb://dist/function.zip
   ```

### Docker Deployment

```bash
# Build Docker image
docker build -t hackathon-user-service .

# Run locally for testing
docker run -p 8080:8080 \
  -e USERS_TABLE_NAME=users \
  -e IDS_TABLE_NAME=ids \
  -e JWT_SECRET=test-secret \
  hackathon-user-service
```

### DynamoDB Setup

Create required DynamoDB tables:

**Users Table:**

```json
{
  "TableName": "hackathon-users",
  "KeySchema": [
    {
      "AttributeName": "userId",
      "KeyType": "HASH"
    }
  ],
  "AttributeDefinitions": [
    {
      "AttributeName": "userId",
      "AttributeType": "N"
    },
    {
      "AttributeName": "email",
      "AttributeType": "S"
    }
  ],
  "GlobalSecondaryIndexes": [
    {
      "IndexName": "email_index",
      "KeySchema": [
        {
          "AttributeName": "email",
          "KeyType": "HASH"
        }
      ],
      "Projection": {
        "ProjectionType": "ALL"
      }
    }
  ]
}
```

**IDs Table:**

```json
{
  "TableName": "hackathon-ids",
  "KeySchema": [
    {
      "AttributeName": "sequence",
      "KeyType": "HASH"
    }
  ],
  "AttributeDefinitions": [
    {
      "AttributeName": "sequence",
      "AttributeType": "S"
    }
  ]
}
```

## 🔄 CI/CD Pipeline

The project includes comprehensive GitHub Actions workflows:

- **Unit Tests**: Automated testing with Go 1.25
- **Integration Tests**: End-to-end testing scenarios
- **Security Scanning**: Vulnerability detection with govulncheck
- **Linting**: Code quality checks with golangci-lint
- **Build & Deploy**: Automated Docker builds and ECR pushes

Workflows are triggered on pull requests and main branch pushes.

## 🧪 Testing

### Test Coverage

The project maintains high test coverage with comprehensive test suites:

- **Unit Tests**: Business logic and use case testing
- **Integration Tests**: End-to-end API testing
- **Mock Tests**: Interface contract verification

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make coverage

# Run specific test package
go test ./internal/core/usecase/...

# Run with verbose output
go test -v ./internal/...
```

## 🔐 Security

- **Password Security**: bcrypt hashing with salt
- **JWT Security**: HMAC-SHA256 signing with configurable expiration
- **Input Validation**: Comprehensive request validation
- **Dependency Scanning**: Automated vulnerability detection
- **Secure Headers**: Standard security headers in responses

## 📈 Performance

- **Serverless**: Auto-scaling with AWS Lambda
- **Optimized Builds**: Minimal binary size with build optimizations
- **Connection Pooling**: Efficient DynamoDB connection management
- **Structured Logging**: Minimal performance impact logging

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

### Code Standards

- Follow Go best practices and idioms
- Maintain test coverage above 80%
- Use conventional commit messages
- Ensure all CI checks pass

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Team

**FIAP SOAT G20** - Hackathon User Service Team

---

Built with ❤️ using Go and AWS serverless technologies.