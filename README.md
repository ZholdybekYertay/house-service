# DomService

DomService is a small full-stack home-service request management system built with Go, PostgreSQL and a lightweight HTML/CSS/JavaScript frontend.

The project demonstrates REST API development, authentication, authorization, database integration and a layered backend architecture.

## Features

- Create service requests
- View all requests as admin
- Update request status
- Delete requests
- User registration
- User login
- Password hashing with bcrypt
- JWT authentication
- User/admin roles
- Protected admin endpoints
- PostgreSQL database
- Docker Compose
- Graceful server shutdown
- Simple frontend

## Tech Stack

### Backend

- Go
- net/http
- pgx / pgxpool
- JWT
- bcrypt

### Database

- PostgreSQL

### Frontend

- HTML
- CSS
- JavaScript

### DevOps

- Docker
- Docker Compose

## Project Structure

```text
service/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── auth/
│   │   ├── jwt.go
│   │   └── password.go
│   │
│   ├── database/
│   │   └── postgres.go
│   │
│   ├── handler/
│   │   ├── auth_handler.go
│   │   └── request_handler.go
│   │
│   ├── middleware/
│   │   └── auth.go
│   │
│   ├── model/
│   │   ├── request.go
│   │   └── user.go
│   │
│   └── repository/
│       ├── request_repository.go
│       └── user_repository.go
│
├── web/
│   ├── index.html
│   ├── style.css
│   └── app.js
│
├── compose.yaml
├── .env.example
├── .gitignore
├── go.mod
├── go.sum
└── README.md