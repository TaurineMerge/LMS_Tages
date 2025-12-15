erDiagram

    TEST_D {
        uuid id PK "not null, unique"
        uuid course_id FK
        varchar title
        integer min_point
        text description
    }

    QUESTION_D {
        uuid id PK "not null, unique"
        uuid test_id FK "not null"
        text text_of_question
        int order
    }

    ANSWER_D {
        uuid id PK "not null, unique"
        text text
        uuid question_id FK "not null"
        int score "not null"
    }

    TEST_ATTEMPT_B {
        uuid id PK "student_id, test_id, date_of_attempt"
        uuid student_id "not null"
        uuid test_id "not null"
        date date_of_attempt
        int point
        uuid certificate_id
        json attempt_version
    }

    %% RELATIONSHIPS

    TEST_D ||--o{ QUESTION_D : "has questions"
    QUESTION_D ||--o{ ANSWER_D : "has answers"

    TEST_D ||--o{ TEST_ATTEMPT_B : "has attempts"