# Тестовое задание для практики в компании "Kvant"

REST API для управления пользователями и их заказами с использованием Golang и PostgreSQL.

## 🛠 Стек технологий
- **Язык программирования**: Go 1.24+
- **База данных**: PostgreSQL 15+
- **ORM**: GORM
- **HTTP-фреймворк**: Gin
- **Аутентификация**: JWT (библиотека `golang-jwt/jwt`)
- **Контейнеризация**: Docker + Docker Compose
- **Документация**: Swagger (опционально)

---

## 🚀 Запуск через Docker Compose

### Предварительные требования
- Установленные **Docker** и **Docker Compose**
- Порт **5432** и **8080** не заняты

### Инструкция
1. Клонируйте репозиторий:
    ```bash
    git clone https://github.com/PhosFactum/kvant-backend-practicum.git
    cd kvant-backend-practicum
    ```

2. Создайте файл `.env` (шаблон):
    ```bash
    echo "DB_HOST=db
    DB_USER=kvant_user
    DB_PASSWORD=your_strong_password
    DB_NAME=kvant_db
    DB_PORT=5432
    JWT_SECRET=your_jwt_secret" > .env
    ```

3. Запустите сервисы:
    ```bash
    docker-compose up --build
    ```

4. Приложение будет доступно на `http://localhost:8080`  
   База данных — на `localhost:5432` (логин/пароль из `.env`)

5. Проверьте документацию Swagger:
    ```
    http://localhost:8080/swagger/index.html
    ```

---

## 🛠 Запуск без Docker

### Требования
- Установленные **Go 1.24+**, **PostgreSQL 15+**
- Созданная БД с параметрами из `.env`

### Инструкция
1. Настройте базу данных:
    ```sql
    CREATE DATABASE kvant_db;
    CREATE USER kvant_user WITH PASSWORD 'your_strong_password';
    GRANT ALL PRIVILEGES ON DATABASE kvant_db TO kvant_user;
    ```

2. Установите зависимости:
    ```bash
    go mod download
    ```

3. Запустите миграции (если требуется):
    ```bash
    go run cmd/migrate/main.go
    ```

4. Запустите сервер:
    ```bash
    go run cmd/main.go
    ```

5. Пример запроса:
    ```bash
    curl -X POST http://localhost:8080/users \
    -H "Content-Type: application/json" \
    -d '{"name":"John Doe", "email":"john@example.com", "age":30, "password":"secret"}'
    ```

---

## 🔒 Авторизация
- Для защищённых эндпоинтов требуется JWT-токен в заголовке:  
  `Authorization: Bearer <your_token>`
- Получить токен:
    ```bash
    curl -X POST http://localhost:8080/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"john@example.com", "password":"secret"}'

    ```
---

## 🐛 Устранение неполадок
- **Ошибка подключения к БД**: Проверьте `.env` и доступность PostgreSQL
- **Миграции не применяются**: Запустите вручную `go run cmd/migrate/main.go`
- **Swagger не генерируется**: Установите `swag` и выполните:
    ```bash
    swag init -g ./cmd/main.go --output ./docs
    ```
