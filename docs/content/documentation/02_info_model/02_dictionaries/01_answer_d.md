---
title: "Answer"
description: "Описание словаря 'answer_d'"
version: 1
draft: false
weight: 3
    
---

Словарь, содержащий список ответов.

| Attribute |    Type     | DB column |        MUST        |      Not Null      | Comments                                    |
|:----------|:-----------:|:----------|:------------------:|:------------------:|:--------------------------------------------|
|id|	uuid	|id	|:white_check_mark:|	:white_check_mark:	|Уникальный идентификатор ответа|
|text	|text	|text|	:x:	|:x:	|Текст ответа|
|question_id|	uuid|	question_id|	:white_check_mark:	|:white_check_mark:	Внешний ключ на вопрос|
|score	|int	|score	|:white_check_mark:	|:white_check_mark:	|Количество баллов за выбор этого ответа|


## Быстрый доступ
[← На главную](/documentation/02_info_model/)