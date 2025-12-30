-- ============================================================
-- SEED v2 для схемы testing (single + multi + edge cases)
-- Формат attempt_version: attemptNo + answers[{order, questionId, answerId, answerIds}]
-- ============================================================

CREATE EXTENSION IF NOT EXISTS pgcrypto;

TRUNCATE TABLE
  testing.content_d,
  testing.test_attempt_b,
  testing.answer_d,
  testing.question_d,
  testing.draft_b,
  testing.test_d
CASCADE;

DO $$
DECLARE
  -- tests
  t_mixed UUID;
  t_math  UUID;

  -- questions (mixed)
  q1 UUID; q2 UUID; q3 UUID; q4 UUID;

  -- questions (math)
  q5 UUID;

  -- answers ids
  q1_a_ok UUID;
  q2_a_ok1 UUID; q2_a_ok2 UUID; q2_a_wrong UUID;
  q3_a_any UUID;
  q4_a_neg UUID; q4_a_ok UUID;
  q5_a_ok UUID;

  -- any answer for content links
  q1_any UUID; q2_any UUID; q3_any UUID; q4_any UUID; q5_any UUID;

  -- students
  s1 UUID := gen_random_uuid();
  s2 UUID := gen_random_uuid();

  -- attempt ids
  a1 UUID := gen_random_uuid();
  a2 UUID := gen_random_uuid();
  a3 UUID := gen_random_uuid();
