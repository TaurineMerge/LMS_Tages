---
title: "Category"
description: "Описание словаря 'category_d"
version: 1
draft: false
weight: 3
    
---

## Словарь, содержащий список категорий

| Attribute   | Type         | DB column  | MUST               | Not Null           | Comments                                   |
|:------------|:-------------|:-----------|:------------------:|:------------------:|:-------------------------------------------|
| id          | uuid         | id         | :white_check_mark: | :white_check_mark: | Уникальный идентификатор категории         |
| title       | varchar(255) | title      | :white_check_mark: | :white_check_mark: | Название категории                         |
| created_at  | timestamp    | created_at | :white_check_mark: | :white_check_mark: | Дата и время создания                      |
| updated_at  | timestamp    | updated_at | :white_check_mark: | :white_check_mark: | Дата и время последнего обновления         |


## Быстрый доступ
[← На главную](/documentation/02_info_model/)