# Scripts Directory

В этом каталоге находятся вспомогательные скрипты для проекта Quest Auth.

## Доступные скрипты

### 📊 Тестирование и статистика — `test-stats.sh`
Запускает тесты с подробной статистикой.
- Показывает Passed/Failed/Skipped
- Считает процент успеха
- Умеет запускать отдельные группы

```bash
# Примеры
./scripts/test-stats.sh domain        # Доменные тесты
./scripts/test-stats.sh contracts     # Контрактные тесты
./scripts/test-stats.sh integration   # Все интеграционные (с тегом)
./scripts/test-stats.sh http          # HTTP API тесты
./scripts/test-stats.sh handler       # Handler-тесты (без HTTP)
./scripts/test-stats.sh grpc          # gRPC Authenticate
./scripts/test-stats.sh e2e           # E2E тесты
./scripts/test-stats.sh repository    # Репозитории
./scripts/test-stats.sh all           # Все вместе
```

### 📈 Покрытие кода — `coverage-check.sh`
Считается покрытие только для `internal/...`, исключая `tests/`.

```bash
./scripts/coverage-check.sh
```

## Запуск через Makefile

```bash
make test-stats       # Статистика тестов
make coverage-check   # Быстрая проверка покрытия
```

## Права на выполнение

```bash
chmod +x scripts/*.sh
```