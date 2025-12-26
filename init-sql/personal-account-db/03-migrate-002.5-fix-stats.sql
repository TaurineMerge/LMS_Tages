-- =============================================
-- ФИКС СХЕМЫ ДЛЯ СТАТИСТИКИ
-- =============================================

CREATE SCHEMA IF NOT EXISTS tests;

-- 1. Добавляем недостающие колонки в попытки тестов
ALTER TABLE tests.test_attempt_b 
    ADD COLUMN IF NOT EXISTS completed BOOLEAN DEFAULT FALSE,
    ADD COLUMN IF NOT EXISTS passed BOOLEAN,
    ADD COLUMN IF NOT EXISTS attempt_snapshot_s3 VARCHAR(500),
    ADD COLUMN IF NOT EXISTS attempt_version JSONB DEFAULT '{}'::jsonb,
    ADD COLUMN IF NOT EXISTS meta JSONB DEFAULT '{}'::jsonb;

-- 2. Создаем схему для агрегированной статистики
CREATE SCHEMA IF NOT EXISTS stats;

-- 3. Таблица агрегированной статистики студентов
CREATE TABLE IF NOT EXISTS stats.student_stats_aggregated (
    student_id UUID PRIMARY KEY REFERENCES personal_account.student_s(id) ON DELETE CASCADE,
    total_attempts INTEGER DEFAULT 0,
    passed_attempts INTEGER DEFAULT 0,
    failed_attempts INTEGER DEFAULT 0,
    avg_score DOUBLE PRECISION DEFAULT 0.0,
    total_tests_taken INTEGER DEFAULT 0,
    last_attempt_at DATE,
    stats_json JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индексы для статистики
CREATE INDEX IF NOT EXISTS idx_stats_updated_at ON stats.student_stats_aggregated(updated_at);
