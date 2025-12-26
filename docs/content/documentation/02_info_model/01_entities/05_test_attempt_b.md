---
title: "Test"
description: "Описание сущности 'test_attempt_b'"
version: 1
draft: false
weight: 3
    
---

## Сущность попытка теста

| Attribute        | Type     | DB column       | MUST               | Not Null           | Comments                                   |
|:-----------------|:--------:|:----------------|:------------------:|:------------------:|:-------------------------------------------|
| id               | uuid     | id              | :white_check_mark: | :white_check_mark: | Уникальный идентификатор попытки (составной ключ) |
| student_id       | uuid     | student_id      | :white_check_mark: | :white_check_mark: | Внешний ключ на студента                   |
| test_id          | uuid     | test_id         | :white_check_mark: | :white_check_mark: | Внешний ключ на тест                       |
| date_of_attempt  | date     | date_of_attempt | :x:                | :x:                | Дата попытки                               |
| point            | int      | point           | :x:                | :x:                | Количество набранных баллов                |
| certificate_id   | uuid     | certificate_id  | :x:                | :x:                | Внешний ключ на сертификат                 |
| attempt_version  | json     | attempt_version | :x:                | :x:                | Версия попытки (данные о тесте на момент прохождения) |
| attempt_snapshot | varchar  | attempt_snapshot| :x:                | :x:                | Снимок состояния попытки (ответы студента) |
| completed        | boolean  | completed       | :x:                | :x:                | Флаг завершения попытки                    |

## Быстрый доступ
[← На главную](/documentation/02_info_model/)