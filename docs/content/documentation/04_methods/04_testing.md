---
title: "Методы для системы тестирования"
description: "Методы, реализуемые в рамках веб-интерфейса (HTML), UI и internal (для общения с другими модулями)"
draft: false
version: 1
weight: 18
---

## Общее описание

* Методы: **GET, PUT, DELETE, POST**
* Формат входных данных: **PATH**

## Все доступные методы в рамках веб-интерфейса (HTML-форма)

 * GET      /testing/web/tests/new                      - создание нового теста
 * POST     /testing/web/tests/save                     - создание нового теста (через универсальную форму)
 * POST     /testing/web/tests/draft                    - создание нового черновика
 * GET      /testing/web/tests/draft-{id}/edit          - форма редактирования черновика
 * GET      /testing/web/tests/{id}/edit                - форма редактирования теста
 * GET      /testing/web/tests/{id}/preview             - ппредпросмотр теста
 * PUT      /testing/web/tests/draft-{id}               - полное обновление существующего черновика
 * PUT      /testing/web/tests/{id}                     - полное обновление существующего теста
 * PATCH    /testing/web/tests/draft-{id}               - частичное обновление черновика
 * POST     /testing/web/tests/{id}/delete              - удаление теста
 * DELETE   /testing/web/tests/draft-{id}               - удаление черновика
 * POST     /testing/web/tests/create-draft             - создание черновика из существующего теста
 * POST     /testing/web/tests/publish                  - публикация черновика

---

## Все доступные методы для UI


 * GET    /ui/tests/{testId}/take                               - получить тест для прохождения
 * POST   /ui/tests/{testId}/questions/{questionId}/answer      - отправить ответ на конкретный вопрос
 * GET    /ui/tests/{testId}/finish                             - 
 * POST   /ui/tests/{testId}/finish                             - 
 * GET    /ui/tests/{testId}/results                            - получить результаты теста
 * POST   /ui/tests/{testId}/retry                              - пройти "еще раз" - создает новую попытку и перекидывает на метод /take

---

## Все доступные методы для internal (взаимодействие с другими сервисами)

 * GET       testing/internal/health                                                  — публичные маршруты (без авторизации)
 * GET       testing/internal/categories/{categoryId}/courses/{courseId}/test         — получить тест, который привязан к курсу
 * GET       testing/internal/users/{userId}/attempts                                 — получить попытки тестов конкретного пользователя
 * GET       testing/internal/users/{userId}/stats                                    — получить статистику пользователя
 * GET       testing/internal/attempts/{attemptId}                                    — получить конкретную попытку прохождения теста
 * GET       testing/internal/categories/{categoryId}/courses/{courseId}/draft        — получить черновик теста, привязанны к курсу

### Быстрый доступ
[← На главную](/documentation/04_methods/)