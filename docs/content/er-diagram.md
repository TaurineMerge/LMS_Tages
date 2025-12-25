---
title: "ER-диаграмма базы данных LMS"
date: 2025-12-15
layout: "single"
---

## Диаграмма БД админ-панель и публичный сайт

{{< mermaid >}}

erDiagram
    category_d {
        UUID id PK "not null, default gen_random_uuid()"
        VARCHAR(255) title "not null"
        TIMESTAMP created_at "not null, default NOW()"
        TIMESTAMP updated_at "not null, default NOW()"
    }
    course_b {
        UUID id PK "not null, default gen_random_uuid()"
        VARCHAR(255) title "not null"
        TEXT description
        VARCHAR(20) level "not null, check: ['hard', 'medium', 'easy']"
        UUID category_id FK "not null"
        VARCHAR(500) image_key
        VARCHAR(20) visibility "not null, check: ['draft', 'public'], default 'draft'"
        TIMESTAMP created_at "not null, default NOW()"
        TIMESTAMP updated_at "not null, default NOW()"
    }
    lesson_d {
        UUID id PK "not null, default gen_random_uuid()"
        VARCHAR(255) title "not null"
        UUID course_id FK "not null"
        TEXT content "not null, default 'draft'"
        TIMESTAMP created_at "not null, default NOW()"
        TIMESTAMP updated_at "not null, default NOW()"
    }
    category_d ||--o{ course_b : "contains"
    course_b   ||--o{ lesson_d : "contains"

{{< /mermaid >}}

---
## Диаграмма БД тестов

{{< mermaid >}}

erDiagram

    test_d {
        uuid id PK "not null, unique"
        uuid course_id FK
        varchar title
        integer min_point
        text description
    }

    question_d {
        uuid id PK "not null, unique"
        uuid test_id FK "null"
        uuid draft_id FK "null"
        text text_of_question
        int order
    }

    answer_d {
        uuid id PK "not null, unique"
        text text
        uuid question_id FK "not null"
        int score "not null"
    }

    test_attempt_b {
        uuid id PK "student_id, test_id, date_of_attempt"
        uuid student_id "not null"
        uuid test_id "not null"
        date date_of_attempt
        int point
        uuid certificate_id
        json attempt_version
        varchar attempt_snapshot
        boolen completed
    }


    draft_b {
        uuid id PK "not null, unique"
        varchar title
        integer min_point
        text description
        uuid test_id "null"
    }


    content_d {
        uuid id PK "not null, unique"
        int order
        text content
        boolen type_of_content
        uuid question_id FK "not null"
        uuid answer_id FK "not null"
    }

    %% RELATIONSHIPS

    test_d ||--o{ question_d : "has questions"
    question_d ||--o{ answer_d : "has answers"
    test_d ||--o{ test_attempt_b : "has attempts"
    draft_b o|--o| test_d : "has drafts"
    content_d }o--o| question_d: "has content"
    content_d }o--o| answer_d: "has content"
    draft_b ||--o{ question_d: "has questions"

{{< /mermaid >}}

---
## Диаграмма БД ЛК

{{< mermaid >}}

erDiagram
    student_s {
        uuid id PK "not null, unique"
        varchar name
        varchar surname
        date birth_date
        varchar avatar
        json contacts
        varchar email
        varchar phone
    }
    
    certificate_b {
        uuid id PK "not null, unique"
        integer certificate_number
        date created_at
        varchar content
        uuid student_id FK
        uuid course_id FK
    }
    
    test_attempt_b {
        uuid id PK "student_id, test_id, date_of_attempt"
        uuid student_id "not null"
        uuid test_id "not null"
        date date_of_attempt
        int point
        uuid certificate_id
        json attempt_version
    }
    
    
    student_s ||--o{ certificate_b : "receives certificates"
    certificate_b ||--o| test_attempt_b : "issued for attempt"

{{< /mermaid >}}