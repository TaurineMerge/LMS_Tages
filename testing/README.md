# Testing Service (Javalin)

Сервис для управления тестами, попытками прохождения тестов и ответами в системе LMS.

## API Документация

Сервис предоставляет статическую Swagger-документацию.

- **Swagger JSON**: `http://localhost:{APP_PORT}/swagger.json`

Документация описывает все эндпоинты API с примерами запросов и ответов.

## Эндпоинты

### Тесты (/tests)
- `GET /tests` - Получить список тестов (только преподаватели)
- `POST /tests` - Создать тест (только преподаватели)
- `GET /tests/{id}` - Получить тест по ID (все роли)
- `PUT /tests/{id}` - Обновить тест (только преподаватели)
- `DELETE /tests/{id}` - Удалить тест (только преподаватели)

### Попытки тестов (/test-attempts)
- `GET /test-attempts` - Получить список попыток
- `POST /test-attempts` - Создать попытку
- `GET /test-attempts/{id}` - Получить попытку по ID
- `DELETE /test-attempts/{id}` - Удалить попытку

### Ответы (/answers)
- `POST /answers` - Создать ответ
- `GET /answers/by-question?questionId={id}` - Получить ответы вопроса
- `DELETE /answers/by-question?questionId={id}` - Удалить ответы вопроса
- `GET /answers/by-question/correct?questionId={id}` - Получить правильные ответы
- `GET /answers/by-question/count?questionId={id}` - Подсчет ответов
- `GET /answers/by-question/count-correct?questionId={id}` - Подсчет правильных ответов
- `GET /answers/{id}` - Получить ответ по ID
- `PUT /answers/{id}` - Обновить ответ
- `DELETE /answers/{id}` - Удалить ответ
- `GET /answers/{id}/correct` - Проверить правильность ответа

## Deploy
```bash
docker build -t javalin-app . && docker run javalin-app
