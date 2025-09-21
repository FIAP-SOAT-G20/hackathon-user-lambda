# Build stage
FROM golang:1.24-alpine AS build

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application for Lambda (bootstrap is required for custom runtime)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bootstrap ./cmd/api

# Runtime stage using AWS Lambda base image
FROM public.ecr.aws/lambda/provided:al2023

# Copy the binary to the Lambda runtime directory
COPY --from=build /app/bootstrap ${LAMBDA_RUNTIME_DIR}/

# Set the handler
CMD [ "bootstrap" ]
