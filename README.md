# gin-auth-example

![Go Version](https://img.shields.io/badge/Go%20Version-latest-brightgreen.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Depado/gin-auth-example)](https://goreportcard.com/report/github.com/Depado/gin-auth-example)

A simple example demonstrating how to implement cookie-based session 
authentication in Go using Gin. This project showcases basic authentication 
flows including login, logout, and protected routes.

## Features

- Cookie-based session management
- Protected routes using middleware
- Basic authentication flow (login/logout)
- Session persistence using encrypted cookies

## Quick Start

```bash
# Clone the repository
git clone https://github.com/depado/gin-auth-example
cd gin-auth-example

# Run the server
go run main.go
```

The server will start on `localhost:8080`

## API Endpoints

### Public Routes

- `POST /login`: Authenticate user
  - Body: `{"username": "hello", "password": "itsme"}`
- `GET /logout`: End user session

### Protected Routes (requires authentication)

- `GET /private/me`: Get current user information
- `GET /private/status`: Get login status

## Authentication Flow

1. Send a POST request to `/login` with credentials
2. On successful login, a session cookie is set
3. Use this cookie for subsequent requests to protected routes
4. Call `/logout` to end the session

## Testing

Run the test suite:

```bash
go test -v
```

## Security Note

This is a demonstration project. For production use:
- Replace the hard-coded secret key
- Use secure password hashing
- Implement proper user storage
- Use HTTPS
