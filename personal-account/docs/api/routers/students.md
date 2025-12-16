# Students Router

API endpoints для управления студентами.

## Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/v1/students` | Список студентов (пагинация) |
| GET | `/api/v1/students/{id}` | Получить студента по ID |
| POST | `/api/v1/students` | Создать студента |
| PUT | `/api/v1/students/{id}` | Обновить студента |
| DELETE | `/api/v1/students/{id}` | Удалить студента |

## API Reference

::: app.routers.students
    options:
      show_root_heading: false
      members_order: source
