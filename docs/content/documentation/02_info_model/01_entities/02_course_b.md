---
title: "Course"
description: "Описание сущности 'course_b'"
version: 1
draft: false
weight: 3
    
---

## Сущность курс

| Attribute |    Type     | DB column |        MUST        |      Not Null      | Comments                                    |
|:----------|:-----------:|:----------|:------------------:|:------------------:|:--------------------------------------------|
|id	        |uuid         |	id	|:white_check_mark:	|:white_check_mark:	|Уникальный идентификатор курса|
|title      |varchar(255) |	title|	:white_check_mark:|	:white_check_mark:	|Название курса|
|description|	text	  |description	|:x:	|:x:	|Описание курса|
|level	    |varchar(20)  |level	|:white_check_mark:	|:white_check_mark:|	Уровень сложности: hard, medium, easy|
|category_id|      uuid   |category_id	|:white_check_mark:	|:white_check_mark:	|Внешний ключ на категорию|
|image_key  |varchar(500) |	image_key	|:x:	|:x:	|Ключ изображения в S3 хранилище|
|visibility	|varchar(20)  |visibility	|:white_check_mark:	|:white_check_mark:	|Видимость курса: draft, public|
|created_at	|timestamp	  |created_at	|:white_check_mark:	|:white_check_mark:|	Дата и время создания|
|updated_at	|timestamp	  |updated_at	|:white_check_mark:	|:white_check_mark:	|Дата и время последнего обновления|

## Быстрый доступ
[← На главную](/documentation/02_info_model/)