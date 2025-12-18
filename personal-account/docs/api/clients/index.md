# Clients

Модуль клиентов для интеграции с внешними сервисами и валидации контрактов.

## Обзор

Пакет `app.clients` содержит клиенты для взаимодействия с внешними сервисами
и инструменты для валидации API-контрактов.

## Модули

- **[Testing](testing.md)** — клиент для работы с сервисом тестирования,
 - **[Testing Client](testing.md)** — клиент для работы с сервисом тестирования.
 - **[Validation tools](validation/index.md)** — общие утилиты валидации и загрузки схем
     (SchemaLoader и ContractManager), вынесённые в `app.clients.validation`.

## Архитектура

```
app/clients/
├── __init__.py
├── validation/
│   ├── __init__.py
│   ├── contract_manager.py   # Менеджер валидации контрактов
│   ├── schema_loader.py      # Загрузчик JSON-схем
│   └── schemas/              # JSON-схемы контрактов
│       ├── user_stats/
│       ├── attempt_detail/
│       └── attempts_list/
└── testing_client.py         # HTTP client for testing service
```

## Быстрый старт

```python
from app.clients.testing import ContractManager, ContractValidationError

# Инициализация менеджера
manager = ContractManager()

# Валидация данных
try:
    validated = await manager.validate_user_stats({
        "student_id": "550e8400-e29b-41d4-a716-446655440000",
        "total_attempts": 5,
        "passed_attempts": 3,
        # ...
    })
    print("✓ Данные валидны")
except ContractValidationError as e:
    print(f"✗ Ошибка: {e.message}")
```
