# URL Shortener with Go, Gin, GORM, and PostgreSQL

This project is a simple URL shortening service built with Go, Gin, GORM, and PostgreSQL. It allows users to shorten URLs and track the number of times the shortened link is accessed.

## Features

* **Shorten any URL**: Enter a long URL and get a short link.
* **Track click counts**: Keep track of how many times a shortened URL is clicked.
* **Expiration for URLs**: Set a lifespan for the shortened URLs (24 hours).
* **Easy setup with Docker**: Quickly get up and running with Docker and Docker Compose.

## Tech Stack

* **Go** - Backend language
* **Gin** - Web framework for Go
* **GORM** - ORM library for Go
* **PostgreSQL** - Database for storing original and shortened URLs
* **Docker** - Containerization for easy deployment

## Installation

### Prerequisites

* **Go** installed on your machine (Go version 1.18 or higher).
* **PostgreSQL** (or Docker with a PostgreSQL container).
* **Docker** and **Docker Compose** for easy setup (optional).

### Clone the repository

```bash
git clone https://github.com/yourusername/url_shortener.git
cd url_shortener
```

### Set up PostgreSQL

Make sure PostgreSQL is running and create a database. You can use Docker for this:

```bash
docker run --name url_shortener_db -e POSTGRES_PASSWORD=password -d -p 5432:5432 postgres
```

Then, connect to PostgreSQL and create a database:

```bash
psql -h localhost -U postgres
CREATE DATABASE url_shortener;
```

### Configure the database connection

In the `db/db.go` file, make sure the connection string points to your PostgreSQL database:

```go
const (
    DB_USER     = "postgres"
    DB_PASSWORD = "password"
    DB_NAME     = "url_shortener"
    DB_HOST     = "localhost"
    DB_PORT     = "5432"
)
```

### Install dependencies

Run the following command to install the Go dependencies:

```bash
go mod download
```

### Migrate the database

To create the necessary tables in PostgreSQL, run the following command:

```bash
go run main.go
```

This will automatically create the `links` table and other necessary fields in your database.

### Start the application

Run the Go application:

```bash
go run main.go
```

The backend will now be running on `http://localhost:8080`.

## Docker Setup (Optional)

For a more convenient setup, you can run the backend and PostgreSQL using Docker Compose.

### Docker Compose

1. Ensure you have Docker and Docker Compose installed.
2. Inside the project directory, create a `docker-compose.yml` file:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:latest
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
      POSTGRES_DB: url_shortener
    ports:
      - "5432:5432"

  backend:
    build: .
    environment:
      DB_HOST: postgres
      DB_PORT: 5432
      DB_USER: postgres
      DB_PASSWORD: password
      DB_NAME: url_shortener
    ports:
      - "8080:8080"
    depends_on:
      - postgres
```

3. Build and start the services:

```bash
docker-compose up --build
```

This will start both the PostgreSQL database and the Go backend.

## Endpoints

### **POST /shorten**

Shorten a given URL.

**Request body**:

```json
{
  "original": "https://example.com"
}
```

**Response**:

```json
{
  "shortened": "short.ly/abc123"
}
```

### **GET /\:shortened**

Redirect to the original URL.

Example: `GET /abc123` will redirect to `https://example.com`.

### **GET /links**

Get all shortened links with their click counts.

**Response**:

```json
[
  {
    "id": 1,
    "original": "https://example.com",
    "shortened": "abc123",
    "clicks": 5,
    "created_at": "2025-05-05T15:00:00Z",
    "expires_at": "2025-05-06T15:00:00Z"
  }
]
```

## Contributing

Feel free to fork the project and submit pull requests! Contributions are always welcome.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
