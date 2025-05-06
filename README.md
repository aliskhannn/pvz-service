# 🏪 PVZ Service

Сервис для управления пунктами выдачи заказов (ПВЗ), реализованный на Go с использованием чистой архитектуры.

## 🚀 Функциональность

- 🔐 Аутентификация и авторизация пользователей (JWT)
- 🏪 Управление ПВЗ (создание, получение, обновление, удаление)
- 📦 Управление товарами (добавление, удаление)
- 📝 Управление приемками (создание, закрытие)
- 🐳 Docker-контейнеризация
- 📦 Миграции базы данных (goose)

## 📁 Структура проекта

```
.
├── cmd/
│   └── server/
│       └── main.go          # Точка входа в приложение
├── internal/
│   ├── auth/               # Логика аутентификации
│   ├── delivery/
│   │   └── http/          # HTTP-хендлеры
│   ├── repository/        # Репозитории для работы с БД
│   └── usecase/          # Бизнес-логика
├── pkg/
│   └── middleware/       # Промежуточное ПО
├── migrations/          # Миграции базы данных
├── docker-compose.yaml  # Конфигурация Docker
├── Dockerfile          # Сборка приложения
└── .env.example       # Пример переменных окружения
```

## ⚙️ Настройка окружения

Создайте файл `.env` на основе `.env.example`:

```env
# Сервер
HTTP_PORT=8080
GRPC_PORT=3000
METRICS_PORT=9000

# База данных PostgreSQL
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=your_db_name

# JWT токен
JWT_SECRET=your_very_long_and_secure_secret_here

# Goose миграции
GOOSE_DRIVER=postgres
GOOSE_DBSTRING=postgres://<db_user>:<db_password>@localhost:5432/<db_name>
GOOSE_MIGRATION_DIR=./migrations
GOOSE_TABLE=custom.goose_migrations
```

## 🐳 Запуск через Docker

1. Соберите и запустите контейнеры:

```bash
docker-compose up --build
```

2. Применение миграций:

```bash
docker-compose exec app goose -dir migrations up
```

3. Сервис будет доступен по адресу: `http://localhost:8080`

## 🧪 Запуск тестов

### Unit-тесты

```bash
go test ./... -v
```

### Интеграционные тесты

```bash
go test -tags=integration ./... -v
```

## 🔐 Авторизация

Сервис использует JWT для аутентификации. Процесс авторизации:

1. Пользователь отправляет запрос на `/login` с учетными данными
2. Сервер проверяет учетные данные и возвращает JWT-токен
3. Клиент включает токен в заголовок `Authorization: Bearer <token>`
4. Middleware проверяет токен и добавляет пользователя в контекст

Пример запроса:

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password"}'
```

## 📝 API Endpoints

### Аутентификация
- `POST /login` - Вход в систему
- `POST /register` - Регистрация

### ПВЗ
- `POST /pvz` - Создание ПВЗ
- `GET /pvz/{id}` - Получение ПВЗ
- `PUT /pvz/{id}` - Обновление ПВЗ
- `POST /pvz/{id}` - Удаление ПВЗ

### Товары
- `POST /products` - Добавление товара
- `POST /products/{pvzId}/delete_last_product` - Удаление последнего товара

### Приемки
- `POST //receptions` - Создание приемки
- `PUT /receptions/{id}/close_last_reception` - Закрытие приемки

## 👤 Автор

Aliskhan Khutiev
