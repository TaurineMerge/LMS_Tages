---
title: "Flowchart-диаграмма архитектуры LMS"
date: 2025-12-15
layout: "single"
---

{{< mermaid >}}

flowchart TD
    A[Пользователь заходит на сайт] --> B{Тип пользователя?}
    B -->|Зарегистрированный ученик| D[Вход в систему]
    
    D --> M[Keycloak: аутентификация]
    M --> N[JWT токен]
    N --> O[Доступ к личному кабинету]
    O --> STUDENT_ACTION{Действие?}
    
    STUDENT_ACTION -->|Управление профилем| STUDENT_PROFILE_ACTION{Управление профилем?}
    STUDENT_ACTION -->|Пройти уроки| STUDENT_CAT_SELECT[Выбор категории]
    STUDENT_ACTION -->|Пройти тест| R[Тестирование]
    STUDENT_ACTION -->|Выйти| STUDENT_LOGOUT[Выход из личного кабинета]
    
    STUDENT_PROFILE_ACTION -->|Редактировать профиль| STUDENT_EDIT{Что редактировать?}
    STUDENT_PROFILE_ACTION -->|Просмотреть информацию| STUDENT_VIEW{Что просмотреть?}
    STUDENT_PROFILE_ACTION -->|Удалить профиль| STUDENT_DELETE[Запрос на удаление]
    
    STUDENT_DELETE --> STUDENT_DELETE_CONFIRM{Подтверждение удаления?}
    STUDENT_DELETE_CONFIRM -->|Да| STUDENT_DELETE_KEYCLOAK[Удаление в Keycloak]
    STUDENT_DELETE_CONFIRM -->|Нет| STUDENT_PROFILE_ACTION
    
    STUDENT_DELETE_KEYCLOAK --> STUDENT_DELETE_PERSONAL[Удаление в Personal Account]
    
    STUDENT_EDIT -->|ФИО| STUDENT_EDIT_NAME[Редактирование ФИО]
    STUDENT_EDIT -->|Пароль| STUDENT_EDIT_PASS[Смена пароля]
    STUDENT_EDIT -->|Почту| STUDENT_EDIT_EMAIL[Редактирование email]
    STUDENT_EDIT -->|Номер телефона| STUDENT_EDIT_PHONE[Редактирование телефона]
    STUDENT_EDIT -->|Дата рождения| STUDENT_EDIT_BIRTHDATE[Редактирование даты рождения]
    STUDENT_EDIT -->|Аватар| STUDENT_EDIT_AVATAR[Загрузка аватара]
    
    STUDENT_EDIT_NAME --> STUDENT_UPDATE_NAME[Обновление ФИО в Personal Account]
    STUDENT_EDIT_PASS --> STUDENT_UPDATE_PASS[Обновление пароля в Keycloak]
    STUDENT_EDIT_EMAIL --> STUDENT_UPDATE_EMAIL[Обновление email в Personal Account и Keycloak]
    STUDENT_EDIT_PHONE --> STUDENT_UPDATE_PHONE[Обновление телефона в Personal Account]
    STUDENT_EDIT_BIRTHDATE --> STUDENT_UPDATE_BIRTHDATE[Обновление даты рождения в Personal Account]
    STUDENT_EDIT_AVATAR --> STUDENT_UPDATE_AVATAR[Загрузка и обновление аватара в Personal Account]
    
    STUDENT_UPDATE_NAME --> STUDENT_SAVE_PROFILE[Сохранение изменений профиля]
    STUDENT_UPDATE_PASS --> STUDENT_SAVE_PROFILE
    STUDENT_UPDATE_EMAIL --> STUDENT_SAVE_PROFILE
    STUDENT_UPDATE_PHONE --> STUDENT_SAVE_PROFILE
    STUDENT_UPDATE_BIRTHDATE --> STUDENT_SAVE_PROFILE
    STUDENT_UPDATE_AVATAR --> STUDENT_SAVE_PROFILE
    
    STUDENT_VIEW -->|Статистика| STUDENT_VIEW_STATS{Какую статистику?}
    STUDENT_VIEW -->|Сертификаты| STUDENT_VIEW_CERTS[Просмотр сертификатов]
    STUDENT_VIEW -->|Профиль| STUDENT_VIEW_PROFILE[Просмотр профиля]
    
    STUDENT_VIEW_STATS -->|Курсы| STUDENT_STATS_COURSES[Статистика по курсам]
    STUDENT_VIEW_STATS -->|Тесты| STUDENT_STATS_TESTS[Статистика по тестам]
    
    STUDENT_STATS_COURSES --> STUDENT_SHOW_STATS[Отображение статистики]
    STUDENT_STATS_TESTS --> STUDENT_SHOW_STATS
    
    STUDENT_VIEW_CERTS --> STUDENT_SHOW_CERTS[Отображение сертификатов]
    
    STUDENT_VIEW_PROFILE --> STUDENT_SHOW_PROFILE[Отображение данных профиля]
    
    STUDENT_CAT_SELECT --> STUDENT_CAT_VIEW[Просмотр курсов в категории]
    STUDENT_CAT_VIEW --> STUDENT_COURSE_SELECT[Выбор курса]
    STUDENT_COURSE_SELECT --> STUDENT_LESSON_SELECT[Выбор урока]
    STUDENT_LESSON_SELECT --> Q[Прохождение урока]
    
    Q --> T[Отметка пройденных уроков]
    T --> U[Обновление прогресса в Personal Account]
    
    R --> STUDENT_TEST_CAT_SELECT[Выбор категории для теста]
    STUDENT_TEST_CAT_SELECT --> STUDENT_TEST_COURSE_SELECT[Выбор курса]
    STUDENT_TEST_COURSE_SELECT --> STUDENT_TEST_SELECT[Выбор теста]
    STUDENT_TEST_SELECT --> W[Прохождение вопросов]
    
    W --> X[Проверка ответов Testing System]
    X --> Y{Результат?}
    Y -->|Успешно| S[Генерация сертификата]
    Y -->|Неудача| STUDENT_ACTION

{{< /mermaid >}}

## Навигация
- [Документация Account →](/account/documentation)
- [← Назад к Account](/account/)