# WYSIWYG Редактор для уроков

## Описание
Добавлен полноценный WYSIWYG редактор для редактирования HTML контента уроков в админ-панели.

## Функционал

### Форматирование текста:
- **Жирный** (Ctrl+B)
- *Курсив* (Ctrl+I)
- <u>Подчеркнутый</u> (Ctrl+U)
- ~~Зачеркнутый~~

### Заголовки:
- H1, H2, H3
- Параграфы

### Списки:
- Маркированные списки
- Нумерованные списки

### Другое:
- Ссылки
- Изображения (по URL)
- Блоки кода
- Цитаты
- Выравнивание текста (слева, по центру, справа)
- Отмена/Повтор действий
- Очистка форматирования

## Миграция базы данных

Для работы редактора необходимо применить миграцию:

### Автоматически (при первом запуске):
Если вы запускаете базу данных впервые, миграция применится автоматически из файла:
```
init-sql/knowledge-base-db/04-add-html-content-to-lessons.sql
```

### Вручную (если база уже запущена):

1. Подключитесь к базе данных:
```bash
docker exec -it app-db psql -U postgres -d knowledge_base_db
```

2. Выполните миграцию:
```sql
ALTER TABLE knowledge_base.lesson_d 
ADD COLUMN IF NOT EXISTS html_content TEXT;

COMMENT ON COLUMN knowledge_base.lesson_d.html_content IS 'HTML контент урока из WYSIWYG редактора';
```

Или:
```bash
docker exec -i app-db psql -U postgres -d knowledge_base_db < init-sql/knowledge-base-db/04-add-html-content-to-lessons.sql
```

## Файлы изменений

### Backend (Go):
- `models/lesson.go` - добавлено поле `HTMLContent`
- `handlers/dto/request/lesson.go` - добавлено поле в DTO
- `handlers/web/lesson_web_handler.go` - обработка html_content в формах
- `repositories/lesson.go` - SQL запросы с полем html_content

### Frontend:
- `static/js/wysiwyg-editor.js` - JavaScript код редактора
- `static/css/wysiwyg-editor.css` - Стили редактора
- `templates/pages/lesson-form.hbs` - Форма с редактором
- `templates/layouts/main.hbs` - Подключение CSS и JS

### База данных:
- `init-sql/knowledge-base-db/04-add-html-content-to-lessons.sql` - Миграция

## Использование

1. Перейдите в админ-панель
2. Откройте редактирование урока или создайте новый
3. Используйте панель инструментов для форматирования текста
4. HTML контент сохраняется автоматически при отправке формы

## Технические детали

- Редактор использует встроенный браузерный API `contentEditable`
- HTML сохраняется в поле `html_content` в таблице `lesson_d`
- Поддержка вставки текста без форматирования (только plain text)
- Responsive дизайн для мобильных устройств
