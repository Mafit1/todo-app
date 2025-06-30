# Todo API Server

REST API server on Go for task management with MySQL support and Redis caching.

## Getting started

### Requirements
- Go 1.20+
- MySQL 8.0+
- Redis 8.0.2+

### Setup
1. Clone repository:
   ```
   git clone https://github.com/Mafit1/todo-app.git
   cd todo-app
   ```

2. Edit .env.example and rename it to .env

3. Run Docker compose:
   ```
   docker-compose up --build
   # or
   docker-compose up --build -d
   ```

4. Logs (Optional):
   ```
   docker-compose logs -f app
   docker-compose logs -f db
   docker-compose logs -f redis
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

[go-redis](https://github.com/redis/go-redis) - Redis Go client

[godotenv](https://github.com/joho/godotenv) - Handling .env files
