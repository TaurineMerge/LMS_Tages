-- =============================================
-- МИГРАЦИЯ 002: Интеграция S3 и статистика
-- =============================================
-- Дата: 2025-12-17
-- Описание: 
--   1. Обновить таблицу certificate_b для работы с S3
--   2. Добавить таблицы для сырых данных статистики из сервиса testing
--   3. Добавить индексы для быстрого поиска необработанных данных




-- =============================================
-- МИГРАЦИЯ 002: Интеграция S3 и статистика
-- =============================================
-- Дата: 2025-12-17
-- Описание: 
--   1. Обновить таблицу certificate_b для работы с S3
--   2. Добавить таблицы для сырых данных статистики из сервиса testing
--   3. Добавить индексы для быстрого поиска необработанных данных

-- =============================================
-- ЧАСТЬ 0: Создание базовых схем и таблиц (дублирование из migrate-001.sql)
-- =============================================

-- 0.1. Создать схему personal_account
CREATE SCHEMA IF NOT EXISTS personal_account;

-- 0.2. Создать таблицу student_s (полная структура из migrate-001.sql)
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

-- 0.3. Создать таблицу certificate_b (базовая, из migrate-001.sql)
CREATE TABLE IF NOT EXISTS personal_account.certificate_b (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID REFERENCES personal_account.student_s(id) ON DELETE CASCADE,
    certificate_number SERIAL ,
    pdf_s3_key VARCHAR(500) ,
    snapshot_s3_key VARCHAR(500) ,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 0.4. Создать схему knowledge_base
CREATE SCHEMA IF NOT EXISTS knowledge_base;

-- 0.5. Создать таблицу category_d
CREATE TABLE IF NOT EXISTS knowledge_base.category_d (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 0.6. Создать таблицу course_b
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

-- 0.7. Создать схему tests
CREATE SCHEMA IF NOT EXISTS tests;

-- 0.8. Создать таблицу test_d
CREATE TABLE IF NOT EXISTS tests.test_d (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    course_id UUID REFERENCES knowledge_base.course_b(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    min_point INTEGER DEFAULT 0,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 0.9. Создать таблицу test_attempt_b (базовая, из migrate-001.sql)
CREATE TABLE IF NOT EXISTS tests.test_attempt_b (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES personal_account.student_s(id) ON DELETE CASCADE,
    test_id UUID NOT NULL REFERENCES tests.test_d(id) ON DELETE CASCADE,
    date_of_attempt DATE NOT NULL,
    point INTEGER,
    passed BOOLEAN,
    completed BOOLEAN DEFAULT FALSE,
    result JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);



-- =============================================
-- ЧАСТЬ 1: Обновление сертификатов (S3)
-- =============================================

-- 1.1. Добавить новые колонки для S3 (если их ещё нет)
-- issuer_id and status columns removed by request; keep only S3-related fields managed elsewhere

-- 1.2. Удалить устаревшие колонки (если они существуют)
-- Примечание: эти колонки были в старой версии, теперь их нет
ALTER TABLE personal_account.certificate_b
    DROP COLUMN IF EXISTS content,
    DROP COLUMN IF EXISTS course_id,
    DROP COLUMN IF EXISTS test_attempt_id;

-- 1.3. Убедиться, что S3 колонки объявлены как NOT NULL
-- (выполняется только если они уже существуют как nullable)

-- 1.4. Добавить индексы для сертификатов
CREATE INDEX IF NOT EXISTS idx_certificate_student ON personal_account.certificate_b(student_id);
CREATE INDEX IF NOT EXISTS idx_certificate_number ON personal_account.certificate_b(certificate_number);
CREATE INDEX IF NOT EXISTS idx_certificate_created ON personal_account.certificate_b(created_at);

-- =============================================
-- ЧАСТЬ 2: Таблицы для сырых данных статистики
-- =============================================

-- 2.1. Создать schema для интеграций
CREATE SCHEMA IF NOT EXISTS integration;

-- 2.2. Таблица для сырых статистик пользователей (из UserStats контракта)
-- Эти данные приходят с сервиса testing при запросе статистики
CREATE TABLE IF NOT EXISTS integration.raw_user_stats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL REFERENCES personal_account.student_s(id) ON DELETE CASCADE,
    payload JSONB NOT NULL,
    -- Payload содержит:
    -- {
    --   "user_id": "uuid",
    --   "attempts_total": int,
    --   "attempts_passed": int,
    --   "best_score": int | null,
    --   "last_attempt_at": datetime | null,
    --   "per_test": [{ "test_id", "test_title", "attempts", "best_score", "passed_count" }]
    -- }
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN DEFAULT FALSE,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2.3. Таблица для сырых попыток тестов (из AttemptDetail контракта)
-- Эти данные приходят при завершении/обновлении попытки
CREATE TABLE IF NOT EXISTS integration.raw_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    external_attempt_id UUID NOT NULL UNIQUE,
    -- ID попытки из сервиса testing (уникальный)
    student_id UUID NOT NULL REFERENCES personal_account.student_s(id) ON DELETE CASCADE,
    test_id UUID,
    -- Может быть неизвестен на начальном этапе (заполнится при обработке)
    payload JSONB NOT NULL,
    -- Payload содержит:
    -- {
    --   "attempt_id": "uuid",
    --   "student_id": "uuid",
    --   "test_id": "uuid",
    --   "date_of_attempt": "date",
    --   "point": int | null,
    --   "completed": bool,
    --   "passed": bool | null,
    --   "certificate_id": "uuid" | null,
    --   "attempt_version": object | null (S3 key or version data),
    --   "attempt_snapshot_s3": string | null (S3 path),
    --   "meta": object | null
    -- }
    received_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed BOOLEAN DEFAULT FALSE,
    processing_attempts INTEGER DEFAULT 0,
    -- Счётчик попыток обработки (для отладки)
    error_message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- ЧАСТЬ 3: Индексы для эффективного поиска
-- =============================================

-- 3.1. Индексы для raw_user_stats
CREATE INDEX IF NOT EXISTS idx_raw_user_stats_processed 
    ON integration.raw_user_stats(processed, received_at);
CREATE INDEX IF NOT EXISTS idx_raw_user_stats_student 
    ON integration.raw_user_stats(student_id);
CREATE INDEX IF NOT EXISTS idx_raw_user_stats_received 
    ON integration.raw_user_stats(received_at DESC);

-- 3.2. Индексы для raw_attempts
CREATE INDEX IF NOT EXISTS idx_raw_attempts_processed 
    ON integration.raw_attempts(processed, received_at);
CREATE INDEX IF NOT EXISTS idx_raw_attempts_student 
    ON integration.raw_attempts(student_id);
CREATE INDEX IF NOT EXISTS idx_raw_attempts_external_id 
    ON integration.raw_attempts(external_attempt_id);
CREATE INDEX IF NOT EXISTS idx_raw_attempts_test 
    ON integration.raw_attempts(test_id) WHERE test_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_raw_attempts_received 
    ON integration.raw_attempts(received_at DESC);

-- =============================================
-- ЧАСТЬ 4: Добавить поля в test_attempt_b
-- =============================================

-- 4.1. Добавить поля если их нет (для полноты контракта AttemptDetail)
ALTER TABLE tests.test_attempt_b
    ADD COLUMN IF NOT EXISTS completed BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS passed BOOLEAN,
    ADD COLUMN IF NOT EXISTS certificate_id UUID 
        REFERENCES personal_account.certificate_b(id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS attempt_snapshot_s3 TEXT,
    ADD COLUMN IF NOT EXISTS attempt_version JSONB,
    ADD COLUMN IF NOT EXISTS meta JSONB DEFAULT '{}'::jsonb;

-- 4.2. Добавить индекс на certificate_id
CREATE INDEX IF NOT EXISTS idx_test_attempt_certificate 
    ON tests.test_attempt_b(certificate_id) WHERE certificate_id IS NOT NULL;

-- =============================================
-- ЧАСТЬ 5: Триггеры для updated_at
-- =============================================

-- 5.1. Добавить триггеры обновления updated_at для новых таблиц
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Триггер для raw_user_stats
DROP TRIGGER IF EXISTS update_raw_user_stats_updated_at ON integration.raw_user_stats;
CREATE TRIGGER update_raw_user_stats_updated_at
    BEFORE UPDATE ON integration.raw_user_stats
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Триггер для raw_attempts
DROP TRIGGER IF EXISTS update_raw_attempts_updated_at ON integration.raw_attempts;
CREATE TRIGGER update_raw_attempts_updated_at
    BEFORE UPDATE ON integration.raw_attempts
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Триггер для certificate_b (если не существует)
DROP TRIGGER IF EXISTS update_certificate_updated_at ON personal_account.certificate_b;
CREATE TRIGGER update_certificate_updated_at
    BEFORE UPDATE ON personal_account.certificate_b
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- =============================================
-- ЧАСТЬ 6: Представления (Views) для удобства
-- =============================================

-- 6.1. View для получения необработанных попыток с информацией о студенте
CREATE OR REPLACE VIEW integration.v_unprocessed_attempts AS
SELECT
    ra.id as raw_id,
    ra.external_attempt_id,
    ra.student_id,
    s.email as student_email,
    s.name as student_name,
    ra.test_id,
    ra.payload,
    ra.received_at,
    ra.processing_attempts,
    ra.error_message
FROM integration.raw_attempts ra
JOIN personal_account.student_s s ON ra.student_id = s.id
WHERE ra.processed = FALSE
ORDER BY ra.received_at ASC;

-- 6.2. View для получения необработанных статистик с информацией о студенте
CREATE OR REPLACE VIEW integration.v_unprocessed_user_stats AS
SELECT
    rus.id as raw_id,
    rus.student_id,
    s.email as student_email,
    s.name as student_name,
    rus.payload,
    rus.received_at,
    rus.error_message
FROM integration.raw_user_stats rus
JOIN personal_account.student_s s ON rus.student_id = s.id
WHERE rus.processed = FALSE
ORDER BY rus.received_at ASC;

-- =============================================
-- ЧАСТЬ 7: Комментарии для документации
-- =============================================

COMMENT ON SCHEMA integration IS 'Схема для интеграций с внешними сервисами (testing и др.)';

COMMENT ON TABLE integration.raw_user_stats IS 
'Таблица для хранения сырых данных статистики пользователя из сервиса testing.
Данные здесь хранятся как есть (полный JSON), затем периодически обрабатываются воркером.
После обработки флаг processed устанавливается в TRUE.';

COMMENT ON TABLE integration.raw_attempts IS 
'Таблица для хранения сырых данных попыток тестов из сервиса testing.
Каждая новая/обновленная попытка сохраняется здесь перед обработкой.
Уникальное ограничение на external_attempt_id предотвращает дубликаты.';

COMMENT ON COLUMN integration.raw_attempts.external_attempt_id IS 
'ID попытки из сервиса testing (уникальный идентификатор для интеграции)';

COMMENT ON COLUMN integration.raw_attempts.test_id IS 
'Может быть NULL до обработки. Заполняется воркером при связывании с нашей БД';

COMMENT ON COLUMN integration.raw_attempts.processing_attempts IS 
'Счётчик попыток обработки. Помогает отследить, сколько раз обработчик пытался обработать эту запись';

-- =============================================
-- ЧАСТЬ 8: Проверка целостности данных
-- =============================================

-- 8.1. Убедиться, что все необходимые индексы существуют
-- (PostgreSQL автоматически пропустит существующие индексы)

-- 8.2. Логирование факта успешного выполнения миграции
-- Вывести информацию о созданных таблицах и индексах
DO $$
BEGIN
    RAISE NOTICE '✓ Migration 002 completed successfully:';
    RAISE NOTICE '  - Updated certificate_b table with S3 fields';
    RAISE NOTICE '  - Created integration schema';
    RAISE NOTICE '  - Created raw_user_stats table';
    RAISE NOTICE '  - Created raw_attempts table';
    RAISE NOTICE '  - Added indexes for performance';
    RAISE NOTICE '  - Created views for convenience';
END $$;

-- =============================================
-- МИГРАЦИЯ ЗАВЕРШЕНА
-- =============================================
-- Next steps:
-- 1. Запустить миграцию: docker-compose up personal-account-db
-- 2. Проверить таблицы: \dt integration.* в psql
-- 3. Создать Python репозитории для работы с этими таблицами
-- 4. Реализовать webhook для приёма данных из testing
-- 5. Создать воркер для обработки сырых данных