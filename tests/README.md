# Tests Directory

Этот каталог содержит все тесты проекта Quest Auth, организованные по уровням и назначению.

## 📁 Структура тестов

```
tests/
├── domain/                        # 🏗️ Доменные (unit) тесты (Email, Phone, User)
├── contracts/                     # 🤝 Контрактные тесты (UserRepository, UnitOfWork, EventPublisher)
│   └── mocks/                     # Mock реализации для контрактов
└── integration/                   # 🔗 Интеграционные тесты (build tag: integration)
    ├── core/                      # Общие компоненты для интеграционных тестов
    │   ├── assertions/            # Переиспользуемые проверки (HTTP, поля, токены)
    │   ├── case_steps/            # Шаги тестирования (HTTP, handlers, gRPC)
    │   ├── storage/               # Доступ к БД для проверок (events)
    │   └── test_data_generators/  # Генераторы тестовых данных (пользователь)
    └── tests/                     # Группы интеграционных тестов
        ├── auth_http_tests/       # HTTP API (register, login) + валидации уровня HTTP
        ├── auth_handler_tests/    # Handlers (register, login, authenticate) без HTTP
        ├── auth_grpc_tests/       # gRPC Authenticate handler
        ├── auth_e2e_tests/        # E2E: HTTP и токены, bubbling доменных ошибок
        └── repository_tests/      # Репозитории и события (EventRepository, UserRepository)
```

## 🧪 Типы тестов

### 1) Domain Tests (`tests/domain/`)
Unit‑тесты доменной логики и value objects.

```bash
go test ./tests/domain -v
```

### 2) Contract Tests (`tests/contracts/`)
Проверка совместимости реализаций портов (репозитории, UoW, EventPublisher).

```bash
go test ./tests/contracts -v
```

### 3) Integration Tests (`tests/integration/`)
Тесты взаимодействия компонентов с реальной БД и транспортами.

```bash
# Все интеграционные тесты
go test -tags=integration ./tests/integration/... -v

# По группам
go test -tags=integration ./tests/integration/tests/auth_http_tests -v
go test -tags=integration ./tests/integration/tests/auth_handler_tests -v
go test -tags=integration ./tests/integration/tests/auth_grpc_tests -v
go test -tags=integration ./tests/integration/tests/repository_tests -v
go test -tags=integration ./tests/integration/tests/auth_e2e_tests -v
```

## 🚀 Быстрый старт

```bash
# Unit + Contracts
go test ./tests/domain -v && go test ./tests/contracts -v

# Все интеграционные
go test -tags=integration ./tests/integration/... -v
```

## 🔧 Требования для integration

- PostgreSQL (можно через Docker Compose)
- Тестовая БД создаётся автоматически в тестовом контейнере
- Build tag: `-tags=integration`

## 🔁 Переиспользуемые компоненты

- `core/case_steps`: шаги (HTTP, handlers, gRPC)
- `core/assertions`: проверки HTTP, полей ответа, токенов
- `core/storage`: утилиты доступа к событиям в БД (EventStorage)
- `core/test_data_generators`: генерация пользователей и конвертеры запросов

## ✅ Best Practices

- Именование: `TestFunction_Scenario_ExpectedResult`
- Структура: Pre‑condition → Act → Assert
- Изоляция: каждый тест независим
- Данные: используйте `test_data_generators`
- Читаемость: тесты — живая документация
