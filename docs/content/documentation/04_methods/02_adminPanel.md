---
title: "Методы для системы администрирования"
description: "Методы, реализуемые в рамках админ панели"
draft: false
version: 1
weight: 18
---

## Общее описание

* Методы: **GET, PUT, DELETE, POST**
* Формат входных данных: **PATH**

## Все доступные методы для управления категориями курсов

* GET admin/categories                                                       - получение списка категорий
* POST admin/categories                                                      - создание новой категории на основе JSON в теле запроса
* GET admin/categories/:category_id                                          - получение категории по ID
* PUT admin/categories/:category_id                                          - обновление категории по ID на основе JSON в теле запроса
* DELETE admin/categories/:category_id                                       - удаление категории по ID

## Все доступные методы для управления курсами

* Get  admin/categories/:category_id/courses                                 - возвращает список курсов для категорий с фильтрами и пагинацией
* Post admin/categories/:category_id/courses/create                          - создает новый курс для категории на основе JSON  в теле запроса
* Get  admin/categories/:category_id/courses/:course_id                      - возвращает курс по ID в категории
* Post admin/categories/:category_id/courses/:course_id/update               - обновляет курс по ID  в категории на основе JSON  в теле запроса
* Post admin/categories/:category_id/courses/:course_id/delete               - удаляет курс по ID  в категории на основе JSON  в теле запроса

## Все доступные методы для управления уроками

* Get    admin/categories/:category_id/courses/:course_id/lessons             - возвращает список уроков для курса с пагинацией и сортировкой
* Post   admin/categories/:category_id/courses/:course_id/lessons             - создает новый урок на основе данных из тела запроса
* Get    admin/categories/:category_id/courses/:course_id/lessons/:lesson_id  - возвращает урок по его ID
* Put    admin/categories/:category_id/courses/:course_id/lessons/:lesson_id  - обновляет существующий урок по его ID на основе данных из тела запроса
* Delete admin/categories/:category_id/courses/:course_id/lessons/:lesson_id  - удаляет урок по его ID

## Все доступные методы для загрузки изображений

* POST /upload/image                                           - загрузка изображения из multipart формы в S#-совместимое хранилище
* POST /upload/image-from-url                                  - загрузка изображения по указанному URL в S3-совместимое хранилище                     

### Быстрый доступ
[← На главную](/documentation/04_methods/)