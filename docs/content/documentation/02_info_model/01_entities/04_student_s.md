---
title: "Student"
description: "Описание сущности 'student_s'"
version: 1
draft: false
weight: 3
    
---

## Сущность студент

| Attribute | Type     | DB column | MUST               | Not Null           | Comments                                   |
|:----------|:--------:|:----------|:------------------:|:------------------:|:-------------------------------------------|
| id        | uuid     | id        | :white_check_mark: | :white_check_mark: | Уникальный идентификатор студента          |
| name      | varchar  | name      | :x:                | :x:                | Имя студента                               |
| surname   | varchar  | surname   | :x:                | :x:                | Фамилия студента                           |
| birth_date| date     | birth_date| :x:                | :x:                | Дата рождения студента                     |
| avatar    | varchar  | avatar    | :x:                | :x:                | Ссылка на аватар студента                  |
| contacts  | json     | contacts  | :x:                | :x:                | Контактная информация в формате JSON       |
| email     | varchar  | email     | :white_check_mark: | :white_check_mark: | Email студента                             |
| phone     | varchar  | phone     | :x:                | :x:                | Телефон студента                           |


## Быстрый доступ
[← На главную](/documentation/02_info_model/)