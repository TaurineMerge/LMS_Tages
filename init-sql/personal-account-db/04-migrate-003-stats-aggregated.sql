-- Migration: Create aggregated stats table
-- This table stores pre-computed statistics for each student

CREATE SCHEMA IF NOT EXISTS stats;

CREATE TABLE IF NOT EXISTS stats.student_stats_aggregated (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    student_id UUID NOT NULL UNIQUE REFERENCES personal_account.student_s(id) ON DELETE CASCADE,
    total_attempts INTEGER DEFAULT 0,
    passed_attempts INTEGER DEFAULT 0,
    failed_attempts INTEGER DEFAULT 0,
    avg_score NUMERIC(5,2) DEFAULT 0,
    total_tests_taken INTEGER DEFAULT 0,
    last_attempt_at TIMESTAMP WITH TIME ZONE,
    stats_json JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_student_stats_student_id ON stats.student_stats_aggregated(student_id);
CREATE INDEX IF NOT EXISTS idx_student_stats_updated_at ON stats.student_stats_aggregated(updated_at);