# 🚗 SpotSync Server

A production-ready **Smart Parking & EV Charging Reservation REST API** built with **Go (Golang)**, **Echo**, **GORM**, **PostgreSQL**, and **Clean Architecture**.

This backend allows drivers to reserve parking and EV charging spots while enabling administrators to manage parking zones and monitor reservations efficiently.

---

## 🌐 Live API

**Base URL**

```text
https://spotsync-server.onrender.com
```

Health Check

```http
GET /
```

Example Response

```json
{
  "success": true,
  "message": "SpotSync API is running..."
}
```

---

# ✨ Features

### 🔐 Authentication & Authorization

* User Registration
* User Login
* JWT Authentication
* Role-Based Authorization (Driver & Admin)
* Password Hashing with bcrypt

### 🚗 Parking Zone Management

* Create Parking Zones (Admin)
* Update Parking Zones (Admin)
* Delete Parking Zones (Admin)
* View All Parking Zones
* View Single Parking Zone
* Dynamic Available Spot Calculation

### 📅 Reservation Management

* Reserve Parking Spot
* View My Reservations
* Cancel Reservation
* View All Reservations (Admin)

### ⚡ EV Charging Reservation

* Reserve limited EV charging spots
* Prevent over-capacity reservations
* Atomic reservation handling using Transactions

### 🔒 Concurrency Control

* Database Transactions
* Row-Level Locking (`FOR UPDATE`)
* Prevent Race Conditions
* Prevent Overbooking

### 🏛 Clean Architecture

* DTO Layer
* Handler Layer
* Service Layer
* Repository Layer
* Models Layer
* Dependency Injection

---

# 🛠 Tech Stack

| Technology | Description          |
| ---------- | -------------------- |
| Go 1.22+   | Programming Language |
| Echo       | HTTP Framework       |
| GORM       | ORM                  |
| PostgreSQL | Database             |
| JWT        | Authentication       |
| bcrypt     | Password Hashing     |
| Validator  | Request Validation   |
| Render     | Deployment           |

---

# 📂 Project Structure

```text
spotsync-server/
│
├── cmd/
│   └── main.go
│
├── internal/
│   ├── config/
│   ├── dto/
│   ├── handler/
│   ├── middleware/
│   ├── models/
│   ├── repository/
│   ├── routes/
│   ├── server/
│   ├── service/
│   └── utils/
│
├── go.mod
├── go.sum
└── README.md
```

---

# 🔑 API Endpoints

## Authentication

| Method | Endpoint                |
| ------ | ----------------------- |
| POST   | `/api/v1/auth/register` |
| POST   | `/api/v1/auth/login`    |

---

## Parking Zones

| Method | Endpoint            |
| ------ | ------------------- |
| POST   | `/api/v1/zones`     |
| GET    | `/api/v1/zones`     |
| GET    | `/api/v1/zones/:id` |

---

## Reservations

| Method | Endpoint                               |
| ------ | -------------------------------------- |
| POST   | `/api/v1/reservations`                 |
| GET    | `/api/v1/reservations/my-reservations` |
| DELETE | `/api/v1/reservations/:id`             |
| GET    | `/api/v1/reservations`                 |

---

# 🚀 Getting Started

## Clone Repository

```bash
git clone https://github.com/Uttamdevsharma/spotsync-server.git
```

```bash
cd spotsync-server
```

---

## Install Dependencies

```bash
go mod tidy
```

---

## Environment Variables

Create a `.env` file in the project root.

```env
PORT=8080

DATABASE_URL=your_postgresql_connection_string

JWT_SECRET=your_super_secret_key

JWT_EXPIRATION=24h

BCRYPT_COST=10
```

---

## Run the Project

```bash
go run ./cmd
```

The server will start at

```text
http://localhost:8080
```

---

# 📦 Build

```bash
go build -o app ./cmd
```

Run

```bash
./app
```

---

# 🔒 Security

* JWT Authentication
* bcrypt Password Hashing
* Request Validation
* Role-Based Authorization
* Clean Layer Separation
* Secure Password Storage

---

# ⚙️ Business Rules

* Only authenticated users can create reservations.
* Only administrators can manage parking zones.
* Drivers can only cancel their own reservations.
* Available spots are calculated dynamically.
* Parking zones cannot exceed their maximum capacity.
* Reservation creation is protected using transactions and row-level locking.

---

# 📄 HTTP Status Codes

| Code | Description           |
| ---- | --------------------- |
| 200  | OK                    |
| 201  | Created               |
| 400  | Bad Request           |
| 401  | Unauthorized          |
| 403  | Forbidden             |
| 404  | Not Found             |
| 409  | Conflict              |
| 500  | Internal Server Error |

---

# 👨‍💻 Author

**Uttam Kumar Dev Sharma**

Backend Developer

GitHub: https://github.com/Uttamdevsharma
