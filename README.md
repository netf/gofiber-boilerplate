# Golang Fiber Boilerplate

[![CI](https://github.com/netf/gofiber-boilerplate/actions/workflows/ci.yml/badge.svg)](https://github.com/netf/gofiber-boilerplate/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/netf/gofiber-boilerplate)](https://goreportcard.com/report/github.com/netf/gofiber-boilerplate)

This is a boilerplate for building RESTful APIs using Golang and the Fiber web framework. It includes a basic Todo application structure with database integration, JWT authentication, and Swagger documentation.

## Features

- Fiber web framework
- PostgreSQL database with GORM
- JWT authentication
- Swagger API documentation
- Docker support
- Air for live reloading during development
- Structured logging with zerolog
- Configuration management with Viper
- Error handling and custom middleware
- API versioning

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose
- PostgreSQL

## Getting Started

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/gofiber-boilerplate.git
   cd gofiber-boilerplate
   ```

2. Set up environment variables:
   Copy the `.env.example` file to `.env` and update the values as needed.

3. Run the application using Docker:
   ```
   make docker-up
   ```
   This will start the application and a PostgreSQL database in Docker containers.

4. Access the API at `http://localhost:8080/api/v1`

5. View the Swagger documentation at `http://localhost:8080/swagger/`

## Development

To run the application in development mode with live reloading:
```
make run
```

## Testing

Run the tests with:
```
make test
```

## Building

To build the application:
```
make build
```

## API Endpoints

- `POST /api/v1/todos`: Create a new todo
- `GET /api/v1/todos`: List all todos
- `GET /api/v1/todos/:id`: Get a specific todo
- `PUT /api/v1/todos/:id`: Update a todo
- `DELETE /api/v1/todos/:id`: Delete a todo

For detailed API documentation, refer to the Swagger UI.

## Project Structure

- `cmd/`: Application entry point
- `config/`: Configuration management
- `internal/`: Internal application code
  - `api/`: API-related code
    - `handlers/`: HTTP request handlers
    - `middleware/`: Custom middleware
    - `routes/`: API route definitions
    - `utils/`: API utility functions
  - `db/`: Database connection and migrations
  - `models/`: Data models
  - `repositories/`: Data access layer
  - `services/`: Business logic
- `docs/`: Swagger documentation
- `migrations/`: Database migration files
- `docker/`: Docker-related files

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License.
