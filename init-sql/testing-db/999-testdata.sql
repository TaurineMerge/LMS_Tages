-- ============================================================
-- SEED для схемы testing (UUID генерятся автоматически)
-- PostgreSQL
-- ============================================================

-- 1) Включаем генерацию UUID (если ещё не включено)
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- 2) Делаем автогенерацию UUID для id-колонок (если вдруг не настроено)
ALTER TABLE IF EXISTS testing.test_d        ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE IF EXISTS testing.question_d    ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE IF EXISTS testing.answer_d      ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE IF EXISTS testing.test_attempt_b ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE IF EXISTS testing.draft_b       ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE IF EXISTS testing.content_d     ALTER COLUMN id SET DEFAULT gen_random_uuid();

-- 3) Чистим таблицы перед заполнением (чтобы не было дублей и конфликтов)
TRUNCATE TABLE
  testing.content_d,
  testing.test_attempt_b,
  testing.answer_d,
  testing.question_d,
  testing.draft_b,
  testing.test_d
CASCADE;

-- 4) Заполняем тестовыми данными
DO $$
DECLARE
  -- tests
  t_alg UUID;
  t_geo UUID;

  -- questions (algebra)
  q1 UUID; q2 UUID; q3 UUID;
  -- questions (geometry)
  q4 UUID; q5 UUID; q6 UUID;

  -- correct answers (для attempt_version)
  q1_ok UUID; q2_ok UUID; q3_ok UUID;
  q4_ok UUID; q5_ok UUID; q6_ok UUID;

  -- any answer (чтобы content_d мог ссылаться, т.к. answer_id NOT NULL)
  q1_any UUID; q2_any UUID; q3_any UUID;
  q4_any UUID; q5_any UUID; q6_any UUID;

  -- students
  s1 UUID := gen_random_uuid();
  s2 UUID := gen_random_uuid();

