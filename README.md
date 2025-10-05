# 🚀 Attendance System - Backend

A robust and scalable Go backend for managing employee attendance with RESTful API endpoints.

![Go](https://img.shields.io/badge/Go-1.25.1-blue)
![Gin](https://img.shields.io/badge/Gin-Gonic-green)
![MySQL](https://img.shields.io/badge/MySQL-Database-orange)
![JWT](https://img.shields.io/badge/JWT-Auth-purple)

## 🚀 Features

### 🔐 Authentication & Authorization

- **JWT-based authentication** with secure token management
- **Role-based access control** (Admin, Manager, Employee)
- **Password hashing** with bcrypt
- **Token refresh** mechanism

### 👥 Employee Management

- **CRUD operations** for employee profiles
- **Department assignment** and management
- **Role management** system
- **Bulk operations** support

### 🏢 Department Management

- **Department CRUD operations**
- **Attendance rules configuration**
- **Late tolerance settings**
- **Working hour configurations**

### ⏰ Attendance Tracking

- **Clock-in/clock-out** functionality
- **Automatic status calculation** (Late, On Time, Absent)
- **Working hours tracking**
- **Attendance history** with filters

### 📊 Reporting & Analytics

- **Real-time attendance reports**
- **Department-wise analytics**
- **Excel export** functionality
- **Summary statistics**

### 🔧 Technical Features

- **RESTful API** design
- **CORS enabled** for frontend integration
- **Structured logging** with Zerolog
- **Input validation** with Go Validator
- **Error handling** middleware

## 🛠 Technology Stack

### Core Framework

- **Go 1.25.1** - Programming language
- **Gin Gonic** - HTTP web framework
- **GORM** - ORM library
- **MySQL** - Database

### Authentication & Security

- **JWT v5** - JSON Web Tokens
- **bcrypt** - Password hashing
- **CORS** - Cross-Origin Resource Sharing

### Utilities & Libraries

- **Zerolog** - Structured logging
- **Godotenv** - Environment configuration
- **Excelize** - Excel file manipulation
- **Validator v10** - Input validation

### Development Tools

- **Go Modules** - Dependency management

## ⚙️ Installation

### Prerequisites

- **Go 1.25.1** or higher
- **MySQL 8.0** or higher
- **Git**

### Setup Instructions

1. **Clone the repository**

```bash
git clone https://github.com/ImranNursandi/attendance-system-backend
```

### Setup Instructions

2. **Install dependencies**

```bash
go mod download
```

3. **Environment Configuration**
   Create a .env file from the template:

```bash
cp .env.example .env
```

4. **Configure environment variables**

```env

APP_ENV=development
APP_PORT=8080
GIN_MODE=debug

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=Password
DB_NAME=attendance_system

# JWT Configuration
JWT_SECRET=super-secret-jwt-key-here
JWT_EXPIRY=24h

# CORS Configuration
CORS_ALLOW_ORIGIN=*
CORS_ALLOW_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOW_HEADERS=*
```

5. **Run database migrations**

```bash
mysql -u root -p < database/migration.sql
```

6. **Start the development server**

```bash
go run main.go
```

The server will start on http://localhost:8080

### Base URL

```text
http://localhost:8080/api/v1
```
