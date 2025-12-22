-- Добавление поля html_content для WYSIWYG редактора в таблицу lesson_d
-- Дата: 2025-12-22

-- Добавляем новую колонку для хранения HTML контента из WYSIWYG редактора
ALTER TABLE knowledge_base.lesson_d 
ADD COLUMN IF NOT EXISTS html_content TEXT;

-- Комментарий к колонке
COMMENT ON COLUMN knowledge_base.lesson_d.html_content IS 'HTML контент урока из WYSIWYG редактора';

-- Опционально: можно добавить индекс для полнотекстового поиска по html_content
-- CREATE INDEX IF NOT EXISTS idx_lesson_html_content_fts 
-- ON knowledge_base.lesson_d USING gin(to_tsvector('russian', html_content));