BEGIN
  -- -----------------------------
  -- TESTS
  -- -----------------------------
  INSERT INTO testing.test_d (course_id, title, min_point, description)
  VALUES (gen_random_uuid(), 'Алгебра: линейные уравнения', 2, 'Тест на базовые линейные уравнения')
  RETURNING id INTO t_alg;

  INSERT INTO testing.test_d (course_id, title, min_point, description)
  VALUES (gen_random_uuid(), 'Геометрия: треугольники', 3, 'Тест на свойства треугольников')
  RETURNING id INTO t_geo;

  -- -----------------------------
  -- DRAFTS
  -- -----------------------------
  INSERT INTO testing.draft_b (title, min_point, description, test_id)
  VALUES ('Черновик: Алгебра (v0)', 2, 'Черновик перед публикацией', t_alg);

  INSERT INTO testing.draft_b (title, min_point, description, test_id)
  VALUES ('Черновик: Геометрия (v0)', 3, 'Черновик перед публикацией', t_geo);

  -- -----------------------------
  -- QUESTIONS (Алгебра)
  -- -----------------------------
  INSERT INTO testing.question_d (test_id, text_of_question, "order")
  VALUES (t_alg, 'Решите уравнение: 2x + 3 = 7', 1)
  RETURNING id INTO q1;

  INSERT INTO testing.question_d (test_id, text_of_question, "order")
  VALUES (t_alg, 'Решите уравнение: 5x = 20', 2)
  RETURNING id INTO q2;

  INSERT INTO testing.question_d (test_id, text_of_question, "order")
  VALUES (t_alg, 'Решите уравнение: x - 9 = -2', 3)
  RETURNING id INTO q3;

  -- ANSWERS for q1 (x=2)
  INSERT INTO testing.answer_d (text, question_id, score) VALUES
    ('x = 1', q1, 0),
    ('x = 2', q1, 1),
    ('x = 3', q1, 0),
    ('x = 4', q1, 0);

  SELECT id INTO q1_ok  FROM testing.answer_d WHERE question_id = q1 AND score = 1 LIMIT 1;
  SELECT id INTO q1_any FROM testing.answer_d WHERE question_id = q1 ORDER BY id LIMIT 1;

  -- ANSWERS for q2 (x=4)
  INSERT INTO testing.answer_d (text, question_id, score) VALUES
    ('x = 2', q2, 0),
    ('x = 3', q2, 0),
    ('x = 4', q2, 1),
    ('x = 5', q2, 0);

  SELECT id INTO q2_ok  FROM testing.answer_d WHERE question_id = q2 AND score = 1 LIMIT 1;
  SELECT id INTO q2_any FROM testing.answer_d WHERE question_id = q2 ORDER BY id LIMIT 1;

  -- ANSWERS for q3 (x=7)
  INSERT INTO testing.answer_d (text, question_id, score) VALUES
    ('x = 6', q3, 0),
    ('x = 7', q3, 1),
    ('x = 8', q3, 0),
    ('x = 9', q3, 0);

  SELECT id INTO q3_ok  FROM testing.answer_d WHERE question_id = q3 AND score = 1 LIMIT 1;
  SELECT id INTO q3_any FROM testing.answer_d WHERE question_id = q3 ORDER BY id LIMIT 1;

  -- -----------------------------
  -- QUESTIONS (Геометрия)
  -- -----------------------------
  INSERT INTO testing.question_d (test_id, text_of_question, "order")
  VALUES (t_geo, 'Сумма углов треугольника равна…', 1)
  RETURNING id INTO q4;

  INSERT INTO testing.question_d (test_id, text_of_question, "order")
  VALUES (t_geo, 'В равнобедренном треугольнике углы при основании…', 2)
  RETURNING id INTO q5;

  INSERT INTO testing.question_d (test_id, text_of_question, "order")
  VALUES (t_geo, 'Если два угла треугольника 30° и 60°, то третий угол равен…', 3)
  RETURNING id INTO q6;

  -- ANSWERS for q4 (180)
  INSERT INTO testing.answer_d (text, question_id, score) VALUES
    ('90°',  q4, 0),
    ('120°', q4, 0),
    ('180°', q4, 1),
    ('360°', q4, 0);

  SELECT id INTO q4_ok  FROM testing.answer_d WHERE question_id = q4 AND score = 1 LIMIT 1;
  SELECT id INTO q4_any FROM testing.answer_d WHERE question_id = q4 ORDER BY id LIMIT 1;

  -- ANSWERS for q5 (равны)
  INSERT INTO testing.answer_d (text, question_id, score) VALUES
    ('всегда прямые', q5, 0),
    ('равны',         q5, 1),
    ('сумма 90°',     q5, 0),
    ('не определены', q5, 0);

  SELECT id INTO q5_ok  FROM testing.answer_d WHERE question_id = q5 AND score = 1 LIMIT 1;
  SELECT id INTO q5_any FROM testing.answer_d WHERE question_id = q5 ORDER BY id LIMIT 1;

  -- ANSWERS for q6 (90)
  INSERT INTO testing.answer_d (text, question_id, score) VALUES
    ('60°',  q6, 0),
    ('80°',  q6, 0),
    ('90°',  q6, 1),
    ('120°', q6, 0);

  SELECT id INTO q6_ok  FROM testing.answer_d WHERE question_id = q6 AND score = 1 LIMIT 1;
  SELECT id INTO q6_any FROM testing.answer_d WHERE question_id = q6 ORDER BY id LIMIT 1;

  -- -----------------------------
  -- CONTENT (у тебя и question_id и answer_id NOT NULL)
  -- type_of_content: TRUE = "контент вопроса", FALSE = "контент ответа" (пример)
  -- -----------------------------
  INSERT INTO testing.content_d ("order", content, type_of_content, question_id, answer_id) VALUES
    (1, 'Подсказка: перенеси число вправо и приведи подобные.', TRUE,  q1, q1_any),
    (1, 'Теория: сумма внутренних углов треугольника — 180°.',  TRUE,  q4, q4_any),
    (2, 'Пояснение: правильный ответ даёт 1 балл.',             FALSE, q1, q1_ok);

  -- -----------------------------
  -- TEST ATTEMPTS
  -- attempt_version заполним JSON-ом с выбранными ответами
  -- -----------------------------
  INSERT INTO testing.test_attempt_b (student_id, test_id, date_of_attempt, point, certificate_id, attempt_version)
  VALUES
    (s1, t_alg, DATE '2025-12-18', 3, NULL,
      json_build_object('answers', json_build_array(
        json_build_object('question', q1, 'answer', q1_ok),
        json_build_object('question', q2, 'answer', q2_ok),
        json_build_object('question', q3, 'answer', q3_ok)
      ))
    ),
    (s2, t_geo, DATE '2025-12-18', 2, NULL,
      json_build_object('answers', json_build_array(
        json_build_object('question', q4, 'answer', q4_ok),
        json_build_object('question', q5, 'answer', q5_ok),
        json_build_object('question', q6, 'answer', q6_any) -- намеренно ошибочный/любой
      ))
    );

END $$;

-- Быстрые проверки
-- SELECT * FROM testing.test_d;
-- SELECT * FROM testing.draft_b;
-- SELECT * FROM testing.question_d ORDER BY test_id, "order";
-- SELECT * FROM testing.answer_d;
-- SELECT * FROM testing.test_attempt_b;
-- SELECT * FROM testing.content_d;