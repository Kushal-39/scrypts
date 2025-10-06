# Scrypts

A secure-by-design notes application built in Go.

## Features

- User registration and login with bcrypt password hashing
- JWT-based authentication for all protected routes
- Notes CRUD: create, read, update, delete
- AES-256-GCM encryption for note content
- Key derivation using scrypt
- SQLite database for persistent storage
- Input validation and security checks
- Rate limiting and password complexity enforcement
- Ready for Docker and cloud deployment
- CI/CD pipeline with security scanning

## API Endpoints

- POST /register — Register a new user
- POST /login — Login and receive a JWT token
- POST /notes — Create a new encrypted note
- GET /notes — List all notes for the authenticated user
- GET /notes/{id} — Read a specific note
- PUT /notes/{id} — Update a note
- DELETE /notes/{id} — Delete a note

## Security Highlights

- Zero trust architecture: all sensitive operations require authentication
- No secrets or keys stored in code or repository
- HTTPS-ready for production deployments
- Secrets managed via environment variables or secret managers
- Audit logging and input sanitization

## Getting Started

1. Clone the repo
2. Set environment variables (JWT_SECRET, DB_PATH, etc.)
3. Run the server:
   go run cmd/scrypts/main.go
4. Test endpoints with curl or Postman

## Docker

- Build and run with Docker Compose
- Secrets passed securely via environment variables

## Contributing

Pull requests and issues are welcome!
Please follow security best practices and code style guidelines.

## License

MIT
