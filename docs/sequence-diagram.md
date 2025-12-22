---
title: "Sequence Диаграммы LMS"
date: $(date +%Y-%m-%d)
---

## Аутентификация
```mermaid
sequenceDiagram
    participant С as Студент (Гость)
    participant N as nginx
    participant K as Keycloak
    participant PA as Personal Account (Python)
    participant PS as Public Side (Go)
    participant DB_A as PostgreSQL Auth
    participant DB_L as PostgreSQL LMS
    
    Note over С,DB_L: Шаг 1: Регистрация через Keycloak
    С->>N: POST /auth/register
    N->>K: Проксирует запрос
    K->>DB_A: Сохраняет нового пользователя
    DB_A-->>K: Успех
    K-->>N: 201 Created + user_id
    N-->>С: Успешная регистрация
    
    Note over С,DB_L: Шаг 2: Аутентификация
    С->>N: POST /auth/login (email/password)
    N->>K: Проксирует запрос
    K->>DB_A: Проверяет учетные данные
    DB_A-->>K: Валидны
    K->>DB_A: Создает сессию и токены
    DB_A-->>K: Успех
    K-->>N: JWT Access + Refresh токены
    N-->>С: Токены
    
    Note over С,DB_L: Шаг 3: Получение профиля
    С->>N: GET /api/personal/profile (с JWT)
    N->>PA: Проксирует запрос
    PA->>K: Валидирует JWT /auth/verify
    K-->>PA: {valid: true, user_id: X, roles: ["student"]}
    PA->>DB_L: Ищет student_s по user_id
    DB_L-->>PA: Не найден (первый вход)
    PA->>DB_L: Создает запись student_s
    DB_L-->>PA: Успех (uuid студента)
    PA-->>N: Профиль студента с avatar=null
    N-->>С: Профиль создан
    
    Note over С,DB_L: Шаг 4: Просмотр доступных курсов
    С->>N: GET /api/public/courses?visibility=public
    N->>PS: Проксирует запрос
    PS->>DB_L: SELECT из course_b WHERE visibility='public'
    DB_L-->>PS: Список курсов
    PS-->>N: JSON с курсами
    N-->>С: Каталог курсов

```

---

## Прохождение теста и получение сертификата

```mermaid
sequenceDiagram
    participant Student as Студент
    participant PA as Personal Account
    participant TS as Testing System
    participant DB as DB (testing_schema)

    Note over Student, DB: Предварительно: Студент аутентифицирован через Keycloak, имеет валидный JWT токен
    
    Student->>PA: POST /api/personal/tests/{id}/start (с JWT в заголовке)
    Note right of PA: PA уже имеет user_id из JWT<br/>или проверяет через Keycloak
    
    PA->>TS: Перенаправление запроса с user_id и test_id
    TS->>DB: Создание test_attempt_b (student_id, test_id, status='in_progress')
    DB-->>TS: attempt_id
    TS-->>Student: Тест начат, вопросы
    
    loop Ответы на вопросы
        Student->>TS: POST /api/tests/attempt/{id}/answer (с JWT)
        TS->>DB: Сохранение ответа в attempt_version (JSON)
    end

    Student->>TS: POST /api/tests/attempt/{id}/finish
    TS->>DB: Расчет балла, проверка min_point
    DB-->>TS: Результат проверки (points, min_point)
    
    alt Балл >= min_point
        TS->>DB: Создание certificate_b
        DB-->>TS: certificate_id
        TS->>DB: Обновление test_attempt_b (status='passed', point, certificate_id)
        TS-->>Student: Успех + данные сертификата
    else Балл < min_point
        TS->>DB: Обновление test_attempt_b (status='failed', point)
        TS-->>Student: Тест не пройден + рекомендации для повторения
    end
