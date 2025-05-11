# Тестовое задание для практики в компании "Kvant"

REST API для управления пользователями и их заказами с использованием Golang и PostgreSQL.

## 🛠 Стек технологий
- **Язык программирования**: `Go` (Golang)
- **База данных**: `PostgreSQL`
- **ORM**: `GORM`
- **HTTP-фреймворк**: `Gin`
- **HTTP-библиотека**: `net/http`
- **Аутентификация**: JWT (библиотека `golang-jwt/jwt`)
- **Логирование**: Стандартная библиотека `log`
- **Контейнеризация**: `Docker` + `Docker Compose`
- **Документация**: `Swagger` (опционально)

---

## 🚀 Запуск проекта

### Предварительные требования
- Установленные **Docker** и **Docker Compose**
- **Go 1.20+**
- **Git**

### Инструкция (без Docker)
1. Клонируйте репозиторий:
    ```
    git clone https://github.com/PhosFactum/kvant-backend-practicum.github
    cd kvant-backend-practicum
    ```

2. Установите зависимости:
    ```
    go mod download
    ```

3. Сгенерируйте документацию и запустите сервер:
    ```
    swag init -g ./cmd/main.go --output ./docs
    go run ./cmd/main.go
    ```

4. Попробуйте отправить запрос для проверки (ниже пример):
    ```
    curl -i "http://localhost:8080/users page=2&limit=5&min_age=20&max_age=30"
    ```

5. Проверьте все хэндлеры по этому адресу (после запуска сервера):
    ```
    http://localhost:8080/swagger/index.html#/
    ```
