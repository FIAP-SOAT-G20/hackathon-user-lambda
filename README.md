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

**Request:**

```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Success (201):**

```json
{
  "user_id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Error (409):**

```json
{
  "error": "email already registered"
}
```

**Error (400):**

```json
{
  "error": "invalid input"
}
```

### POST /users/login
Authenticate user and receive JWT token

**Request:**

```json
{
  "email": "john@example.com",
  "password": "password123"
}
```

**Success (200):**

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error (401):**

```json
{
  "error": "invalid credentials"
}
```

### GET /users/me
Get current user profile (requires authentication)

**Header:** `Authorization: Bearer <jwt-token>`

**Success (200):**

```json
{
  "user_id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

**Error (401):**

```json
{
  "error": "missing bearer token"
}
```

**Error (404):**

```json
{
  "error": "user not found"
}
```

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