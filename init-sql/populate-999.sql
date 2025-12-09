DO $$
DECLARE
    -- Массивы с названиями для генерации
    category_titles TEXT[] := ARRAY['Programming', 'Data Science', 'Design', 'Marketing', 'Business', 'Languages'];
    course_prefixes TEXT[] := ARRAY['Introduction to', 'Advanced', 'Mastering', 'Essentials of', 'Complete Guide to'];
    course_subjects TEXT[] := ARRAY['Go', 'Python', 'JavaScript', 'SQL', 'Docker', 'React', 'Figma', 'SEO', 'Agile', 'German'];
    lesson_prefixes TEXT[] := ARRAY['Chapter', 'Module', 'Unit', 'Part', 'Section'];
    
    -- ID для связывания
    v_category_id UUID;
    v_course_id UUID;
    v_lesson_id UUID;
    
    -- JSON массив для контента урока
    v_lesson_content JSONB;
    v_content_item JSONB;
    
    -- Счетчики и случайные значения
    i INT;
    j INT;
    k INT;
    l INT;
    num_courses INT;
    num_lessons INT;
    num_contents INT;
BEGIN
    RAISE NOTICE 'Starting to populate database with test data...';

    -- 1. Создание категорий
    FOR i IN 1..array_length(category_titles, 1) LOOP
        INSERT INTO knowledge_base.category_d (title) 
        VALUES (category_titles[i]) 
        RETURNING id INTO v_category_id;

        -- 2. Создание курсов для каждой категории
        num_courses := floor(random() * 11) + 5; -- от 5 до 15 курсов
        FOR j IN 1..num_courses LOOP
            INSERT INTO knowledge_base.course_b (title, description, category_id, level, visibility) 
            VALUES (
                course_prefixes[floor(random() * array_length(course_prefixes, 1)) + 1] || ' ' || course_subjects[floor(random() * array_length(course_subjects, 1)) + 1] || ' ' || j,
                'A comprehensive course on ' || course_subjects[floor(random() * array_length(course_subjects, 1)) + 1],
                v_category_id,
                CASE floor(random() * 3)
                    WHEN 0 THEN 'easy'
                    WHEN 1 THEN 'medium'
                    ELSE 'hard'
                END,
                'public'
            ) RETURNING id INTO v_course_id;

            -- 3. Создание уроков для каждого курса
            num_lessons := floor(random() * 16) + 5; -- от 5 до 20 уроков
            FOR k IN 1..num_lessons LOOP
                
                -- 4. Создание контента для каждого урока
                v_lesson_content := '[]'::jsonb;
                num_contents := floor(random() * 21) + 10; -- от 10 до 30 элементов контента
                FOR l IN 1..num_contents LOOP
                    v_content_item := jsonb_build_object(
                        'type', 'text',
                        'value', 'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed non risus. Suspendisse lectus tortor, dignissim sit amet, adipiscing nec, ultricies sed, dolor. Cras elementum ultrices diam. Maecenas ligula massa, varius a, semper congue, euismod non, mi. Proin porttitor, orci nec nonummy molestie, enim est eleifend mi, non fermentum diam nisl sit amet erat.'
                    );
                    v_lesson_content := v_lesson_content || jsonb_build_array(v_content_item);
                END LOOP;

                INSERT INTO knowledge_base.lesson_d (title, course_id, content) 
                VALUES (
                    lesson_prefixes[floor(random() * array_length(lesson_prefixes, 1)) + 1] || ' ' || k || ': Getting Started',
                    v_course_id,
                    v_lesson_content
                );
            END LOOP;
        END LOOP;
    END LOOP;

    RAISE NOTICE 'Finished populating database.';
END $$;
