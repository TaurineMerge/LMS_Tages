---
title: "Sequence-диаграммы админ-панели LMS"
date: 2025-12-15
layout: "single"
---

## Вход в админ-панель и создание категории

{{< mermaid >}}

sequenceDiagram
    participant Admin as Админ (Teacher)
    participant N as nginx
    participant AP as Admin Panel (Go)
    participant K as Keycloak
    participant KB as Knowledge Base DB
    
    Note over Admin,KB: Шаг 1: Вход в админ-панель
    Admin->>N: GET localhost/admin
    N->>AP: Проксирует запрос
    AP->>K: Проверка JWT токена и роли
    K-->>AP: {valid: true, user_id: X, roles: ["teacher"]}
    AP-->>N: HTML страницы админ-панели
    N-->>Admin: Отображает админ-панель
    
    Note over Admin,KB: Шаг 2: Создание категории
    Admin->>N: POST /admin/categories (title: "Новая категория")
    N->>AP: Проксирует запрос
    AP->>K: Проверка JWT токена и роли
    K-->>AP: {valid: true, user_id: X, roles: ["teacher"]}
    AP->>KB: INSERT INTO category_d (title) VALUES ("Новая категория")
    KB-->>AP: Успех (id категории)
    AP-->>N: 201 Created + id категории
    N-->>Admin: Категория успешно создана

{{< /mermaid >}}

---

## Создание курса с загрузкой изображения

{{< mermaid >}}

sequenceDiagram
    participant Admin as Админ (Teacher)
    participant N as nginx
    participant AP as Admin Panel (Go)
    participant K as Keycloak
    participant M as MinIO (S3)
    participant KB as Knowledge Base DB
    
    Note over Admin,KB: Шаг 1: Загрузка формы создания курса
    Admin->>N: GET /admin/categories/category_{id}/courses/new
    N->>AP: Проксирует запрос
    AP->>K: Проверка JWT токена и роли
    K-->>AP: {valid: true, user_id: X, roles: ["teacher"]}
    AP-->>N: HTML формы создания курса
    N-->>Admin: Отображает форму
    
    Note over Admin,KB: Шаг 2: Отправка данных курса с изображением
    Admin->>N: POST /admin/categories/category_{id}/courses (multipart/form-data:<br/>title, description, level, category_id, visibility, image_file)
    N->>AP: Проксирует запрос с файлом
    
    Note over AP,K: Проверка прав доступа
    AP->>K: Проверка JWT токена и роли
    K-->>AP: {valid: true, user_id: X, roles: ["teacher"]}
    
    Note over AP,M: Загрузка изображения в S3
    AP->>AP: Парсинг multipart формы<br/>Извлечение image_file
    AP->>M: PUT /go/{uuid}.jpg<br/>Загрузка файла
    M-->>AP: Успех, возвращает image_key
    
    Note over AP,KB: Сохранение курса в БД
    AP->>KB: INSERT INTO course_b<br/>(title, description, level, category_id, visibility, image_key)<br/>VALUES (..., "images/{uuid}.jpg")
    KB-->>AP: Успех (id курса)
    AP-->>N: 201 Created + id курса
    N-->>Admin: Курс успешно создан

{{< /mermaid >}}

---

## Создание и редактирование урока

{{< mermaid >}}

sequenceDiagram
    participant Admin as Админ (Teacher)
    participant N as nginx
    participant AP as Admin Panel (Go)
    participant K as Keycloak
    participant KB as Knowledge Base DB
    
    Note over Admin,KB: Вариант A: Создание урока
    Admin->>N: POST /admin/categories/category_{id}/courses/course_{id}/lessons<br/>(course_id, title, description, content, visibility)
    N->>AP: Проксирует запрос
    AP->>K: Проверка JWT токена и роли
    K-->>AP: {valid: true, user_id: X, roles: ["teacher"]}
    AP->>KB: INSERT INTO lesson_d<br/>(course_id, title, description, content, visibility)<br/>VALUES (...)
    KB-->>AP: Успех (id урока)
    AP-->>N: 201 Created
    N-->>Admin: Урок создан
    
    Note over Admin,KB: Вариант B: Редактирование урока
    Admin->>N: PUT /admin/categories/category_{id}/courses/course_{id}/lessons/lessons/{id}<br/>(title, description, content, visibility)
    N->>AP: Проксирует запрос
    AP->>K: Проверка JWT токена и роли
    K-->>AP: {valid: true, user_id: X, roles: ["teacher"]}
    AP->>KB: UPDATE lesson_d SET<br/>title=..., description=..., content=..., visibility=...<br/>WHERE id = {lesson_id}
    KB-->>AP: Успех
    AP-->>N: 200 OK
    N-->>Admin: Урок обновлен

{{< /mermaid >}}

---

## Кнопка "Создать тест"

{{< mermaid >}}

sequenceDiagram
    participant Admin as Админ (Teacher)
    participant N as nginx
    participant AP as Admin Panel (Go)
    participant TB as Test Builder Service
    participant KB as Knowledge Base DB
    
    Note over Admin,KB: Шаг 1: Нажатие кнопки в админке
    Admin->>N: GET /admin/courses/course_{id}/create-test (пример)
    N->>AP: Проксирует запрос
    AP->>KB: Получение данных курса
    KB-->>AP: Данные курса
    AP->>AP: Генерация ссылки на конструктор тестов<br/>с параметрами (course_id, title)
    
    Note over AP,TB: Перенаправление на внешний сервис
    AP-->>N: 302 Redirect → https://test-builder.com/create?course_id=...
    N-->>Admin: Перенаправление на конструктор тестов
    
    Note over Admin,TB: Работа во внешнем сервисе
    Admin->>TB: Работает в конструкторе тестов
    TB-->>Admin: Тест создан/сохранен
    
    Note over Admin,AP: Возврат в админ-панель
    Admin->>N: GET /admin/courses/course_{id}
    N->>AP: Проксирует запрос
    AP-->>N: Страница курса
    N-->>Admin: Отображает страницу курса<br/>(теперь с кнопкой "Редактировать тест")

{{< /mermaid >}}