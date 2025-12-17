CREATE SCHEMA IF NOT EXISTS testing;

CREATE TABLE IF NOT EXISTS testing.test_d (
    id UUID PRIMARY KEY NOT NULL,
    course_id UUID,
    title VARCHAR,
    min_point INTEGER,
    description TEXT,
    UNIQUE(id)
);

CREATE TABLE IF NOT EXISTS testing.question_d (
    id UUID PRIMARY KEY NOT NULL,
    test_id UUID NOT NULL,
    text_of_question TEXT,
    "order" INTEGER,
    UNIQUE(id),
    CONSTRAINT fk_question_test FOREIGN KEY (test_id) REFERENCES testing.test_d(id)
);

CREATE TABLE IF NOT EXISTS testing.answer_d (
    id UUID PRIMARY KEY NOT NULL,
    text TEXT,
    question_id UUID NOT NULL,
    score INTEGER NOT NULL,
    UNIQUE(id),
    CONSTRAINT fk_answer_question FOREIGN KEY (question_id) REFERENCES testing.question_d(id)
);

CREATE TABLE IF NOT EXISTS testing.test_attempt_b (
    id UUID PRIMARY KEY,
    student_id UUID NOT NULL,
    test_id UUID NOT NULL,
    date_of_attempt DATE,
    point INTEGER,
    certificate_id UUID,
    attempt_version JSON,
    UNIQUE(student_id, test_id, date_of_attempt),
    CONSTRAINT fk_attempt_test FOREIGN KEY (test_id) REFERENCES testing.test_d(id)
);

CREATE TABLE IF NOT EXISTS testing.draft_b (
    id UUID PRIMARY KEY NOT NULL,
    title VARCHAR,
    min_point INTEGER,
    description TEXT,
    test_id UUID NOT NULL,
    UNIQUE(id),
    CONSTRAINT fk_draft_test FOREIGN KEY (test_id) REFERENCES testing.test_d(id)
);

CREATE TABLE IF NOT EXISTS testing.content_d (
    id UUID PRIMARY KEY NOT NULL,
    "order" INTEGER,
    content TEXT,
    type_of_content BOOLEAN,
    question_id UUID NOT NULL,
    answer_id UUID NOT NULL,
    UNIQUE(id),
    CONSTRAINT fk_content_question FOREIGN KEY (question_id) REFERENCES testing.question_d(id),
    CONSTRAINT fk_content_answer FOREIGN KEY (answer_id) REFERENCES testing.answer_d(id)
);