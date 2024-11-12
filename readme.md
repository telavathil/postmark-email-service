# Postmark Email Service

A Go service that handles email sending via Postmark API. Supports both AWS Lambda deployment and local HTTP server execution.

## Features

- Send emails via Postmark API
- Input validation
- Runs as AWS Lambda function or local HTTP server
- Environment-based configuration
- JSON request/response format

## Prerequisites

- Go 1.22 or later
- Postmark account and server token
- AWS account (for Lambda deployment)
- Docker (optional)

## Installation

Clone the repository and install dependencies:
```bash
git clone https://github.com/yourusername/postmark-service.git
cd postmark-service
go mod download
```

## Configuration

Create a `.env` file in the project root:
```bash
POSTMARK_SERVER_TOKEN=your_token_here
PORT=8080  # Optional, defaults to 8080
```

## Usage

### Running Locally

Build and run the server:
```bash
go build
./postmark-service
```

The server will start on http://localhost:8080

### API Endpoint

**POST /send-email**

Request:
```json
{
    "from": "sender@example.com",
    "to": ["recipient@example.com"],
    "subject": "Test Email",
    "htmlBody": "<p>Hello World</p>",
    "textBody": "Hello World"
}
```

Response:
```json
{
    "success": true,
    "message": "Email sent successfully"
}
```

### AWS Lambda Deployment

Build and package for AWS Lambda:
```bash
# Build for Lambda
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

# Create deployment package
zip function.zip bootstrap
```

## Project Structure
```
.
├── config/
│   └── config.go         # Configuration management
├── models/
│   └── email.go          # Data models
├── server/
│   ├── server.go         # Server setup
│   ├── handlers.go       # Request handlers
│   └── responses.go      # Response formatting
├── main.go               # Application entry point
├── go.mod               # Go module file
└── .env                 # Environment variables (not in repo)
```

## Development

Run tests and format code:
```bash
# Run tests
go test ./...

# Format code
go fmt ./...
```
