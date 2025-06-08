# Todo API Server

REST API server on Go for task management with MySQL support.

## Getting started

### Requirements
- Go 1.20+
- MySQL 8.0+

### Installation
1. Clone repository:
   ```
   git clone https://github.com/Mafit1/todo-app.git
   
2. Setup MySQL database:
   ```
   CREATE DATABASE todo_db;
   USE todo_db;

   CREATE TABLE todo (
     id INT AUTO_INCREMENT PRIMARY KEY,
     title VARCHAR(255) NOT NULL,
     completed BOOLEAN DEFAULT FALSE
   );
   
3. Create .env:
   ```
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=yourpassword
   DB_NAME=todo_db

4. Run server:
   ```
   go run cmd/main.go

### Dependencies
[Echo](https://echo.labstack.com/) - HTTP framework

[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) - MySQL driver

[godotenv](https://github.com/joho/godotenv) - Handling .env files
