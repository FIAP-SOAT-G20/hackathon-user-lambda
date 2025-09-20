# Hackathon User Service (AWS Lambda + DynamoDB + JWT)

This microservice provides JWT authentication and basic user management (register, login, me) on AWS Lambda with
DynamoDB, written in Go, following clean architecture principles with hexagonal architecture pattern.

## Features

- User registration with email validation
- JWT-based authentication
- User profile retrieval
- Clean architecture with separated concerns
- Comprehensive test coverage with mocks
- DynamoDB integration with sequential ID generation

## API Endpoints

### POST /users/register
Register a new user
- **Request**: `{ "name": string, "email": string, "password": string }`
- **Success (201)**: `{ "userId": number, "name": string, "email": string }`
- **Error (409)**: `{ "error": "email already registered" }`
- **Error (400)**: `{ "error": "invalid input" }`

### POST /users/login
Authenticate user and receive JWT token
- **Request**: `{ "email": string, "password": string }`
- **Success (200)**: `{ "token": string }`
- **Error (401)**: `{ "error": "invalid credentials" }`

### GET /users/me
Get current user profile (requires authentication)
- **Header**: `Authorization: Bearer <jwt-token>`
- **Success (200)**: `{ "userId": number, "name": string, "email": string }`
- **Error (401)**: `{ "error": "missing/invalid bearer token" }`
- **Error (404)**: `{ "error": "user not found" }`

## Requirements

- Go >= 1.21 (recommended 1.22)
- AWS CLI (for deployment)
- Make (for build automation)

## Project Structure

```
├── cmd/api/                    # Lambda entry point
├── internal/
│   ├── core/                   # Business logic layer
│   │   ├── domain/             # Domain entities
│   │   ├── dto/                # Data transfer objects
│   │   ├── port/               # Interface definitions
│   │   │   └── mocks/          # Generated mocks
│   │   └── usecase/            # Business use cases
│   ├── adapter/                # External adapters
│   │   ├── controller/         # HTTP handlers
│   │   └── presenter/          # Response formatting
│   └── infrastructure/         # Infrastructure layer
│       ├── config/             # Configuration management
│       ├── datasource/         # Database implementation
│       ├── logger/             # Logging utilities
│       └── pkg/                # Shared utilities
└── dist/                       # Build artifacts
```

## Environment Variables

The Lambda function expects the following environment variables:

- **USERS_TABLE_NAME**: DynamoDB table name for users
- **IDS_TABLE_NAME**: DynamoDB table name for sequential ID generation
- **AWS_REGION**: AWS region
- **JWT_SECRET**: HMAC secret for JWT signing/verification
- **JWT_EXPIRATION**: Token expiration duration (e.g., "24h")

For local development, you can create a `.env` file:

```env
JWT_SECRET=your-secure-secret-here
AWS_REGION=us-east-1
USERS_TABLE_NAME=hackathon-users-local
IDS_TABLE_NAME=hackathon-ids-local
JWT_EXPIRATION=24h
```

## Development

### Building

Build the Lambda function for Linux:

```bash
make build
```

This creates a `bootstrap` executable in the `dist/` directory.

### Packaging

Create a deployment package:

```bash
make package
```

This creates `dist/function.zip` ready for Lambda deployment.

### Testing

Run all tests:

```bash
make test
```

Run tests with coverage:

```bash
make coverage
```

### Code Generation

Generate mocks for testing:

```bash
make mock
```

### Code Formatting

```bash
make fmt
```

### Clean Build Artifacts

```bash
make clean
```

## Deployment

1. **Build and package** the function:
   ```bash
   make package
   ```

2. **Deploy manually** using AWS CLI or AWS Console:
   - Upload `dist/function.zip` to your Lambda function
   - Set the handler to `bootstrap`
   - Configure environment variables
   - Set runtime to `provided.al2023`

3. **Set up DynamoDB tables** with the following structure:

   **Users Table**:
   - Primary Key: `userId` (Number)
   - GSI: `email_index` with `email` (String) as key

   **IDs Table**:
   - Primary Key: `sequence` (String)
   - Used for generating sequential user IDs

## Architecture

This project follows Clean Architecture principles:

- **Domain Layer**: Core business entities and rules
- **Use Case Layer**: Application-specific business logic
- **Interface Adapters**: Controllers, presenters, and gateways
- **Infrastructure**: External concerns (database, HTTP, logging)

The hexagonal architecture ensures:
- Business logic independence from external frameworks
- Easy testing with comprehensive mocks
- Flexibility to change external dependencies

## Security Considerations

- Never commit secrets to version control
- Use AWS Secrets Manager or SSM Parameter Store for JWT_SECRET in production
- Passwords are hashed using bcrypt with salt
- JWT tokens have configurable expiration
- Input validation on all endpoints
