---
title: "/answer"
description: "Методы, реализуемые в рамках сущности answer"
draft: false
version: 1
weight: 18
---

Методы, реализуемые в рамках сущности answer.

## Общее описание

* Архитектура: **REST**
* Протокол взаимодействия: **HTTPS**
* Методы: **GET, PUT, DELETE, POST**
* Формат входных данных: **PATH**
* Формат выходных данных: **JSON (body)**

## Все доступные методы

 * POST   /answers                                    - создание ответа (teacher)
 * GET    /answers/by-question?questionId={id}        - получить все ответы вопроса (student, teacher)
 * DELETE /answers/by-question?questionId={id}        - удалить все ответы вопроса (teacher)
 * GET    /answers/by-question/correct?questionId={id} - получить правильные ответы (student, teacher)
 * GET    /answers/by-question/count?questionId={id}   - подсчет ответов (student, teacher)
 * GET    /answers/by-question/count-correct?questionId={id} - подсчет правильных ответов (student, teacher)
 * 
 * GET    /answers/{id}                               - получить ответ по ID (student, teacher)
 * PUT    /answers/{id}                               - обновить ответ (teacher)
 * DELETE /answers/{id}                               - удалить ответ (teacher)
 * GET    /answers/{id}/correct                       - проверить правильность ответа (student, teacher)