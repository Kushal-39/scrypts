# Scrypts

A secure, encrypted notes application with Go backend and Next.js frontend.

## Overview

Scrypts is a full-stack web application for creating and managing encrypted notes. It features end-to-end encryption, JWT authentication, and a modern TypeScript frontend. Perfect for learning secure web development or as a foundation for production applications.

## Architecture

- **Backend**: Go with JWT auth, AES-GCM encryption, SQLite storage
- **Frontend**: Next.js 13 + TypeScript + React
- **Communication**: REST API with CORS enabled
- **Security**: Scrypt password hashing, per-user encryption keys, server-side decryption

## Features

### Backend
- User registration with scrypt password hashing
- JWT-based authentication with configurable expiry
- AES-GCM encryption for note content (server stores encrypted data)
- Server-side decryption for GET requests (plaintext response)
- Per-user encryption keys wrapped with master key
- SQLite persistence with WAL mode and foreign keys
- CORS middleware for frontend integration
- TLS/HTTPS support with configurable ports
- Input validation and ownership checks
- Efficient database queries with indexes

### Frontend
- TypeScript + React with Next.js 13
- User registration and login UI
- Full CRUD interface for notes
- JWT token storage in localStorage
- Real-time note updates
- Responsive design
- Error handling and user feedback

## Quick Start

### Prerequisites

- Go 1.20+
- Node.js 18+ and npm/yarn
- Terminal access

### 1. Start the Backend

```bash
# From project root
cd /home/syko/go_stuff/scrypts

# Build the backend
go build -o scrypts ./cmd/scrypts

# Run the backend (HTTP mode for local dev)
./scrypts
```

The backend will start on `http://localhost:8080`.

### 2. Start the Frontend

In a new terminal:

```bash
# Navigate to frontend directory
cd /home/syko/go_stuff/scrypts/frontend

# Install dependencies (first time only)
npm install

# Run development server
npm run dev
```

The frontend will start on `http://localhost:3000`.

### 3. Use the App

