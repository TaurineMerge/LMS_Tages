CREATE INDEX IF NOT EXISTS idx_question_test_id ON testing.question_d(test_id);
CREATE INDEX IF NOT EXISTS idx_answer_question_id ON testing.answer_d(question_id);
CREATE INDEX IF NOT EXISTS idx_attempt_test_id ON testing.test_attempt_b(test_id);
CREATE INDEX IF NOT EXISTS idx_attempt_student_id ON testing.test_attempt_b(student_id);
CREATE INDEX IF NOT EXISTS idx_draft_test_id ON testing.draft_b(test_id);
CREATE INDEX IF NOT EXISTS idx_content_question_id ON testing.content_d(question_id);
CREATE INDEX IF NOT EXISTS idx_content_answer_id ON testing.content_d(answer_id);