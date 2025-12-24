CREATE SCHEMA IF NOT EXISTS knowledge_base;

CREATE TABLE IF NOT EXISTS knowledge_base.category_d (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS knowledge_base.course_b (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    level VARCHAR(20) NOT NULL CHECK (level IN ('hard', 'medium', 'easy')),
    visibility VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (visibility IN ('draft', 'public')),
    category_id UUID NOT NULL REFERENCES knowledge_base.category_d(id) ON DELETE RESTRICT,
    image_key VARCHAR(500),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS knowledge_base.lesson_d (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL DEFAULT '',
    course_id UUID NOT NULL REFERENCES knowledge_base.course_b(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
