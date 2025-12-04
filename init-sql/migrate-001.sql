-- =============================================
-- ДОМЕННАЯ ОБЛАСТЬ: ЛИЧНЫЙ КАБИНЕТ
-- =============================================

CREATE SCHEMA IF NOT EXISTS personal_account;

-- 1.1. Студенты (справочник)
CREATE TABLE IF NOT EXISTS personal_account.student_s (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    surname VARCHAR(100) NOT NULL,
    birth_date DATE,
    avatar VARCHAR(500),
    contacts JSONB DEFAULT '{}'::jsonb,
    email VARCHAR(255) NOT NULL UNIQUE,
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 1.2. Сертификаты (бизнес-сущность)
CREATE TABLE IF NOT EXISTS personal_account.certificate_b (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    certificate_number SERIAL UNIQUE,
    created_at DATE NOT NULL DEFAULT CURRENT_DATE,
    content TEXT,
    student_id UUID REFERENCES personal_account.student_s(id) ON DELETE CASCADE,
    course_id UUID, -- Внешний ключ будет обновлен позже
    test_attempt_id UUID -- Внешний ключ будет обновлен позже
);

-- 1.3. Посещения уроков
CREATE TABLE IF NOT EXISTS personal_account.visit_students_for_lessons (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES personal_account.student_s(id) ON DELETE CASCADE,
    lesson_id UUID NOT NULL, -- Внешний ключ будет обновлен позже
    UNIQUE(student_id, lesson_id)
);

-- =============================================
-- ДОМЕННАЯ ОБЛАСТЬ: БАЗА ЗНАНИЙ
-- =============================================

CREATE SCHEMA IF NOT EXISTS knowledge_base;

-- 2.1. Категории (справочник)
CREATE TABLE IF NOT EXISTS knowledge_base.category_d (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2.2. Курсы (бизнес-сущность)
CREATE TABLE IF NOT EXISTS knowledge_base.course_b (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    level VARCHAR(20) CHECK (level IN ('hard', 'medium', 'easy')) DEFAULT 'medium',
    category_id UUID NOT NULL REFERENCES knowledge_base.category_d(id) ON DELETE RESTRICT,
    visibility VARCHAR(20) CHECK (visibility IN ('draft', 'public', 'private')) DEFAULT 'draft',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2.3. Уроки (документы)
CREATE TABLE IF NOT EXISTS knowledge_base.lesson_d (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    course_id UUID NOT NULL REFERENCES knowledge_base.course_b(id) ON DELETE CASCADE,
    content JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- ДОМЕННАЯ ОБЛАСТЬ: ТЕСТЫ
-- =============================================

CREATE SCHEMA IF NOT EXISTS tests;

-- 3.1. Тесты (документы)
CREATE TABLE IF NOT EXISTS tests.test_d (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES knowledge_base.course_b(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    min_point INTEGER DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3.2. Вопросы тестов (документы)
CREATE TABLE IF NOT EXISTS tests.question_d (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    test_id UUID NOT NULL REFERENCES tests.test_d(id) ON DELETE CASCADE,
    text_of_question TEXT NOT NULL,
    "order" INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3.3. Варианты ответов (справочник)
CREATE TABLE IF NOT EXISTS tests.answer_d (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    text TEXT NOT NULL,
    question_id UUID NOT NULL REFERENCES tests.question_d(id) ON DELETE CASCADE,
    score INTEGER NOT NULL DEFAULT 0,
    "order" INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 3.4. Попытки тестов (бизнес-сущность)
CREATE TABLE IF NOT EXISTS tests.test_attempt_b (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES personal_account.student_s(id) ON DELETE CASCADE,
    test_id UUID NOT NULL REFERENCES tests.test_d(id) ON DELETE CASCADE,
    date_of_attempt DATE NOT NULL DEFAULT CURRENT_DATE,
    point INTEGER DEFAULT 0,
    result JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- ОБНОВЛЕНИЕ ВНЕШНИХ КЛЮЧЕЙ
-- =============================================

-- Обновляем внешние ключи для certificate_b
ALTER TABLE personal_account.certificate_b 
    ADD CONSTRAINT fk_certificate_course 
    FOREIGN KEY (course_id) 
    REFERENCES knowledge_base.course_b(id) ON DELETE CASCADE;

ALTER TABLE personal_account.certificate_b 
    ADD CONSTRAINT fk_certificate_test_attempt 
    FOREIGN KEY (test_attempt_id) 
    REFERENCES tests.test_attempt_b(id) ON DELETE CASCADE;

-- Обновляем внешний ключ для visit_students_for_lessons
ALTER TABLE personal_account.visit_students_for_lessons 
    ADD CONSTRAINT fk_visit_lesson 
    FOREIGN KEY (lesson_id) 
    REFERENCES knowledge_base.lesson_d(id) ON DELETE CASCADE;

-- =============================================
-- ИНДЕКСЫ
-- =============================================

-- Индексы для личного кабинета
CREATE INDEX IF NOT EXISTS idx_student_email ON personal_account.student_s(email);
CREATE INDEX IF NOT EXISTS idx_certificate_student ON personal_account.certificate_b(student_id);
CREATE INDEX IF NOT EXISTS idx_certificate_course ON personal_account.certificate_b(course_id);
CREATE INDEX IF NOT EXISTS idx_certificate_number ON personal_account.certificate_b(certificate_number);
CREATE INDEX IF NOT EXISTS idx_visit_student ON personal_account.visit_students_for_lessons(student_id);
CREATE INDEX IF NOT EXISTS idx_visit_lesson ON personal_account.visit_students_for_lessons(lesson_id);

-- Индексы для базы знаний
CREATE INDEX IF NOT EXISTS idx_course_category ON knowledge_base.course_b(category_id);
CREATE INDEX IF NOT EXISTS idx_course_visibility ON knowledge_base.course_b(visibility);
CREATE INDEX IF NOT EXISTS idx_lesson_course ON knowledge_base.lesson_d(course_id);

-- Индексы для тестов
CREATE INDEX IF NOT EXISTS idx_test_course ON tests.test_d(course_id);
CREATE INDEX IF NOT EXISTS idx_question_test ON tests.question_d(test_id);
CREATE INDEX IF NOT EXISTS idx_question_order ON tests.question_d(test_id, "order");
CREATE INDEX IF NOT EXISTS idx_answer_question ON tests.answer_d(question_id);
CREATE INDEX IF NOT EXISTS idx_attempt_student ON tests.test_attempt_b(student_id);
CREATE INDEX IF NOT EXISTS idx_attempt_test ON tests.test_attempt_b(test_id);
CREATE INDEX IF NOT EXISTS idx_attempt_composite ON tests.test_attempt_b(student_id, test_id, date_of_attempt);

-- =============================================
-- ТРИГГЕР ДЛЯ ОБНОВЛЕНИЯ UPDATED_AT
-- =============================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Создаем триггеры для всех таблиц с updated_at
DO $$ 
DECLARE 
    tbl text;
    schema_name text;
    full_table_name text;
BEGIN 
    FOR schema_name, tbl IN 
        SELECT nspname, relname 
        FROM pg_class c
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE nspname IN ('personal_account', 'knowledge_base', 'tests')
        AND relkind = 'r'
        AND relname IN (
            'student_s', 'course_b', 'lesson_d', 'test_d', 
            'question_d', 'answer_d', 'certificate_b', 'test_attempt_b', 'category_d'
        )
    LOOP
        full_table_name := schema_name || '.' || tbl;
        EXECUTE format('
            DROP TRIGGER IF EXISTS update_%s_updated_at ON %s;
            CREATE TRIGGER update_%s_updated_at
            BEFORE UPDATE ON %s
            FOR EACH ROW
            EXECUTE FUNCTION update_updated_at_column();
        ', tbl, full_table_name, tbl, full_table_name);
    END LOOP;
END $$;