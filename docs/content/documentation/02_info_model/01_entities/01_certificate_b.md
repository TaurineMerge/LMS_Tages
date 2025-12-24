---
title: "Certificate"
description: "Описание сущности 'certificate_b'"
version: 1
draft: false
weight: 3
    
---

## Сущность сертификат

| Attribute        |    Type     | DB column        |        MUST        |      Not Null      | Comments                                    |
|:-----------------|:-----------:|:----------       |:------------------:|:------------------:|:--------------------------------------------|
|id                |	uuid     |	id              |:white_check_mark:  |	:white_check_mark:|	Уникальный идентификатор сертификата|
|certificate_number|integer	     |certificate_number|	:x:	             |:x:	              |Номер сертификата|
|created_at	       |date         |	created_at	    |   :x:	             |:x:	              |Дата выдачи сертификата|
|content	       |varchar      |	content         |	:x:              |	:x:	              |Содержимое сертификата (например, путь к файлу)|
|student_id	       |uuid         |	student_id	    |:white_check_mark:  |	:white_check_mark:|Внешний ключ на студента|
|course_id	       |uuid         |	course_id	    |:white_check_mark:  |	:white_check_mark:|Внешний ключ на курс|

## Быстрый доступ
[← На главную](/documentation/02_info_model/)