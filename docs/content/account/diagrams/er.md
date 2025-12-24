---
title: "ER-диаграмма базы данных LMS"
date: 2025-12-15
layout: "single"
---

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