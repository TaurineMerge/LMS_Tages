---
title: "Test"
description: "Описание словаря 'test_d"
version: 1
draft: false
weight: 3
    
---

## Словарь, содержащий тест целиком

| Attribute   | Type    | DB column  | MUST               | Not Null           | Comments                                   |
|:------------|:--------|:-----------|:------------------:|:------------------:|:-------------------------------------------|
| id          | uuid    | id         | :white_check_mark: | :white_check_mark: | Уникальный идентификатор теста             |
| course_id   | uuid    | course_id  | :white_check_mark: | :white_check_mark: | Внешний ключ на курс                       |
| title       | varchar | title      | :x:                | :x:                | Название теста                             |
| min_point   | integer | min_point  | :x:                | :x:                | Минимальное количество баллов для прохождения |
| description | text    | description| :x:                | :x:                | Описание теста                             |

## Быстрый доступ
[← На главную](/documentation/02_info_model/)