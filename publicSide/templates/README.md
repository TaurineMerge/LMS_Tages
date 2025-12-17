# Фронтенд для publicSide

## Описание

Добавлен фронтенд на Handlebars для отображения курсов по категориям. Страница отображает список курсов с карточками, включающими:
- Название курса
- Уровень сложности (легкий, средний, сложный)
- Описание курса (обрезанное до 150 символов)
- Статистику (рейтинг, количество студентов, длительность)

## Структура проекта

```
publicSide/
├── views/                      # Handlebars шаблоны
│   ├── layouts/               # Layout шаблоны
│   │   └── main.hbs          # Основной layout
│   ├── partials/             # Переиспользуемые компоненты
│   │   └── course-card.hbs   # Карточка курса
│   └── pages/                # Страницы
│       └── courses.hbs       # Страница списка курсов
├── static/                    # Статические файлы
│   └── css/
│       └── styles.css        # Стили для фронтенда
├── internal/
│   ├── handler/
│   │   └── course.go         # Handler для отображения курсов
│   ├── service/
│   │   └── course.go         # Сервис для работы с курсами
│   ├── repository/
│   │   └── course.go         # Репозиторий для курсов и категорий
│   └── template/
│       └── engine.go         # Конфигурация Handlebars engine
└── cmd/
    └── main.go               # Точка входа (обновлен)
```

## Новые эндпоинты

### HTML страницы
- `GET /categories/:categoryId/courses` - Страница со списком курсов по категории

### API (существующие)
- `GET /api/v1/categories/:categoryId/courses/:courseId/lessons` - Список уроков курса
- `GET /api/v1/categories/:categoryId/courses/:courseId/lessons/:lessonId` - Детали урока

## Использование

### Запуск проекта

1. Убедитесь, что база данных запущена и содержит данные:
```bash
docker-compose up -d postgres
```

2. Запустите publicSide сервис:
```bash
cd publicSide
go run cmd/main.go
```

3. Откройте браузер и перейдите на страницу курсов:
```
http://localhost:8081/categories/{category_id}/courses
```

Замените `{category_id}` на реальный UUID категории из базы данных.

### Получение category_id

Вы можете получить список категорий через SQL запрос:
```sql
SELECT id, title FROM knowledge_base.category_d;
```

Или добавить API эндпоинт для получения списка категорий.

## Особенности реализации

### Дизайн
- Современный дизайн с карточками курсов
- Адаптивная верстка (responsive design)
- Цветовая индикация уровня сложности:
  - Зеленый - легкий уровень
  - Желтый - средний уровень
  - Красный - сложный уровень

### Безопасность
- Только публичные курсы (`visibility = 'public'`) отображаются на странице
- Валидация UUID параметров
- Защита от SQL инъекций через параметризованные запросы

### Производительность
- Использование connection pool для базы данных
- Трейсинг запросов через OpenTelemetry
- Оптимизированные SQL запросы

## Customization

### Изменение стилей
Отредактируйте файл `static/css/styles.css` для изменения внешнего вида.

### Добавление новых полей в карточку курса
1. Обновите структуру `CourseView` в `internal/handler/course.go`
2. Обновите шаблон `views/pages/courses.hbs`

### Добавление изображений курсов
Добавьте поле `image_url` в таблицу `knowledge_base.course_b` и обновите:
1. Domain модель `internal/domain/course.go`
2. Repository `internal/repository/course.go`
3. Handler `internal/handler/course.go`
4. Шаблон `views/pages/courses.hbs`

## Зависимости

Новые зависимости:
- `github.com/gofiber/template/handlebars/v2` - Handlebars template engine для Fiber

## Troubleshooting

### Ошибка "Category not found"
Убедитесь, что category_id существует в базе данных.

### Курсы не отображаются
Проверьте, что в категории есть курсы со статусом `visibility = 'public'`.

### Стили не применяются
Убедитесь, что статические файлы доступны по пути `/static/css/styles.css` и сервер запущен корректно.
