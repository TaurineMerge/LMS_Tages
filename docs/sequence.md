sequenceDiagram
    title Прохождение теста и получение сертификата
    participant S as Студент
    participant FE as Фронтенд
    participant NG as Nginx
    participant PA as Personal Account
    participant TS as Testing System
    participant KC as Keycloak
    participant DB as PostgreSQL LMS

    Note over S,DB: 1. Начало теста
    S->>FE: Нажимает "Начать тест"
    FE->>NG: POST /api/tests/{testId}/start
    NG->>TS: Проксирует запрос
    
    Note over TS,KC: 2. Проверка авторизации
    TS->>KC: GET /auth/verify (JWT)
    KC-->>TS: 200 OK (user_id, roles)
    
    Note over TS,DB: 3. Получение вопросов
    TS->>DB: SELECT * FROM test_d WHERE id = ?
    TS->>DB: SELECT * FROM question_d WHERE test_id = ?
    TS->>DB: SELECT * FROM answer_d WHERE question_id IN (?)
    DB-->>TS: Данные теста
    
    TS-->>FE: 200 OK (вопросы теста)
    FE-->>S: Отображает вопросы
    
    Note over S,DB: 4. Ответы на вопросы
    loop Для каждого вопроса
        S->>FE: Выбирает ответ
        FE->>NG: POST /api/tests/{testId}/answer/{questionId}
        NG->>TS: Проксирует запрос
        TS->>DB: INSERT ответ в test_attempt_b
        DB-->>TS: Подтверждение
        TS-->>FE: 200 OK
    end
    
    Note over S,DB: 5. Завершение теста
    S->>FE: Нажимает "Завершить тест"
    FE->>NG: POST /api/tests/{testId}/finish
    NG->>TS: Проксирует запрос
    
    Note over TS,DB: 6. Расчет результатов
    TS->>DB: SELECT score FROM answer_d WHERE id IN (...)
    DB-->>TS: Баллы
    TS->>DB: SELECT min_point FROM test_d WHERE id = ?
    DB-->>TS: Минимальный балл
    TS->>DB: UPDATE test_attempt_b SET point=?, status='completed'
    
    Note over TS,DB: 7. Проверка сертификата
    TS->>DB: SELECT * FROM certificate_b WHERE student_id=? AND test_id=?
    DB-->>TS: Результат
    
    alt Результат >= min_point и нет сертификата
        TS->>DB: INSERT INTO certificate_b (...)
        DB-->>TS: certificate_id
        TS->>DB: UPDATE test_attempt_b SET certificate_id=?
    else Результат < min_point
        Note over TS,PA: Тест не пройден
    end
    
    Note over TS,PA: 8. Сохранение в статистику ЛК
    TS->>PA: POST /api/personal/statistics/test-result
    PA->>KC: GET /auth/verify (JWT)
    KC-->>PA: 200 OK
    
    alt Результат >= min_point
        PA->>DB: Обновляет успешную статистику
        PA->>DB: Записывает в test_history
        TS->>PA: POST /api/personal/notifications/certificate-earned
    else
        PA->>DB: Обновляет неудачную статистику
        PA->>DB: Записывает в test_history
        TS->>PA: POST /api/personal/notifications/test-failed
    end
    
    PA-->>TS: 200 OK
    
    Note over TS,FE: 9. Отправка результатов
    TS->>DB: SELECT детали попытки
    DB-->>TS: Данные
    TS-->>FE: 200 OK (результаты + детали)
    FE-->>S: Отображает результаты
    
    Note over S,PA: 10. Просмотр статистики в ЛК
    S->>FE: Переходит в "Статистика"
    FE->>NG: GET /api/personal/statistics
    NG->>PA: Проксирует
    PA->>KC: Проверка JWT
    KC-->>PA: 200 OK
    PA->>DB: SELECT статистику
    DB-->>PA: Данные
    PA-->>FE: 200 OK
    FE-->>S: Показывает статистику