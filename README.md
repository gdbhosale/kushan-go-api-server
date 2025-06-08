# Kushan GO API Server

A modern, production-ready REST API server built with Go, providing a robust foundation for Customer Relationship Management (CRM) systems. This API server features:

### Key Features
- **JWT Authentication** - Secure token-based authentication system
- **Role-Based Access Control (RBAC)** - Granular permission management
- **RESTful API Design** - Following REST best practices with standard HTTP methods
- **Swagger Documentation** - Complete API documentation with interactive Swagger UI
- **PostgreSQL Database** - Reliable data persistence with migrations
- **Middleware Support** - Built-in logging, rate limiting, and CORS
- **Environment Configuration** - Easy configuration via .env files
- **Hot Reload Development** - Fast development cycle with Air

### API Documentation

The API documentation is available via Swagger UI at `/swagger/index.html` when running the server locally. It provides detailed information about all available endpoints, request/response formats, and authentication requirements.

### Tech Stack
- Go 1.21+
- PostgreSQL
- JWT for authentication
- Swagger/OpenAPI for documentation
- Air for hot reloading
- Make for build automation

## TODO's

- Linting - golangci-lint - https://github.com/evilmartians/lefthook
- Testing
- Dependency injection at compile - google/wire

## Go Installation on Mac:

```
brew install go
```

## Project Setup

```
cd cmd/
go get

cp .env.local .env
```

Edit `.env` file for Database Credentials

## Run Project

You can run project by following options:

1. Run in development environment:

    ```
    make run-dev
    ```

2. Run & Watch in development environment (Using [Air](https://github.com/cosmtrek/air)):

    ```
    make run-watch
    ```

3. Build & Run `app` using `make`:

    ```
    make clean
    make build
    make run
    ```

## Utilities

Update Swagger Doc:

    ```
    make swagger
    ```




## References:
- https://github.com/swaggo/swag
- https://github.com/swaggo/http-swagger

## Required Libraries

- jwt-go - https://github.com/golang-jwt/jwt

## VSCode Extensions

- GO https://marketplace.visualstudio.com/items?itemName=golang.Go
- Code Spell Checker https://marketplace.visualstudio.com/items?itemName=streetsidesoftware.code-spell-checker
- indent-rainbow https://marketplace.visualstudio.com/items?itemName=oderwat.indent-rainbow