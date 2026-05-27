# Chirpy - Twitter-like Social Media API

A REST API backend for a Twitter-like social media application built with Go, PostgreSQL, and JWT authentication.

## Features

- **User Management**: Create users, update profile, manage passwords with Argon2id hashing
- **JWT Authentication**: Access tokens (1 hour) and refresh tokens (60 days) with HS256 signing
- **Chirps (Posts)**: Create, read, delete chirps with author filtering and sorting
- **Chirpy Red Membership**: Upgrade users via Polka webhooks with API key validation
- **Profanity Filtering**: Automatic censoring of inappropriate words
- **Secure Webhooks**: Polka integration with API key authentication

## Setup

### Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Goose (database migrations)
- SQLC (SQL code generation)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repo-url>
   cd chirpy
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Generate database code**
   ```bash
   sqlc generate
   ```

4. **Set up environment variables** (`.env` file)
   ```
   DB_URL=postgres://postgres:postgres@localhost:5432/chirpy?sslmode=disable
   PLATFORM=dev
   JWT_SECRET=<your-jwt-secret>
   POLKA_KEY=f271c81ff7084ee5b99a5091b42d486e
   ```

   Generate JWT_SECRET:
   ```bash
   openssl rand -base64 64
   ```

5. **Run database migrations**
   ```bash
   goose -dir ./sql/schema postgres "$DB_URL" up
   ```

6. **Start the server**
   ```bash
   go run main.go
   ```

   Server runs on `http://localhost:8080`

## API Endpoints

### Health Check

```
GET /api/healthz
```

Returns server health status.

### User Management

#### Create User
```
POST /api/users
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response: 201 Created
{
  "id": "uuid",
  "email": "user@example.com",
  "created_at": "2026-05-26T...",
  "updated_at": "2026-05-26T...",
  "is_chirpy_red": false
}
```

#### Login
```
POST /api/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

Response: 200 OK
{
  "id": "uuid",
  "email": "user@example.com",
  "token": "eyJ...",           (access token)
  "refresh_token": "abc123...", (refresh token)
  "is_chirpy_red": false,
  "created_at": "2026-05-26T...",
  "updated_at": "2026-05-26T..."
}
```

#### Update User
```
PUT /api/users
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "email": "newemail@example.com",
  "password": "newpassword123"
}

Response: 200 OK
{
  "id": "uuid",
  "email": "newemail@example.com",
  "is_chirpy_red": false,
  "created_at": "2026-05-26T...",
  "updated_at": "2026-05-26T..."
}
```

### Token Management

#### Refresh Access Token
```
POST /api/refresh
Authorization: Bearer <refresh_token>

Response: 200 OK
{
  "token": "eyJ..." (new access token)
}
```

#### Revoke Refresh Token
```
POST /api/revoke
Authorization: Bearer <refresh_token>

Response: 204 No Content
```

### Chirps (Posts)

#### Create Chirp
```
POST /api/chirps
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "body": "This is my first chirp!"
}

Response: 201 Created
{
  "id": "uuid",
  "body": "This is my first chirp!",
  "user_id": "uuid",
  "created_at": "2026-05-26T...",
  "updated_at": "2026-05-26T..."
}
```

#### Get All Chirps
```
GET /api/chirps
GET /api/chirps?author_id=<user-id>
GET /api/chirps?sort=asc
GET /api/chirps?sort=desc
GET /api/chirps?author_id=<user-id>&sort=desc

Response: 200 OK
[
  {
    "id": "uuid",
    "body": "...",
    "user_id": "uuid",
    "created_at": "2026-05-26T...",
    "updated_at": "2026-05-26T..."
  },
  ...
]
```

Query Parameters:
- `author_id` (optional): Filter chirps by author UUID
- `sort` (optional): Sort by created_at - `asc` (default) or `desc`

#### Get Chirp by ID
```
GET /api/chirps/{chirpID}

Response: 200 OK
{
  "id": "uuid",
  "body": "...",
  "user_id": "uuid",
  "created_at": "2026-05-26T...",
  "updated_at": "2026-05-26T..."
}
```

