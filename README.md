# Todo API Server

REST API server on Go for task management with MySQL support.

## Getting started

### Requirements
- Go 1.20+
- MySQL 8.0+

### Setup
1. Clone repository:
   ```
   git clone https://github.com/Mafit1/todo-app.git
   ```
   
2. Setup MySQL database:
   ```
   CREATE DATABASE todo_db;
   USE todo_db;

   CREATE TABLE todo (
     id INT AUTO_INCREMENT PRIMARY KEY,
     title VARCHAR(255) NOT NULL,
     completed BOOLEAN DEFAULT FALSE
   );
   ```
   
3. Setup Redis (Docker example)
   ```
   docker run -d --name redis -p 6379:6379 redis redis-server --requirepass your_password_here
   ```

4. Edit .env.example and rename it to .env

5. Run server:
   ```
   go run cmd/main.go
   ```

### API Endpoints

| Method | Endpoint | Description | Sample request body |
| ------ | -------- | ----------- | ------------------- |
| POST   | /todo   | Create task | {"title":"Task 1","completed":false} |
| GET   | /todo   | Get all tasks | - |
| GET   | /todo/:id   | Get task with id | - |
| PUT   | /todo/:id   | Complete task update | {"title":"Updated","completed":true} |
| PATCH   | /todo/:id   | Partial task update | {"completed":true} |
| DELETE   | /todo/:id   | Delete task | - |

### Dependencies
[Echo](https://echo.labstack.com/) - HTTP framework

[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) - MySQL driver

[godotenv](https://github.com/joho/godotenv) - Handling .env files
