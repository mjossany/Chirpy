# Chirpy 🐦

A Twitter-like social media API built with Go and PostgreSQL. Chirpy allows users to create accounts, post short messages (chirps), and interact with a clean, RESTful API.

## Features

### Core Functionality
- **User Management**: Registration, login, profile updates
- **Chirps**: Create, read, list, and delete short messages (max 140 characters)
- **Authentication**: JWT-based auth with refresh tokens
- **Content Moderation**: Built-in profanity filter
- **Premium Subscriptions**: Chirpy Red upgrade system via Polka webhooks
- **Admin Dashboard**: Metrics and system administration

### Security Features
- bcrypt password hashing
- JWT access tokens (1 hour expiry)
- Refresh tokens (60 days expiry)
- API key authentication for webhooks
- Bearer token authentication for protected endpoints

## Tech Stack

- **Backend**: Go 1.24.5
- **Database**: PostgreSQL 13
- **Authentication**: JWT with refresh tokens
- **ORM**: SQLC for type-safe SQL
- **Containerization**: Docker & Docker Compose
- **Password Hashing**: bcrypt

## Quick Start

### Prerequisites
- Go 1.24.5+
- Docker & Docker Compose
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/mjossany/Chirpy.git
   cd Chirpy
   ```

2. **Start the database**
   ```bash
   docker-compose up -d
   ```

3. **Set up environment variables**
   Create a `.env` file in the root directory:
   ```env
   DB_URL=postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable
   PLATFORM=dev
   JWT_SECRET=your-jwt-secret-here
   POLKA_KEY=your-polka-api-key-here
   ```

4. **Run database migrations**
   ```bash
   # Install goose for migrations (if needed)
   go install github.com/pressly/goose/v3/cmd/goose@latest
   
   # Run migrations
   goose -dir sql/schema postgres "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
   ```

5. **Install dependencies**
   ```bash
   go mod download
   ```

6. **Run the application**
   ```bash
   go run .
   ```

The server will start on `http://localhost:8080`

## API Documentation

### Authentication

All protected endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <jwt_token>
```

Webhook endpoints require an API key:
```
Authorization: ApiKey <api_key>
```

### Endpoints

#### Health Check
- **GET** `/api/healthz` - Health check endpoint

#### User Management
- **POST** `/api/users` - Create a new user
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```

- **PUT** `/api/users` - Update user information (requires auth)
  ```json
  {
    "email": "newemail@example.com",
    "password": "newpassword123"
  }
  ```

#### Authentication
- **POST** `/api/login` - User login
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```

- **POST** `/api/refresh` - Refresh access token
  ```json
  {
    "token": "refresh_token_here"
  }
  ```

- **POST** `/api/revoke` - Revoke refresh token
  ```json
  {
    "token": "refresh_token_here"
  }
  ```

#### Chirps
- **GET** `/api/chirps` - List all chirps (supports sorting and filtering)
- **GET** `/api/chirps/{chirpID}` - Get a specific chirp
- **POST** `/api/chirps` - Create a new chirp (requires auth)
  ```json
  {
    "body": "This is my first chirp!"
  }
  ```
- **DELETE** `/api/chirps/{chirpID}` - Delete a chirp (requires auth, owner only)

#### Webhooks
- **POST** `/api/polka/webhooks` - Polka payment webhook (requires API key)

#### Admin
- **GET** `/admin/metrics` - View admin dashboard with hit metrics
- **POST** `/admin/reset` - Reset application metrics

#### Static Files
- **GET** `/app/*` - Serve static files from the root directory

## Database Schema

### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    is_chirpy_red BOOLEAN NOT NULL DEFAULT false
);
```

### Chirps Table
```sql
CREATE TABLE chirp (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    body TEXT NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);
```

### Refresh Tokens Table
```sql
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    revoked_at TIMESTAMP
);
```

## Project Structure

```
Chirpy/
├── main.go                      # Application entry point
├── go.mod                       # Go module dependencies
├── docker-compose.yml           # Database container setup
├── sqlc.yaml                    # SQLC configuration
├── json.go                      # JSON response helpers
├── handle_*.go                  # HTTP route handlers
├── handler_*.go                 # Additional handlers
├── internal/
│   ├── auth/                    # Authentication utilities
│   │   ├── jwt.go              # JWT token management
│   │   └── hash.go             # Password hashing
│   └── database/               # Database layer (SQLC generated)
│       ├── models.go           # Database models
│       ├── db.go               # Database connection
│       └── *.sql.go            # Generated query functions
├── sql/
│   ├── schema/                 # Database migrations
│   │   └── *.sql               # Migration files
│   └── queries/                # SQL queries for SQLC
│       └── *.sql               # Query definitions
└── assets/
    └── logo.png                # Application assets
```

## Development

### Database Migrations

Using Goose for database migrations:

```bash
# Create a new migration
goose -dir sql/schema create migration_name sql

# Run migrations
goose -dir sql/schema postgres $DB_URL up

# Rollback migrations
goose -dir sql/schema postgres $DB_URL down
```

### Code Generation

This project uses SQLC for type-safe database queries:

```bash
# Install SQLC
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate code from SQL queries
sqlc generate
```

### Content Moderation

The application includes a basic profanity filter that replaces the following words with "****":
- kerfuffle
- sharbert
- fornax

### Testing

Run tests for the authentication module:
```bash
go test ./internal/auth/...
```

## Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `DB_URL` | PostgreSQL connection string | Yes | - |
| `PLATFORM` | Platform identifier | Yes | - |
| `JWT_SECRET` | Secret key for JWT signing | Yes | - |
| `POLKA_KEY` | API key for Polka webhooks | Yes | - |

## API Response Formats

### Success Response
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

### Error Response
```json
{
  "error": "Error message describing what went wrong"
}
```

## Chirpy Red Premium

Users can upgrade to Chirpy Red premium status through the Polka payment integration. When a user upgrades, Polka sends a webhook to `/api/polka/webhooks` with the event type `user.upgraded`.

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is part of a learning exercise and is not intended for production use.

## Contact

For questions or support, please open an issue on GitHub. 