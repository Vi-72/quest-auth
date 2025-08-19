# Quest Auth Service

HTTP-сервис аутентификации с поддержкой регистрации пользователей, авторизации и управления профилями.

## ✨ Основные возможности

- 🔐 **Аутентификация**: регистрация, вход и управление сессиями
- 👤 **Управление пользователями**: создание, получение и обновление профилей
- 📧 **Валидация контактов**: проверка email и телефонов с форматированием
- 🔒 **Безопасность**: хеширование паролей с bcrypt
- 🔄 **Domain Events**: отслеживание пользовательских действий
- 🏗️ **Clean Architecture**: четкое разделение слоев и ответственности
- ⚡ **Оптимизированная БД**: индексы для быстрого поиска

## 🔧 Запуск

### 📦 Требования
- Go 1.23+
- PostgreSQL

### 🚀 Быстрый старт

1. **Настройка переменных окружения:**
```bash
cp config.example .env
# Отредактируйте .env файл под вашу конфигурацию
```

2. **Запуск:**
```bash
go run ./cmd/app
```

Сервер запускается на порту, указанном в переменной `HTTP_PORT` (по умолчанию 8080).

### 🌐 API Endpoints

#### Аутентификация
- `POST /api/v1/auth/register` - Регистрация нового пользователя
- `POST /api/v1/auth/login` - Вход в систему

#### Пользователи  
- `GET /api/v1/users/{user_id}` - Получение профиля пользователя

#### Система
- `GET /health` - Проверка состояния сервиса

### 📖 Примеры использования API

#### Регистрация пользователя
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "phone": "+1234567890", 
    "name": "John Doe",
    "password": "securepassword123"
  }'
```

#### Вход в систему
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "securepassword123"
  }'
```

#### Получение профиля
```bash
curl -X GET http://localhost:8080/api/v1/users/{user_id}
```

### 🏗️ Структура проекта

```
quest-auth/
├── cmd/                    # 🚀 Точка входа
│   ├── app/                # Главное приложение
│   ├── composition_root.go # DI контейнер
│   └── config.go           # Конфигурация
├── internal/               # 🏗️ Основной код приложения
│   ├── adapters/           # Адаптеры (Hexagonal Architecture)
│   │   ├── in/http/        # HTTP handlers & validations
│   │   └── out/postgres/   # Репозитории БД
│   │       ├── userrepo/   # User repository
│   │       └── eventrepo/  # Event store
│   ├── core/               # Бизнес-логика (DDD)
│   │   ├── application/    # Use cases & handlers
│   │   │   └── usecases/   # Auth use cases
│   │   │       └── auth/   # Register, Login, GetUser
│   │   ├── domain/         # Доменная модель
│   │   │   └── model/      # Domain objects
│   │   │       ├── auth/   # User aggregate
│   │   │       └── kernel/ # Shared value objects (Email, Phone)
│   │   └── ports/          # Интерфейсы
│   └── pkg/                # Общие пакеты
│       ├── ddd/            # DDD building blocks
│       └── errs/           # Error types
├── tests/                  # 🧪 Все тесты проекта
├── configs/                # ⚙️ Конфигурационные файлы
└── Makefile                # 🛠️ Команды для разработки
```

### 🎯 Доменная модель

**User (Пользователь)** - Aggregate Root
- ID, Email, Phone, Name
- Password Hash (bcrypt)
- Timestamps (CreatedAt, UpdatedAt)
- Методы: Register, Login, ChangePhone, ChangeName, SetPassword
- Domain Events (UserRegistered, UserLoggedIn, UserPhoneChanged, etc.)

**Email (Электронная почта)** - Value Object
- Валидация формата через regex
- Нормализация (trim, lowercase)
- Обеспечение уникальности

**Phone (Телефон)** - Value Object  
- Валидация международного формата (+1234567890)
- Обеспечение уникальности
- Поддержка различных стран

### 🔐 Безопасность

1. **Хеширование паролей**: bcrypt с настраиваемой сложностью
2. **Валидация входных данных**: многоуровневая система проверок
3. **Уникальность контактов**: проверка email/phone при регистрации
4. **Защита от SQL injection**: используется GORM ORM
5. **Структурированные ошибки**: RFC 7807 Problem Details

