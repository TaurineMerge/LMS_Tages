# Testing Client

Клиент для работы с сервисом тестирования и валидации API-контрактов.

## Обзор

Модуль `app.clients.testing` предоставляет высокопроизводительные инструменты
для валидации API-контрактов с использованием JSON Schema и библиотеки
[fastjsonschema](https://github.com/horejsek/python-fastjsonschema).

### Ключевые возможности

- ⚡ **Высокая производительность** — скомпилированные валидаторы fastjsonschema
  (в 10-100 раз быстрее стандартного jsonschema)
- 🎯 **Версионирование схем** — поддержка v1, v2 и автоматический выбор "latest"
- 💾 **Умное кэширование** — скомпилированные валидаторы кэшируются
- 📊 **Детальные ошибки** — структурированные объекты ошибок с JSON-путями
- 🔄 **Async First** — нативная поддержка async/await
- 🔒 **Type Safe** — полная типизация для IDE

---

## ContractManager

::: app.clients.testing.contract_manager.ContractManager
    options:
      show_root_heading: true
      show_source: true
      members_order: source

---

## ContractValidationError

::: app.clients.testing.contract_manager.ContractValidationError
    options:
      show_root_heading: true
      show_source: true

---

## SchemaLoader

::: app.clients.testing.schema_loader.SchemaLoader
    options:
      show_root_heading: true
      show_source: true
      members_order: source

---

## Примеры использования

### Базовая валидация

```python
from app.clients.testing import ContractManager, ContractValidationError

async def validate_stats():
    manager = ContractManager()
    
    data = {
        "student_id": "550e8400-e29b-41d4-a716-446655440000",
        "total_attempts": 5,
        "passed_attempts": 3,
        "failed_attempts": 2,
        "average_score": 78.5,
        "best_score": 95.0,
        "total_time_spent": 3600,
        "last_attempt_at": "2024-01-15T10:30:00Z"
    }
    
    try:
        validated = await manager.validate_user_stats(data)
        print("✓ Данные валидны!")
        return validated
    except ContractValidationError as e:
        print(f"✗ Ошибка валидации: {e.message}")
        raise
```

### Валидация с указанием версии

```python
async def validate_with_version():
    manager = ContractManager()
    
    # Явное указание версии схемы
    validated = await manager.validate(
        data={"student_id": "...", ...},
        contract_name="user_stats",
        version="v1"
    )
    
    # Использование "latest" (по умолчанию)
    validated = await manager.validate(
        data={"student_id": "...", ...},
        contract_name="user_stats"
    )
```

### Обработка ошибок

```python
async def handle_validation_errors():
    manager = ContractManager()
    
    invalid_data = {
        "student_id": "not-a-uuid",  # Невалидный UUID
        "total_attempts": -1,         # Отрицательное число
    }
    
    try:
        await manager.validate_user_stats(invalid_data)
    except ContractValidationError as e:
        print(f"Сообщение: {e.message}")
        # Вывод:
        # Contract 'user_stats' validation failed:
        #   1. [student_id] 'not-a-uuid' does not match pattern...
        #   2. [total_attempts] -1 is less than minimum 0
```

### Использование в FastAPI endpoint

```python
from fastapi import APIRouter, HTTPException
from app.clients.testing import ContractManager, ContractValidationError

router = APIRouter()
contract_manager = ContractManager()

@router.post("/stats/validate")
async def validate_statistics(data: dict):
    """Валидация статистики пользователя."""
    try:
        validated = await contract_manager.validate_user_stats(data)
        return {"status": "valid", "data": validated}
    except ContractValidationError as e:
        raise HTTPException(
            status_code=422,
            detail={"message": e.message}
        )
```

---

## JSON-схемы

Схемы хранятся в директории `app/clients/testing/schemas/` и организованы
по контрактам и версиям:

```
schemas/
├── user_stats/
│   └── v1.json
├── attempt_detail/
│   └── v1.json
└── attempts_list/
    └── v1.json
```

### Создание новой схемы

1. Создайте директорию для контракта:
   ```bash
   mkdir -p app/clients/testing/schemas/my_contract
   ```

2. Создайте файл схемы `v1.json`:
   ```json
   {
     "$schema": "https://json-schema.org/draft/2020-12/schema",
     "$id": "my_contract_v1",
     "title": "MyContract",
     "type": "object",
     "required": ["field1", "field2"],
     "properties": {
       "field1": {"type": "string"},
       "field2": {"type": "integer", "minimum": 0}
     }
   }
   ```

3. Используйте в коде:
   ```python
   validated = await manager.validate(data, "my_contract", "v1")
   ```

---

## Производительность

fastjsonschema компилирует JSON-схемы в оптимизированный Python-код,
что обеспечивает значительный прирост производительности:

| Библиотека     | Время валидации | Относительно |
|----------------|-----------------|--------------|
| jsonschema     | 100 мс          | 1x           |
| fastjsonschema | 1-10 мс         | 10-100x      |

Валидаторы кэшируются после первой компиляции, поэтому повторные
валидации выполняются мгновенно.
