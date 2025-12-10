---
title: "ER"
date: 2025-12-09
---

```mermaid
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
    
    course_b {
        uuid id PK "not null, unique"
        varchar title
        text description
        varchar level "hard, medium, easy"
        uuid category_id FK "not null"
        varchar visibility "draft, public, private"
    }
    
    category_d {
        uuid id PK "not null, unique"
        varchar title
    }
    
    lesson_d {
        uuid id PK "not null, unique"
        varchar title
        uuid course_id FK "not null"
        json content
    }
    
    test_d {
        uuid id PK "not null, unique"
        uuid course_id FK
        varchar title
        integer min_point
        text description
    }
    
    question_d {
        uuid id PK "not null, unique"
        uuid test_id FK "not null"
        text text_of_question
        int order
    }
    
    answer_d {
        uuid id PK "not null, unique"
        text text
        uuid question_id FK "not null"
        int score "not null"
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
    
    visit_students_for_lessons {
        uuid id PK "not null, unique"
        uuid student_id "not null, unique"
        uuid lesson_id FK "not null, unique"
    }
    
    course_b ||--o{ lesson_d : "has lessons"
    course_b ||--|| category_d : "belongs to category"
    course_b ||--o| test_d : "has test"
    test_d ||--|{ question_d : "contains questions"
    question_d ||--|{ answer_d : "has answers"
    student_s ||--o{ certificate_b : "receives certificates"
    student_s ||--o{ test_attempt_b : "attempts tests"
    student_s }|--|| visit_students_for_lessons : "visits lessons"
    lesson_d ||--o{ visit_students_for_lessons : "visited by students"
    test_d }o--|{ test_attempt_b : "attempted by students"
    certificate_b ||--o| test_attempt_b : "issued for attempt"