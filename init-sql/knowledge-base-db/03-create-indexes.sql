CREATE INDEX IF NOT EXISTS idx_course_created_at ON knowledge_base.course_b (created_at);
CREATE INDEX IF NOT EXISTS idx_course_level ON knowledge_base.course_b (level);
CREATE INDEX IF NOT EXISTS idx_course_visibility ON knowledge_base.course_b (visibility);

CREATE INDEX IF NOT EXISTS idx_lesson_created_at ON knowledge_base.lesson_d (created_at);
