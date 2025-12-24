---
title: "Sequence-диаграммы архитектуры LMS"
date: 2025-12-15
layout: "single"
---

## Авторизация и аутентификация пользователя

{{< mermaid >}}

sequenceDiagram
    participant С as Студент (Гость)
    participant N as nginx
    participant K as Keycloak
    participant PA as Personal Account (Python)
    participant PS as Public Side (Go)
    participant DB_A as KC db
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
    С->>N: GET localhost/categories/category_{id}/courses/
    N->>PS: Проксирует запрос
    PS->>DB_L: SELECT из course_b WHERE visibility='public'
    DB_L-->>PS: Список курсов
    PS-->>N: JSON с курсами
    N-->>С: Каталог курсов

{{< /mermaid >}}