---
title: "Lesson"
description: "Описание словаря 'lesson_d'"
version: 1
draft: false
weight: 3
    
---

## Словарь, содержащий список уроков

| Attribute   | Type         | DB column  | MUST               | Not Null           | Comments                                   |
|:------------|:-------------|:-----------|:------------------:|:------------------:|:-------------------------------------------|
| id          | uuid         | id         | :white_check_mark: | :white_check_mark: | Уникальный идентификатор урока             |
| title       | varchar(255) | title      | :white_check_mark: | :white_check_mark: | Название урока                             |
| course_id   | uuid         | course_id  | :white_check_mark: | :white_check_mark: | Внешний ключ на курс                       |
| visibility  | varchar(20)  | visibility | :white_check_mark: | :white_check_mark: | Видимость урока: draft, public             |
| content     | text         | content    | :white_check_mark: | :white_check_mark: | Содержимое урока (текст, HTML и т.д.)      |
| created_at  | timestamp    | created_at | :white_check_mark: | :white_check_mark: | Дата и время создания                      |
| updated_at  | timestamp    | updated_at | :white_check_mark: | :white_check_mark: | Дата и время последнего обновления         |

## Быстрый доступ
[← На главную](/documentation/02_info_model/)