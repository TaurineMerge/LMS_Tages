---
title: "Content"
description: "Описание словаря 'content_d'"
version: 1
draft: false
weight: 3
    
---

## Словарь, содержащий контент, необходимый для теста

| Attribute       | Type    | DB column       | MUST               | Not Null           | Comments                                   |
|:----------------|:--------|:----------------|:------------------:|:------------------:|:-------------------------------------------|
| id              | uuid    | id              | :white_check_mark: | :white_check_mark: | Уникальный идентификатор контента          |
| order           | int     | order           | :x:                | :x:                | Порядок контента в рамках вопроса/ответа   |
| content         | text    | content         | :x:                | :x:                | Содержимое (ссылка на медиа-файл или текст) |
| type_of_content | boolean | type_of_content | :x:                | :x:                | Тип контента (0 - для вопроса, 1 - для ответа) |
| question_id     | uuid    | question_id     | :white_check_mark: | :white_check_mark: | Внешний ключ на вопрос                      |
| answer_id       | uuid    | answer_id       | :white_check_mark: | :white_check_mark: | Внешний ключ на ответ                       |

## Быстрый доступ
[← На главную](/documentation/02_info_model/)