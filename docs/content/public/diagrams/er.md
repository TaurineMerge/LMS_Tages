---
title: "ER-диаграмма для схемы knowledge_base в knowledge_base_db"
date: 2025-12-15
layout: "single"
---

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
        VARCHAR(20) visibility "not null, check: ['draft', 'public'], default 'draft'"
        TEXT content "not null, default 'draft'"
        TIMESTAMP created_at "not null, default NOW()"
        TIMESTAMP updated_at "not null, default NOW()"
    }
    category_d ||--o{ course_b : "contains"
    course_b   ||--o{ lesson_d : "contains"

{{< /mermaid >}}