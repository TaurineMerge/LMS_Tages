---
title: "ER-диаграмма для схемы testing"
date: 2025-12-15
layout: "single"
---

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