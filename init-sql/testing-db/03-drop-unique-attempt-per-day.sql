-- Позволяет делать несколько попыток одного и того же теста в один день.
-- Если ограничение уже было создано ранее, удаляем его.

ALTER TABLE IF EXISTS testing.test_attempt_b
  DROP CONSTRAINT IF EXISTS test_attempt_b_student_id_test_id_date_of_attempt_key;

-- Индекс для быстрых выборок попыток студента по конкретному тесту (НЕ UNIQUE)
CREATE INDEX IF NOT EXISTS idx_ta_student_test_date
  ON testing.test_attempt_b (student_id, test_id, date_of_attempt);
