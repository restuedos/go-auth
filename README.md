# Gin-Golang Backend Boilerplate

This is a boilerplate for building a backend API with Gin (a Golang web framework). It provides a structured foundation with authentication, Swagger API documentation, and Docker integration for scalability and maintainability.

## Features

- JWT-based authentication
- Swagger API documentation
- Modular project structure
- Database integration support
- Middleware for request handling
- Makefile for simplified commands
- Docker support for containerized deployment

## Default Admin Credentials

- **Email:** admin@example.com
- **Password:** admin123

## Folder Structure

```
project-root/
├── cmd/                         # Application entry point
├── docs/                        # Swagger API documentation
├── internal/
│   ├── delivery/
│   │   ├── http/
│   │   │   ├── handler/          # Request handlers
│   │   │   ├── handler/requests  # Request DTOs
│   │   │   ├── handler/responses # Response DTOs
│   │   │   ├── middleware/       # HTTP middlewares
│   │   │   ├── routes/           # API route definitions
│   ├── domain/
│   │   ├── entity/               # Business entities/models
│   │   ├── repository/           # Repository interfaces
│   ├── infrastructure/
│   │   ├── cache/                # Caching connection with Redis
│   │   ├── database/             # Database connection and setup
│   ├── usecase/                  # Business logic layer
├── tests/                        # Unit test
```

## Installation

### Prerequisites

- Golang installed (Go 1.24+ recommended)
- Docker installed (optional but recommended for deployment)
- Make (optional, for running commands easily)

### Setup

```sh
git clone https://github.com/restuedos/go-auth.git
cd go-auth
make init  # Install dependencies and prepare the project
```

## Running the Application

### Using Makefile

```sh
make run  # Start the application
```

### Using Go Commands

```sh
go run cmd/main.go
```

### Using Docker

```sh
make docker-build
make docker-run
```

## API Documentation

This boilerplate uses Swagger for API documentation. To generate and serve the documentation, run:

```sh
make doc
```

Then, open your browser and navigate to:

```
http://localhost:8080/swagger/index.html
```

## Authentication

- The API uses JWT for authentication.
- Users must obtain a token by logging in and then include it in the `Authorization` header for protected routes.

## Available Makefile Commands

| Command             | Description                                    |
| ------------------- | ---------------------------------------------- |
| `make run`          | Start the application using Air                |
| `make init`         | Install dependencies, generate docs, and build |
| `make build`        | Build the application binary                   |
| `make docker-build` | Build Docker image                             |
| `make docker-run`   | Run application using Docker Compose           |
| `make test`         | Run tests                                      |
| `make clean`        | Remove binaries and temporary files            |
| `make doc`          | Generate Swagger documentation                 |
| `make tidy`         | Clean up Go modules                            |

## License

This boilerplate is open-source and available under the MIT License.

---

Feel free to customize this README to match your project requirements!
