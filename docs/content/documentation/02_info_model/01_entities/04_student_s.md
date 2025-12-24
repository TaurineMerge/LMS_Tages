---
title: "Student"
description: "Описание сущности 'student_s'"
version: 1
draft: false
weight: 3
    
---

## Сущность студент

| Attribute |    Type     | DB column |        MUST        |      Not Null      | Comments                                    |
|:----------|:-----------:|:----------|:------------------:|:------------------:|:--------------------------------------------|
| id        | smallserial | id        | :white_check_mark: | :white_check_mark: | Уникальный идентификатор студента                        |
| name      |    text     | name      | :x: | :x: | Имя студента |
| surname   |    text     | name      | :x: | :x: | Фамилия студента |
| burthday      |    date     | name      | :x: | :x: | Дата рождения студента |
| avatar      |    text     | name      | :x: | :x: | Ссылка на аватар студента |
| contacts     |    text     | name      | :x: | :x: | Контактная информация в формате JSON |
| email     |    text     | name      | :white_check_mark: | :white_check_mark: | Email студента|
| phone      |    text     | name      | :x: | :x: | Телефон студента |


## Быстрый доступ
[← На главную](/documentation/02_info_model/)