### 🔄 Обработка ошибок

Система использует **правильные HTTP коды**:

```json
// Ошибка валидации (400)
{
  "type": "bad-request",
  "title": "Bad Request", 
  "status": 400,
  "detail": "validation failed: field 'email' is required"
}

// Неверные учетные данные (401)
{
  "type": "unauthorized",
  "title": "Unauthorized",
  "status": 401, 
  "detail": "invalid credentials"
}

// Пользователь не найден (404)
{
  "type": "not-found",
  "title": "Not Found",
  "status": 404,
  "detail": "user with ID 'uuid' not found"
}
```

### ⚡ Производительность

#### 🗂️ Индексы БД

```sql
-- Уникальность и быстрый поиск
CREATE UNIQUE INDEX idx_users_email ON users(email);
CREATE UNIQUE INDEX idx_users_phone ON users(phone);

-- Поиск по ID
CREATE INDEX idx_users_id ON users(id);

-- События
CREATE INDEX idx_events_aggregate_id ON events(aggregate_id);
CREATE INDEX idx_events_type ON events(event_type);
```

## 🔄 Domain-Driven Design

### 🏗️ Паттерны

- **Aggregate Root**: User с инкапсуляцией бизнес-логики
- **Value Objects**: Email, Phone с валидацией  
- **Domain Events**: отслеживание пользовательских действий
- **Unit of Work**: атомарные транзакции
- **Repository**: абстракция над хранилищем

### 📡 События

```go
// Автоматически создаются при изменениях
UserRegistered{UserID, Email, Phone, At, ...}
UserLoggedIn{UserID, At, ...}  
UserPhoneChanged{UserID, OldPhone, NewPhone, At, ...}
UserNameChanged{UserID, OldName, NewName, At, ...}
UserPasswordChanged{UserID, At, ...}
```

## 📚 Используемые библиотеки

- [Chi Router](https://github.com/go-chi/chi) - HTTP роутер
- [GORM](https://gorm.io/) - ORM для работы с БД
- [bcrypt](https://golang.org/x/crypto/bcrypt) - Хеширование паролей
- [UUID](https://github.com/google/uuid) - Генерация UUID

## 🧪 Тестирование

### 📊 Запуск тестов

```bash
# Unit тесты доменной логики
go test ./tests/domain -v

# Интеграционные тесты с PostgreSQL  
go test -tags=integration ./tests/integration/... -v

# Все тесты сразу
make test-all
```

### 🔧 Требования для интеграционных тестов

```bash
# PostgreSQL через Docker
docker compose up -d postgres

# Создание тестовой БД (автоматически)
CREATE DATABASE quest_auth_test;
```

## 🚀 Развертывание

### 🐳 Docker

```bash
# Сборка образа
docker build -t quest-auth .

# Запуск с БД
docker compose up
```

### ⚙️ Переменные окружения

```bash
HTTP_PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=username
DB_PASSWORD=secret
DB_NAME=quest_auth
DB_SSLMODE=disable
EVENT_GOROUTINE_LIMIT=10
```

## 🔧 Разработка

Проект следует принципам **Clean Architecture** и **Domain-Driven Design**:

- **Domain Layer**: Богатая доменная модель с бизнес-правилами
- **Application Layer**: Use cases, обработчики команд
- **Infrastructure Layer**: Репозитории, внешние адаптеры  
- **Ports & Adapters**: Инверсия зависимостей, тестируемость

### 🎯 Архитектурные решения

- **Безопасность**: bcrypt + валидация на всех уровнях
- **Производительность**: индексы БД + оптимизированные запросы
- **Масштабируемость**: Event Sourcing ready для интеграций
- **Тестируемость**: DI + порты/адаптеры для изоляции
- **Поддерживаемость**: четкое разделение ответственности

## 📈 Мониторинг

### Health Check
```bash
curl http://localhost:8080/health
# Response: {"status":"ok"}
```

### Логирование
- Структурированные логи всех операций
- Отслеживание аутентификации и ошибок
- Domain events для аудита действий пользователей

---

**Quest Auth Service** - надежный и масштабируемый микросервис аутентификации, готовый к production использованию! 🚀