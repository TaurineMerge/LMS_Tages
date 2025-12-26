---
title: "Draft"
description: "Описание сущности 'draft_b'"
version: 1
draft: false
weight: 3
    
---

## Сущность черновик теста

| Attribute   | Type    | DB column  | MUST               | Not Null           | Comments                                   |
|:------------|:--------|:-----------|:------------------:|:------------------:|:-------------------------------------------|
| id          | uuid    | id         | :white_check_mark: | :white_check_mark: | Уникальный идентификатор черновика         |
| title       | varchar | title      | :x:                | :x:                | Название черновика теста                   |
| min_point   | integer | min_point  | :x:                | :x:                | Минимальное количество баллов для прохождения |
| description | text    | description| :x:                | :x:                | Описание черновика теста                   |
| test_id     | uuid    | test_id    | :x:                | :x:                | Внешний ключ на опубликованный тест        |

## Быстрый доступ
[← На главную](/documentation/02_info_model/)