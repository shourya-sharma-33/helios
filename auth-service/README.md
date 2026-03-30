# Authentication Service API Documentation

This service handles user registration and authentication using Go, Gin, PostgreSQL, and Redis. It features a two-step registration (email verification) and a two-step login (OTP via email).

## Project Structure
```text
auth-service/
├── main.go
├── go.mod
├── config/
│   ├── db.go
│   └── redis.go
├── models/
│   └── user.go
├── controllers/
│   └── auth_controller.go
├── routes/
│   └── auth_routes.go
├── middleware/
│   └── auth_middleware.go
├── utils/
│   ├── token.go
│   ├── jwt.go
│   └── mailer.go
```

## Base URL
`http://localhost:8081`

## Endpoints

### 1. User Registration
Initiates registration by sending a verification link (printed to console in dev).

- **URL:** `/register`
- **Method:** `POST`
- **Headers:** `Content-Type: application/json`
- **Request Body:**
  ```json
  {
    "name": "Full Name",
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **Success Response:**
  - **Code:** 200 OK
  - **Content:** `{"success": true, "message": "Verification email sent"}`

---

### 2. Verify Email
Completes registration and saves user to database.

- **URL:** `/verify/:token`
- **Method:** `GET`
- **Success Response:**
  - **Code:** 200 OK
  - **Content:** `{"success": true, "message": "User verified"}`

---

### 3. User Login (Step 1: Send OTP)
Authenticates credentials and sends a 6-digit OTP to the user's email.

- **URL:** `/login`
- **Method:** `POST`
- **Headers:** `Content-Type: application/json`
- **Request Body:**
  ```json
  {
    "email": "user@example.com",
    "password": "password123"
  }
  ```
- **Success Response:**
  - **Code:** 200 OK
  - **Content:** `{"message": "If email is valid, OTP sent (valid 5 min)"}`

---

### 4. Verify OTP (Step 2: Complete Login)
Verifies the OTP and returns Access and Refresh tokens.

- **URL:** `/verify-otp`
- **Method:** `POST`
- **Headers:** `Content-Type: application/json`
- **Request Body:**
  ```json
  {
    "email": "user@example.com",
    "otp": "123456"
  }
  ```
- **Success Response:**
  - **Code:** 200 OK
  - **Content:**
    ```json
    {
      "message": "Login successful",
      "access_token": "ACCESS_JWT",
      "refresh_token": "REFRESH_JWT",
      "user": {
        "id": "1",
        "email": "user@example.com",
        "name": "Full Name"
      }
    }
    ```

---

### 5. Protected Route
Example of a route protected by JWT middleware.

- **URL:** `/protected/`
- **Method:** `GET`
- **Headers:** `Authorization: Bearer <ACCESS_JWT>`
- **Success Response:**
  - **Code:** 200 OK
  - **Content:** `{"message": "Protected route"}`

## Database Schema
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    password TEXT NOT NULL
);
```

## Requirements
- **PostgreSQL**: Database `helios` must exist.
- **Redis**: Must be running on `localhost:6379`.

## Setup & Running

### Using Docker (Recommended)
This will set up the Auth Service, PostgreSQL, and Redis automatically.

1. Ensure Docker Desktop is running.
2. From the `auth-service/` directory, run:
   ```bash
   docker-compose up --build
   ```
3. The API will be available at `http://localhost:8081`.

### Manual Setup
1. `go mod tidy` to install dependencies.
2. Ensure Postgres and Redis are running locally.
3. Update database credentials in `config/db.go`.
4. `go run .`
