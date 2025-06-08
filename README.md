# Todo API Server

REST API сервер на Go для управления задачами с поддержкой MySQL.

## 🚀 Быстрый старт

### Требования
- Go 1.20+
- MySQL 8.0+

### Установка
1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/Mafit1/todo-app.git
   
2. Настройте MySQL БД:
   ```sql
   CREATE DATABASE todo_db;
   USE todo_db;

   CREATE TABLE todo (
     id INT AUTO_INCREMENT PRIMARY KEY,
     title VARCHAR(255) NOT NULL,
     completed BOOLEAN DEFAULT FALSE
   );
   
3. Создайе .env файл:
   ```ini
     DB_HOST=localhost
     DB_PORT=3306
     DB_USER=root
     DB_PASSWORD=yourpassword
     DB_NAME=todo_db

4. Запустите сервер:
   ```bash
     go run cmd/main.go

### Зависимости
[Echo](https://echo.labstack.com/) - HTTP фреймворк
[go-sql-driver/mysql](https://github.com/go-sql-driver/mysql) - Драйвер MySQL
[godotenv](https://github.com/joho/godotenv) - Загрузка .env файлов
