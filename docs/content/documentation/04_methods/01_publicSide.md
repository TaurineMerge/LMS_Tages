---
title: "Методы для публичной страницы"
description: "Методы, реализуемые в рамках публичной страницы"
draft: false
version: 1
weight: 18
---

## Общее описание

* Методы: **GET, PUT, DELETE, POST**
* Формат входных данных: **PATH**

## Все доступные методы аутентификации

* GET /login                                                        - вход в аккаунт
* GET /logout                                                       - выход из аккаунта
* GET /auth/callback                                                - авторизация через keycloak
* * GET /reg                                                        - регистрация

## Внешние сервисы

* GET /account/profile                                              - получение профиля из ЛК

## Все доступные маршруты веб-приложения

* GET /categories                                                   - получение всех категорий
* GET /categories/category_id                                       - получение категории по ID
* GET /categories/category_id/courses                               - получение всех курсов
* GET /categories/category_id/courses/course_id                     - получение курса по ID
* GET /categories/category_id/courses/course_id/lessons             - получение всех уроков
* GET /categories/category_id/courses/course_id/lessons/lesson_id   - получение урока по ID

### Быстрый доступ
[← На главную](/documentation/04_methods/)