BEGIN
  -- -----------------------------
  -- TESTS
  -- -----------------------------
  INSERT INTO testing.test_d (course_id, title, min_point, description)
  VALUES (gen_random_uuid(), 'MIXED: single + multi + edge cases', 4, 'Проверка интерфейса single/multi и подсчёта баллов')
  RETURNING id INTO t_mixed;

  INSERT INTO testing.test_d (course_id, title, min_point, description)
  VALUES (gen_random_uuid(), 'Math: одно правильное', 1, 'Короткий тест для smoke-check')
  RETURNING id INTO t_math;

  -- drafts (не обязательно, но пусть будут)
  INSERT INTO testing.draft_b (title, min_point, description, test_id)
  VALUES ('Draft: MIXED v0', 4, 'Черновик', t_mixed);

  INSERT INTO testing.draft_b (title, min_point, description, test_id)
  VALUES ('Draft: MATH v0', 1, 'Черновик', t_math);

  -- -----------------------------
  -- QUESTIONS (MIXED)
  -- -----------------------------

  -- Q1 single (1 correct)
  INSERT INTO testing.question_d (test_id, text_of_question, "order")
  VALUES (t_mixed, 'Q1 (single): 2 + 2 = ?', 1)
  RETURNING id INTO q1;

  INSERT INTO testing.answer_d (text, question_id, score) VALUES
    ('3', q1, 0),
    ('4', q1, 2),   -- правильный (score>0) и сразу 2 балла
    ('5', q1, 0);

  SELECT id INTO q1_a_ok FROM testing.answer_d WHERE question_id=q1 AND score>0 ORDER BY score DESC LIMIT 1;
  SELECT id INTO q1_any  FROM testing.answer_d WHERE question_id=q1 ORDER BY id LIMIT 1;

  -- Q2 multi (2 correct)
  INSERT INTO testing.question_d (test_id, text_of_question, "order")
  VALUES (t_mixed, 'Q2 (multi): выбери простые числа', 2)
  RETURNING id INTO q2;

  INSERT INTO testing.answer_d (text, question_id, score) VALUES
    ('2', q2, 1),   -- correct
    ('3', q2, 1),   -- correct
    ('4', q2, 0),   -- wrong
    ('9', q2, 0);   -- wrong

  SELECT id INTO q2_a_ok1 FROM testing.answer_d WHERE question_id=q2 AND text='2' LIMIT 1;
  SELECT id INTO q2_a_ok2 FROM testing.answer_d WHERE question_id=q2 AND text='3' LIMIT 1;
  SELECT id INTO q2_a_wrong FROM testing.answer_d WHERE question_id=q2 AND text='4' LIMIT 1;
  SELECT id INTO q2_any   FROM testing.answer_d WHERE question_id=q2 ORDER BY id LIMIT 1;

  -- Q3 trick (0 correct -> UI будет как single, но баллов не даст)
  INSERT INTO testing.question_d (test_id, text_of_question, "order")
  VALUES (t_mixed, 'Q3 (trick): у этого вопроса 0 правильных (все score=0)', 3)
  RETURNING id INTO q3;

  INSERT INTO testing.answer_d (text, question_id, score) VALUES
    ('вариант A', q3, 0),
    ('вариант B', q3, 0),
    ('вариант C', q3, 0);

  SELECT id INTO q3_a_any FROM testing.answer_d WHERE question_id=q3 ORDER BY id LIMIT 1;
  SELECT id INTO q3_any   FROM testing.answer_d WHERE question_id=q3 ORDER BY id LIMIT 1;

  -- Q4 negative scores (проверка max(0,score))
  INSERT INTO testing.question_d (test_id, text_of_question, "order")
  VALUES (t_mixed, 'Q4 (negative): выбери лучший вариант (есть отрицательные баллы)', 4)
  RETURNING id INTO q4;

  INSERT INTO testing.answer_d (text, question_id, score) VALUES
    ('плохой вариант (-2)', q4, -2),
    ('нормальный (0)',      q4, 0),
    ('лучший (+2)',         q4, 2);

  SELECT id INTO q4_a_neg FROM testing.answer_d WHERE question_id=q4 AND score=-2 LIMIT 1;
  SELECT id INTO q4_a_ok  FROM testing.answer_d WHERE question_id=q4 AND score=2  LIMIT 1;
  SELECT id INTO q4_any   FROM testing.answer_d WHERE question_id=q4 ORDER BY id LIMIT 1;

  -- -----------------------------
  -- QUESTIONS (MATH)
  -- -----------------------------
  INSERT INTO testing.question_d (test_id, text_of_question, "order")
  VALUES (t_math, 'Q5 (single): Реши 5x=20. x=?', 1)
  RETURNING id INTO q5;

  INSERT INTO testing.answer_d (text, question_id, score) VALUES
    ('2', q5, 0),
    ('4', q5, 1),
    ('5', q5, 0);

  SELECT id INTO q5_a_ok FROM testing.answer_d WHERE question_id=q5 AND score>0 LIMIT 1;
  SELECT id INTO q5_any  FROM testing.answer_d WHERE question_id=q5 ORDER BY id LIMIT 1;

  -- -----------------------------
  -- CONTENT (для проверки, что контент подтягивается/не ломает страницы)
  -- type_of_content: трактовку ты сама знаешь; я оставляю как есть (TRUE/FALSE)
  -- -----------------------------
  INSERT INTO testing.content_d ("order", content, type_of_content, question_id, answer_id) VALUES
    (1, 'Подсказка: 2+2=4 — базовая арифметика.', TRUE,  q1, q1_any),
    (1, 'Подсказка: простые числа делятся только на 1 и себя.', TRUE,  q2, q2_any),
    (1, 'Этот вопрос специально без правильных.', FALSE, q3, q3_any),
    (1, 'Отрицательные баллы должны игнорироваться при maxPossible (если у тебя так задумано).', TRUE, q4, q4_any);

  -- -----------------------------
  -- TEST ATTEMPTS
  -- attempt_version в НОВОМ формате
  -- -----------------------------

  -- Attempt 1: student1 завершил MIXED (выбрал single+multi+trick+best)
  INSERT INTO testing.test_attempt_b
    (id, student_id, test_id, date_of_attempt, point, certificate_id, attempt_version, attempt_snapshot, completed)
  VALUES
    (a1, s1, t_mixed, DATE '2025-12-20', 2+1+1+0+2, NULL,
      json_build_object(
        'attemptNo', 1,
        'answers', json_build_array(
          json_build_object('order', 1, 'questionId', q1::text, 'answerId', q1_a_ok::text, 'answerIds', json_build_array(q1_a_ok::text)),
          json_build_object('order', 2, 'questionId', q2::text, 'answerId', NULL,         'answerIds', json_build_array(q2_a_ok1::text, q2_a_ok2::text)),
          json_build_object('order', 3, 'questionId', q3::text, 'answerId', q3_a_any::text,'answerIds', json_build_array(q3_a_any::text)),
          json_build_object('order', 4, 'questionId', q4::text, 'answerId', q4_a_ok::text, 'answerIds', json_build_array(q4_a_ok::text))
        )
      ),
      json_build_object(
        'testId', t_mixed::text,
        'questionsCount', 4,
        'note', 'Завершённая попытка для проверки single/multi'
      )::text,
      TRUE
    );

  -- Attempt 2: student2 начал MIXED, ответил только на Q2 (multi) и НЕ завершил
  INSERT INTO testing.test_attempt_b
    (id, student_id, test_id, date_of_attempt, point, certificate_id, attempt_version, attempt_snapshot, completed)
  VALUES
    (a2, s2, t_mixed, NULL, NULL, NULL,
      json_build_object(
        'attemptNo', 1,
        'answers', json_build_array(
          json_build_object('order', 1, 'questionId', q1::text, 'answerId', NULL, 'answerIds', json_build_array()),
          json_build_object('order', 2, 'questionId', q2::text, 'answerId', NULL, 'answerIds', json_build_array(q2_a_ok1::text, q2_a_wrong::text)),
          json_build_object('order', 3, 'questionId', q3::text, 'answerId', NULL, 'answerIds', json_build_array()),
          json_build_object('order', 4, 'questionId', q4::text, 'answerId', NULL, 'answerIds', json_build_array())
        )
      ),
      json_build_object(
        'testId', t_mixed::text,
        'questionsCount', 4,
        'note', 'Незавершённая попытка (для проверки incomplete state)'
      )::text,
      FALSE
    );

  -- Attempt 3: student1 завершил MATH (smoke)
  INSERT INTO testing.test_attempt_b
    (id, student_id, test_id, date_of_attempt, point, certificate_id, attempt_version, attempt_snapshot, completed)
  VALUES
    (a3, s1, t_math, DATE '2025-12-21', 1, NULL,
      json_build_object(
        'attemptNo', 2,
        'answers', json_build_array(
          json_build_object('order', 1, 'questionId', q5::text, 'answerId', q5_a_ok::text, 'answerIds', json_build_array(q5_a_ok::text))
        )
      ),
      json_build_object(
        'testId', t_math::text,
        'questionsCount', 1,
        'note', 'Smoke test'
      )::text,
      TRUE
    );

END $$;

-- Быстрая проверка
SELECT 'test_d' as table_name, COUNT(*) as row_count FROM testing.test_d
UNION ALL SELECT 'question_d', COUNT(*) FROM testing.question_d
UNION ALL SELECT 'answer_d', COUNT(*) FROM testing.answer_d
UNION ALL SELECT 'test_attempt_b', COUNT(*) FROM testing.test_attempt_b
UNION ALL SELECT 'draft_b', COUNT(*) FROM testing.draft_b
UNION ALL SELECT 'content_d', COUNT(*) FROM testing.content_d
ORDER BY table_name;

-- Список тестов
SELECT id, title, min_point, description FROM testing.test_d ORDER BY title;

-- Посмотреть attempt_version красиво
SELECT
  ta.id,
  ta.student_id,
  t.title AS test_title,
  ta.date_of_attempt,
  ta.point,
  ta.completed,
  jsonb_pretty(ta.attempt_version::jsonb) AS attempt_version
FROM testing.test_attempt_b ta
JOIN testing.test_d t ON t.id = ta.test_id
ORDER BY ta.date_of_attempt NULLS LAST, t.title;