#### Delete Chirp
```
DELETE /api/chirps/{chirpID}
Authorization: Bearer <access_token>

Response: 204 No Content
```

Only the chirp author can delete their own chirps.

### Chirp Validation

#### Validate Chirp Content
```
POST /api/validate_chirp
Content-Type: application/json

{
  "body": "This is too long or contains bad words..."
}

Response: 200 OK
{
  "cleaned_body": "This is [redacted] or contains [redacted]..."
}

Response: 400 Bad Request (if > 140 characters)
{
  "error": "Chirp is too long"
}
```

### Chirpy Red

#### Upgrade User (Webhook)
```
POST /api/polka/webhooks
Authorization: ApiKey f271c81ff7084ee5b99a5091b42d486e
Content-Type: application/json

{
  "event": "user.upgraded",
  "data": {
    "user_id": "uuid"
  }
}

Response: 204 No Content
```

Other events are ignored with 204 response.

## Authentication

This API uses **JWT (JSON Web Token)** authentication:

1. **Login** to get access token + refresh token
2. **Use access token** for authenticated requests (1 hour validity)
3. **Refresh token** when access token expires (valid 60 days)
4. **Logout** by revoking refresh token

### How to Use Access Token

Include in request header:
```
Authorization: Bearer <access_token>
```

Example:
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Authorization: Bearer eyJ..." \
  -H "Content-Type: application/json" \
  -d '{"body": "Hello world"}'
```

## Database Schema

### Users Table
- `id` (UUID) - Primary key
- `email` (TEXT) - Unique email
- `hashed_password` (TEXT) - Argon2id hashed password
- `is_chirpy_red` (BOOLEAN) - Chirpy Red membership status
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Chirps Table
- `id` (UUID) - Primary key
- `user_id` (UUID) - Foreign key to users
- `body` (TEXT) - Chirp content
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Refresh Tokens Table
- `token` (TEXT) - Primary key
- `user_id` (UUID) - Foreign key to users
- `expires_at` (TIMESTAMP)
- `revoked_at` (TIMESTAMP) - NULL if not revoked
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

## Profanity Filter

The following words are automatically censored:
- kerfuffle
- sharbert
- fornax

Replace with `[redacted]` in responses.

## Error Handling

Standard HTTP status codes:
- `200 OK` - Success
- `201 Created` - Resource created
- `204 No Content` - Success with no response body
- `400 Bad Request` - Invalid request
- `401 Unauthorized` - Missing/invalid authentication
- `403 Forbidden` - Authenticated but not authorized
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

## Development

### Run Tests
```bash
bootdev run <test-id>
```

### Generate Code from SQL
```bash
sqlc generate
```

### Database Operations
```bash
# Run migrations
goose -dir ./sql/schema postgres "$DB_URL" up

# Check migration status
goose -dir ./sql/schema postgres "$DB_URL" status

# Rollback last migration
goose -dir ./sql/schema postgres "$DB_URL" down
```

## Project Structure

```
chirpy/
├── main.go                      # Server entry point
├── internal/
│   ├── handler/                 # HTTP handlers
│   │   ├── chirps.go
│   │   ├── users.go
│   │   ├── login.go
│   │   ├── webhook.go
│   │   └── ...
│   ├── auth/                    # Authentication logic
│   │   ├── jwt.go
│   │   ├── password.go
│   │   └── refreshToken.go
│   ├── database/                # Generated SQLC code
│   └── util/                    # Utility functions
├── sql/
│   ├── schema/                  # Database migrations
│   └── queries/                 # SQL queries
├── .env                         # Environment variables
└── README.md
```

## Security Notes

- Never commit `.env` files with real secrets
- JWT_SECRET should be a strong random string
- Passwords are hashed with Argon2id (industry standard)
- Refresh tokens are stored in database with expiration
- API keys for webhooks are required for sensitive operations
- All endpoints validate input and return appropriate error codes

## License

MIT License
