# Scrypts

A secure, encrypted notes application with Go backend and Next.js frontend.

## Overview

Scrypts is a production-ready full-stack web application for creating and managing encrypted notes. It features hardened security, JWT authentication, rate limiting, and a modern TypeScript frontend. Includes comprehensive security measures based on OWASP best practices and VAPT audit recommendations.

## Architecture

- **Backend**: Go with JWT auth, bcrypt password hashing, AES-GCM encryption, SQLite storage
- **Frontend**: Next.js 14 + TypeScript + React 18
- **Security**: Rate limiting, CORS whitelist, security headers, timing attack prevention
- **Communication**: REST API with strict CORS policy
- **Encryption**: Per-user AES-GCM encryption keys, server-side decryption

## Features

### Security (New!)
- **Rate Limiting**: 10 requests/minute per IP on auth endpoints
- **CORS Whitelist**: Strict origin validation with `ALLOWED_ORIGINS` env var
- **Security Headers**: HSTS, CSP, X-Frame-Options, X-Content-Type-Options, etc.
- **Entropy Validation**: Enforces strong secrets (32+ chars, 4.0+ bits/byte)
- **Timing Attack Prevention**: Constant-time operations in authentication
- **User Enumeration Prevention**: Generic errors and random delays
- **Username Validation**: Regex whitelist (alphanumeric, underscore, hyphen only)
- **Password Complexity**: Enforces uppercase, lowercase, digits, special chars (min 8 chars)
- **Bcrypt Cost 12**: Increased from default for stronger password hashing

### Backend
- User registration with bcrypt password hashing (cost: 12)
- JWT-based authentication with configurable expiry
- AES-GCM encryption for note content (server stores encrypted data)
- Server-side decryption for GET requests (plaintext response)
- Per-user encryption keys wrapped with master key
- SQLite persistence with WAL mode and foreign keys
- Rate-limited authentication endpoints (10 req/min per IP)
- CORS middleware with whitelist validation
- Security headers middleware (HSTS, CSP, X-Frame-Options, etc.)
- TLS/HTTPS support with configurable ports and HTTP→HTTPS redirect
- Input validation with regex patterns and ownership checks
- Efficient database queries with indexes

### Frontend
- TypeScript + React 18 with Next.js 14
- User registration and login UI
- Full CRUD interface for notes
- JWT token storage in localStorage
- Real-time note updates with edit mode
- Responsive design
- Error handling and user feedback
- No npm vulnerabilities (regularly updated)

## Quick Start

### Prerequisites

- Go 1.20+
- Node.js 18+ and npm/yarn
- Terminal access

### 1. Start the Backend

```bash
# From project root
cd /home/syko/go_stuff/scrypts

# Set required environment variables
export JWT_SECRET="$(openssl rand -base64 48)"
export MASTER_KEY="$(openssl rand -base64 48)"
export ALLOWED_ORIGINS="http://localhost:3000"

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

### Backend (Required)

**Critical - Application will not start without these:**

- `JWT_SECRET` - JWT signing key (min 32 chars, high entropy required)
  ```bash
  export JWT_SECRET="$(openssl rand -base64 48)"
  ```
- `MASTER_KEY` - Master encryption key for wrapping user keys (min 32 chars, high entropy required)
  ```bash
  export MASTER_KEY="$(openssl rand -base64 48)"
  ```

**Recommended:**

- `ALLOWED_ORIGINS` - Comma-separated list of allowed CORS origins (default: `http://localhost:3000,http://localhost:8080`)
  ```bash
  export ALLOWED_ORIGINS="http://localhost:3000,https://yourdomain.com"
  ```

**Optional:**

- `SCRYPTS_DB_PATH` - Database file path (default: `./scrypts.db`)
- `SCRYPTS_TLS_CERT` - Path to TLS certificate (optional)
- `SCRYPTS_TLS_KEY` - Path to TLS private key (optional)
- `SCRYPTS_HTTPS_PORT` - HTTPS port (default: `8443`)
- `SCRYPTS_HTTP_PORT` - HTTP port or redirector port (default: `8080`)

### Frontend

- `NEXT_PUBLIC_SCRYPTS_API` - Backend API URL (default: `http://localhost:8080`)

## Security Highlights

### Authentication & Authorization
- **Bcrypt password hashing** with cost factor 12 (increased from default)
- **JWT tokens** with configurable expiry
- **Username validation** with regex: `^[a-zA-Z0-9_-]{4,255}$`
- **Password complexity** requirements: min 8 chars, uppercase, lowercase, digit, special char
- **Rate limiting**: 10 requests/minute per IP on `/register` and `/login`
- **Timing attack prevention**: Constant-time operations, dummy hash for non-existent users
- **User enumeration prevention**: Generic error messages with random delays
- **Ownership verification** on all note operations

### Encryption
- **AES-256-GCM** authenticated encryption for all note content
- **Per-user encryption keys** derived from password using scrypt
- **User keys wrapped** with master key for secure storage
- **Server-side decryption** for GET requests (plaintext in response)
- **Nonces stored per-note** for GCM security

### Infrastructure
- **Security Headers**: HSTS, CSP, X-Frame-Options, X-Content-Type-Options, X-XSS-Protection
- **CORS Whitelist**: Strict origin validation (no permissive `*`)
- **Secret Validation**: Enforces 32+ character secrets with entropy checking (4.0+ bits/byte)
- **SQLite with WAL mode** for better concurrency
- **Foreign key constraints** for data integrity
- **Input validation** at HTTP layer with regex patterns
- **UUID validation** for note IDs
- **TLS 1.2+** with secure cipher preferences
- **HTTP to HTTPS redirect** support

