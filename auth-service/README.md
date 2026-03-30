# Authentication Service — Final Backend Implementation

This service implements a production-grade, secure authentication flow using **Go**, **Gin**, **GORM**, and **PostgreSQL**. It features a two-step login with OTP and secure session management via HttpOnly cookies.

## 🏗️ Project Structure
```text
auth-service/
├── main.go
├── config/
│   └── db.go
├── models/
│   └── user.go
├── handlers/
│   └── auth.go
├── middleware/
│   └── auth.go
├── utils/
│   ├── jwt.go
│   ├── otp.go
│   └── hash.go
```

## 📦 Final Backend APIs
All endpoints are prefixed with `/api/v1`.

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `POST` | `/register` | Create a new user account |
| `POST` | `/login` | Step 1: Validate credentials and send OTP |
| `POST` | `/verify-otp` | Step 2: Verify OTP and set session cookies |
| `GET` | `/me` | Get current authenticated user details |
| `POST` | `/refresh` | Refresh access token using refresh token cookie |
| `POST` | `/logout` | Clear session cookies and logout |

---

## 🚀 How to Run the App

### 1. Using Docker (Recommended)
This is the easiest way to start the app along with its database.
1.  Ensure **Docker Desktop** is running.
2.  Open your terminal in the `auth-service/` directory.
3.  Run:
    ```bash
    docker-compose up --build
    ```
    *The app will be available at `http://localhost:8081`.*

---

## 🧪 How to Test the App (Step-by-Step)

Follow these steps in a **separate terminal** to test the full authentication flow.

### Step 1: Register a User
```powershell
curl.exe -X POST http://localhost:8081/api/v1/register `
-H "Content-Type: application/json" `
-d '{"name":"Test User","email":"test@test.com","password":"password123"}'
```

### Step 2: Login (Triggers OTP)
```powershell
curl.exe -X POST http://localhost:8081/api/v1/login `
-H "Content-Type: application/json" `
-d '{"email":"test@test.com","password":"password123"}'
```
👉 **IMPORTANT:** Look at the terminal where Docker is running. You will see a line like `OTP: 123456`. Copy that code.

### Step 3: Verify OTP (Sets Cookies)
Replace `<OTP_CODE>` with the code from the logs.
```powershell
curl.exe -X POST http://localhost:8081/api/v1/verify-otp `
-H "Content-Type: application/json" `
-d '{"email":"test@test.com","otp":"<OTP_CODE>"}' `
-c cookies.txt
```
*The `-c cookies.txt` flag saves the secure session cookies to a file.*

### Step 4: Get Your Profile (Authorized)
```powershell
curl.exe http://localhost:8081/api/v1/me -b cookies.txt
```
*The `-b cookies.txt` flag sends the cookies back to the server for authentication.*

### Step 5: Logout
```powershell
curl.exe -X POST http://localhost:8081/api/v1/logout -b cookies.txt -c cookies.txt
```

---

## 💀 Final Result Summary
- **No Frontend Needed**: Fully testable via CLI.
- **Production-Style**: Secure HttpOnly cookies and JWT Refresh flow.
- **OTP Security**: Prevents unauthorized access even if password is leaked.
