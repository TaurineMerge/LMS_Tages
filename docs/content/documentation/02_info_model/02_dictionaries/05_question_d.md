---
title: "Question"
description: "Описание словаря 'question_d'"
version: 1
draft: false
weight: 3
    
---

## Словарь, содержащий список вопросов

| Attribute        | Type     | DB column       | MUST               | Not Null           | Comments                                    |
|:-----------------|:--------:|:----------------|:------------------:|:------------------:|:--------------------------------------------|
| id               | uuid     | id              | :white_check_mark: | :white_check_mark: | Уникальный идентификатор вопроса            |
| test_id          | uuid     | test_id         | :x:                | :x:                | Внешний ключ на тест (может быть NULL)      |
| draft_id         | uuid     | draft_id        | :x:                | :x:                | Внешний ключ на черновик (может быть NULL)  |
| text_of_question | text     | text_of_question| :x:                | :x:                | Текст вопроса                               |
| order            | int      | order           | :x:                | :x:                | Порядок вопроса в тесте                     |

## Быстрый доступ
[← На главную](/documentation/02_info_model/)