## Project Structure

```
scrypts/
├── cmd/scrypts/
│   └── main.go              # Application entry point with middleware chain
├── internal/
│   ├── auth/
│   │   ├── handler.go       # Registration, login, JWT (with timing attack prevention)
│   │   └── password.go      # Bcrypt password hashing (cost: 12)
│   ├── config/
│   │   └── config.go        # Configuration with entropy validation
│   ├── middleware/
│   │   ├── cors.go          # CORS whitelist middleware
│   │   ├── security.go      # Security headers middleware (NEW)
│   │   └── ratelimit.go     # Rate limiting middleware (NEW)
│   ├── notes/
│   │   └── handler.go       # Notes CRUD handlers
│   ├── storage/
│   │   └── storage.go       # SQLite database layer with regex validation
│   └── utils/
│       └── crypto.go        # AES-GCM encryption utilities
├── frontend/
│   ├── pages/
│   │   ├── _app.tsx         # Next.js app wrapper
│   │   ├── index.tsx        # Login/register page
│   │   └── notes.tsx        # Notes CRUD interface
│   ├── styles/
│   │   └── globals.css      # Global styles
│   ├── package.json         # Frontend dependencies (Next.js 14)
│   ├── tsconfig.json        # TypeScript config
│   └── next.config.js       # Next.js configuration
├── go.mod                   # Go dependencies
├── go.sum                   # Go dependency checksums
└── README.md               # This file
```

## Production Deployment

### Backend

1. **Generate strong secrets** using cryptographically secure random generators:
   ```bash
   export JWT_SECRET="$(openssl rand -base64 48)"
   export MASTER_KEY="$(openssl rand -base64 48)"
   ```
   
2. **Configure CORS whitelist** for your production domains:
   ```bash
   export ALLOWED_ORIGINS="https://yourdomain.com,https://app.yourdomain.com"
   ```

3. **Use a reverse proxy** (nginx/Caddy/Traefik) for TLS termination

4. **Run as systemd service** with limited privileges

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
- Check that your origin is in the `ALLOWED_ORIGINS` environment variable
- For local dev: `export ALLOWED_ORIGINS="http://localhost:3000"`
- For production: Update `ALLOWED_ORIGINS` with your production domain
- Check browser console for specific origin issues

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

### End-to-End Tests
Run the included E2E test script:
```bash
./test.zsh
```

This will:
- Start a test server
- Register a user
- Login and obtain JWT
- Create, read, update, and delete notes
- Clean up test database

### Security Tests
Run the security test suite:
```bash
./test_security.zsh
```

This validates:
- ✅ Username validation (regex enforcement)
- ✅ Security headers (HSTS, CSP, X-Frame-Options, etc.)
- ✅ Rate limiting (10 req/min on auth endpoints)
- ✅ CORS whitelist policy
- ✅ Password complexity requirements

## Development

### Backend
```bash
# Run in development mode (requires env vars)
export JWT_SECRET="$(openssl rand -base64 48)"
export MASTER_KEY="$(openssl rand -base64 48)"
export ALLOWED_ORIGINS="http://localhost:3000"
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

### Completed ✅
- [x] Rate limiting on auth endpoints (10 req/min per IP)
- [x] Security headers middleware (HSTS, CSP, X-Frame-Options, etc.)
- [x] CORS whitelist with environment variable configuration
- [x] Entropy validation for secrets (32+ chars, 4.0+ bits/byte)
- [x] Timing attack prevention in authentication
- [x] User enumeration prevention
- [x] Username regex validation
- [x] Increased bcrypt cost to 12
- [x] Comprehensive security test suite

### Planned Enhancements
- [ ] CSRF protection middleware
- [ ] Implement refresh tokens for session management
- [ ] Add audit logging for authentication events
- [ ] Add health check endpoint
- [ ] Set up monitoring and logging (Prometheus, ELK)
- [ ] Add end-to-end tests with Playwright
- [ ] Implement password reset flow
- [ ] Add two-factor authentication
- [ ] Migrate to PostgreSQL for production scale
- [ ] Add real-time collaboration with WebSockets
- [ ] Implement note sharing and permissions

For detailed security improvements, see [SECURITY_IMPROVEMENTS.md](./SECURITY_IMPROVEMENTS.md).

## Contributing

Pull requests and issues are welcome! Please follow security best practices and code style guidelines.

### Guidelines
- Write tests for new features
- Follow Go and TypeScript best practices
- Document API changes
- Update README when adding features
- Run formatters (gofmt, prettier) before committing
- Review [SECURITY_AUDIT.md](./SECURITY_AUDIT.md) for security considerations

## Security

This project has undergone a comprehensive VAPT (Vulnerability Assessment and Penetration Testing) security audit. All critical and high-priority vulnerabilities have been addressed.

📄 **Security Documentation:**
- [SECURITY_AUDIT.md](./SECURITY_AUDIT.md) - Full VAPT audit report with 21 findings
- [SECURITY_IMPROVEMENTS.md](./SECURITY_IMPROVEMENTS.md) - Detailed implementation of security fixes

🔒 **Security Features:**
- Rate limiting (10 req/min per IP)
- CORS whitelist enforcement
- Security headers (HSTS, CSP, X-Frame-Options, etc.)
- Strong secret validation (32+ chars, entropy checking)
- Timing attack prevention
- User enumeration prevention
- Bcrypt cost 12
- Username regex validation

If you discover a security vulnerability, please email security@yourdomain.com (or open a private security advisory on GitHub).

## License

MIT