1. Open [http://localhost:3000](http://localhost:3000) in your browser
2. Click **Register** to create a new account
3. Enter a username and password
4. Click **Login** to authenticate
5. Create, edit, and delete encrypted notes!

## API Endpoints

### Authentication
- `POST /register` — Register a new user
  - Body: `{"username": "user", "password": "pass"}`
  - Response: `201 Created` or error message

- `POST /login` — Login and receive JWT token
  - Body: `{"username": "user", "password": "pass"}`
  - Response: `{"token": "jwt_token_here"}`

### Notes (Protected - requires JWT)
- `POST /notes` — Create a new encrypted note
  - Header: `Authorization: Bearer <token>`
  - Body: `{"content": "note text"}`
  - Response: `{"id": "note-uuid"}`

- `GET /notes` — List all notes for authenticated user
  - Header: `Authorization: Bearer <token>`
  - Response: `[{"id": "...", "content": "...", "created": ..., "modified": ...}]`

- `PUT /notes` — Update a note
  - Header: `Authorization: Bearer <token>`
  - Body: `{"id": "note-uuid", "content": "updated text"}`
  - Response: `{"status": "updated"}`

- `DELETE /notes` — Delete a note
  - Header: `Authorization: Bearer <token>`
  - Body: `{"id": "note-uuid"}`
  - Response: `{"status": "deleted"}`

## Environment Variables

### Backend

- `SCRYPTS_DB_PATH` - Database file path (default: `./scrypts.db`)
- `SCRYPTS_TLS_CERT` - Path to TLS certificate (optional)
- `SCRYPTS_TLS_KEY` - Path to TLS private key (optional)
- `SCRYPTS_HTTPS_PORT` - HTTPS port (default: `8443`)
- `SCRYPTS_HTTP_PORT` - HTTP port or redirector port (default: `8080`)

### Frontend

- `NEXT_PUBLIC_SCRYPTS_API` - Backend API URL (default: `http://localhost:8080`)

## Security Highlights

### Encryption
- AES-256-GCM authenticated encryption for all note content
- Per-user encryption keys derived from password using scrypt
- User keys wrapped with master key for secure storage
- Server-side decryption for GET requests (plaintext in response)
- Nonces stored per-note for GCM security

### Authentication
- JWT tokens with configurable expiry
- Scrypt password hashing (N=32768, r=8, p=1)
- Username validation and length limits
- Ownership verification on all note operations

### Infrastructure
- SQLite with WAL mode for better concurrency
- Foreign key constraints for data integrity
- Input validation at HTTP layer
- UUID validation for note IDs
- CORS configuration for production deployment
- TLS 1.2+ with secure cipher preferences
- HTTP to HTTPS redirect support

## Project Structure

```
scrypts/
├── cmd/scrypts/
│   └── main.go              # Application entry point
├── internal/
│   ├── auth/
│   │   ├── handler.go       # Registration, login, JWT handling
│   │   └── password.go      # Password hashing with scrypt
│   ├── config/
│   │   └── config.go        # Configuration and secrets
│   ├── middleware/
│   │   └── cors.go          # CORS middleware
│   ├── notes/
│   │   └── handler.go       # Notes CRUD handlers
│   ├── storage/
│   │   └── storage.go       # SQLite database layer
│   └── utils/
│       └── crypto.go        # AES-GCM encryption utilities
├── frontend/
│   ├── pages/
│   │   ├── _app.tsx         # Next.js app wrapper
│   │   ├── index.tsx        # Login/register page
│   │   └── notes.tsx        # Notes CRUD interface
│   ├── styles/
│   │   └── globals.css      # Global styles
│   ├── package.json         # Frontend dependencies
│   ├── tsconfig.json        # TypeScript config
│   └── next.config.js       # Next.js configuration
├── go.mod                   # Go dependencies
├── go.sum                   # Go dependency checksums
└── README.md               # This file
```

## Production Deployment

### Backend

1. **Set secrets via environment** (JWT secret, master wrap key)
2. **Use a reverse proxy** (nginx/Caddy/Traefik) for TLS termination
3. **Run as systemd service** with limited privileges
4. **Configure CORS** to allow only your frontend domain in `internal/middleware/cors.go`
5. **Set up DB backups** and monitoring
6. **Enable logging** and metrics collection
7. **Use production-grade secrets manager** (AWS Secrets Manager, Vault, etc.)

Example systemd service:
```ini
[Unit]
Description=Scrypts API Server
After=network.target

[Service]
Type=simple
User=scrypts
WorkingDirectory=/opt/scrypts
ExecStart=/opt/scrypts/scrypts
Environment=SCRYPTS_DB_PATH=/var/lib/scrypts/scrypts.db
Restart=on-failure

[Install]
WantedBy=multi-user.target
```

### Frontend

1. **Build production bundle**: `npm run build`
2. **Deploy to Vercel, Netlify**, or serve with `npm start`
3. **Set `NEXT_PUBLIC_SCRYPTS_API`** to your production backend URL
4. **Configure CDN** and caching for static assets
5. **Enable production optimizations** in Next.js config

### TLS/HTTPS

For local testing with self-signed certificates:
```bash
openssl req -x509 -newkey rsa:4096 -nodes -days 365 \
  -keyout key.pem -out cert.pem \
  -subj "/CN=localhost" \
  -addext "subjectAltName=DNS:localhost,IP:127.0.0.1"

export SCRYPTS_TLS_CERT=/path/to/cert.pem
export SCRYPTS_TLS_KEY=/path/to/key.pem
export SCRYPTS_HTTPS_PORT=8443
./scrypts
```

For production, use Let's Encrypt via reverse proxy or certbot.

## Troubleshooting

### CORS errors in browser
- Ensure backend is running with CORS middleware enabled
- Check browser console for specific origin issues
- For production, update `internal/middleware/cors.go` with your frontend domain

### "Failed to fetch notes"
- Verify backend is running on port 8080
- Check JWT token is stored in localStorage (browser dev tools → Application)
- Ensure user is logged in and token hasn't expired

### TypeScript errors in frontend
- Run `npm install` to ensure all dependencies are installed
- Check `tsconfig.json` is present in frontend directory
- Clear Next.js cache: `rm -rf .next && npm run dev`

### Port 443 permission denied
- Ports < 1024 require root or special capabilities
- Use `sudo setcap 'cap_net_bind_service=+ep' ./scrypts`
- Or run behind a reverse proxy on privileged ports

## Testing

Run the included end-to-end test script:
```bash
./test.zsh
```

This will:
- Start a test server
- Register a user
- Login and obtain JWT
- Create, read, update, and delete notes
- Clean up test database

## Development

### Backend
```bash
# Run in development mode
go run ./cmd/scrypts

# Build for production
go build -o scrypts ./cmd/scrypts

# Run tests
go test ./...

# Check for errors
go vet ./...
```

### Frontend
```bash
cd frontend

# Development server with hot reload
npm run dev

# Production build
npm run build
npm start

# Type checking
npx tsc --noEmit
```

## Next Steps

- [ ] Add rate limiting on auth endpoints
- [ ] Implement refresh tokens for session management
- [ ] Add HSTS and security headers middleware
- [ ] Set up monitoring and logging (Prometheus, ELK)
- [ ] Add end-to-end tests with Playwright
- [ ] Implement password reset flow
- [ ] Add two-factor authentication
- [ ] Migrate to PostgreSQL for production scale
- [ ] Add real-time collaboration with WebSockets
- [ ] Implement note sharing and permissions

## Contributing

Pull requests and issues are welcome! Please follow security best practices and code style guidelines.

### Guidelines
- Write tests for new features
- Follow Go and TypeScript best practices
- Document API changes
- Update README when adding features
- Run formatters (gofmt, prettier) before committing

## License

MIT
