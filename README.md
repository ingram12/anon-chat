# Anon Chat Backend

Anonymous chat application backend with Proof-of-Work authentication.

## Project Structure

```
anon-chat-backend/
├── cmd/                   # Entry points
│   └── server/           # Main server application
├── internal/             # Internal application code
│   ├── app/              # Application initialization
│   ├── auth/             # Authentication logic
│   ├── chat/             # Chat functionality
│   ├── pow/              # Proof-of-Work implementation
│   ├── storage/          # Database interactions
│   ├── models/           # Data models
│   └── utils/            # Utility functions
├── pkg/                  # Reusable packages
├── configs/              # Configuration files
├── migrations/           # Database migrations
├── scripts/              # Deployment scripts
└── docs/                 # Documentation
```

## Getting Started

1. Install dependencies:
```bash
go mod tidy
```

2. Run the server:
```bash
go run cmd/server/main.go
```

## API Endpoints

### Proof-of-Work

- `GET /get-challenge` - Get a new Proof-of-Work challenge
- `POST /verify-challenge` - Verify a Proof-of-Work solution

#### Verify Challenge Request
```json
{
    "challenge": "string",
    "nonce": "string"
}
```

#### Verify Challenge Response
```json
{
    "success": true,
    "message": "Proof verified successfully"
}
``` 