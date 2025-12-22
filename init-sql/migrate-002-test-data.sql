-- =============================================
-- ТЕСТОВЫЕ ДАННЫЕ ДЛЯ РАЗРАБОТКИ
-- =============================================

-- 1. Студент (соответствует пользователю в Keycloak student.json)
INSERT INTO personal_account.student_s (id, name, surname, email)
VALUES ('d3b07384-d92d-4e12-9096-c8350fdc5fd2', 'student', 'main', 'student@example.com')
ON CONFLICT (email) DO UPDATE SET 
    name = EXCLUDED.name,
    surname = EXCLUDED.surname;

-- 2. Категория и Курс
INSERT INTO knowledge_base.category_d (id, title)
VALUES ('a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d', 'Программирование')
ON CONFLICT (id) DO NOTHING;

INSERT INTO knowledge_base.course_b (id, title, description, level, category_id, visibility)
VALUES ('b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'Python для начинающих', 'Основы языка Python', 'easy', 'a1b2c3d4-e5f6-4a5b-8c9d-0e1f2a3b4c5d', 'public')
ON CONFLICT (id) DO NOTHING;

-- 3. Тест
INSERT INTO tests.test_d (id, course_id, title, min_point, description)
VALUES ('c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', 'b2c3d4e5-f6a7-4b8c-9d0e-1f2a3b4c5d6e', 'Финальный тест по Python', 70, 'Проверка знаний основ Python')
ON CONFLICT (id) DO NOTHING;

-- 4. Попытки (Attempts)
INSERT INTO tests.test_attempt_b (id, student_id, test_id, date_of_attempt, point, result, passed, completed)
VALUES 
('d4e5f6a7-b8c9-4d0e-1f2a-3b4c5d6e7f8a', 'd3b07384-d92d-4e12-9096-c8350fdc5fd2', 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', CURRENT_DATE - INTERVAL '2 days', 45, '{"status": "failed"}'::jsonb, false, true),
('e5f6a7b8-c9d0-4e1f-2a3b-4c5d6e7f8a9b', 'd3b07384-d92d-4e12-9096-c8350fdc5fd2', 'c3d4e5f6-a7b8-4c9d-0e1f-2a3b4c5d6e7f', CURRENT_DATE - INTERVAL '1 day', 85, '{"status": "passed"}'::jsonb, true, true)
ON CONFLICT (id) DO NOTHING;
