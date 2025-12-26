---
title: "Flowchart-диаграмма управления контентом"
date: 2025-12-15
layout: "single"
---

## Управление контентом (админ)

{{< mermaid >}}

flowchart TD
    Start[Вход в систему] --> Auth[Авторизация через Keycloak]
    Auth -->|Успех| CheckRole{Проверка роли}
    
    CheckRole -->|Админ| AdminPanel[Админ-панель]
    CheckRole -->|Студент| StudentPanel[Личный кабинет]
    CheckRole -->|Другая роль| ErrorRole[Доступ запрещен]
    
    AdminPanel --> Dashboard[Дашборд]
    Dashboard --> Menu{Выбор действия}
    
    Menu -->|Управление категориями| Categories[Категории]
    Menu -->|Управление курсами| Courses[Курсы]
    Menu -->|Управление уроками| Lessons[Уроки]
    
    Categories --> CategoryAction{Действие с категорией}
    CategoryAction -->|Создать| CreateCategory[Создать категорию]
    CategoryAction -->|Редактировать| EditCategory[Редактировать категорию]
    CategoryAction -->|Удалить| DeleteCategory[Удалить категорию]
    
    Courses --> CourseAction{Действие с курсом}
    CourseAction -->|Создать| CreateCourse[Создать курс]
    CourseAction -->|Редактировать| EditCourse[Редактировать курс]
    CourseAction -->|Удалить| DeleteCourse[Удалить курс]
    CourseAction -->|Создать тест| TestBuilder[Конструктор тестов<br/><i>внешний сервис</i>]
    
    Lessons --> LessonAction{Действие с уроком}
    LessonAction -->|Создать| CreateLesson[Создать урок]
    LessonAction -->|Редактировать| EditLesson[Редактировать урок]
    LessonAction -->|Удалить| DeleteLesson[Удалить урок]
    
    CreateCategory --> SaveCategory[Сохранить в БД<br/>knowledge_base_db]
    EditCategory --> UpdateCategory[Обновить в БД]
    DeleteCategory --> RemoveCategory[Удалить из БД]
    
    CreateCourse --> UploadImage[Загрузить изображение в S3]
    UploadImage --> SaveCourse[Сохранить курс в БД<br/>с ссылкой на S3]
    EditCourse --> UpdateCourse[Обновить курс в БД]
    DeleteCourse --> RemoveCourse[Удалить курс из БД]
    
    CreateLesson --> SaveLesson[Сохранить урок в БД]
    EditLesson --> UpdateLesson[Обновить урок в БД]
    DeleteLesson --> RemoveLesson[Удалить урок из БД]
    
    SaveCategory --> SuccessCat[Категория сохранена]
    UpdateCategory --> SuccessCat
    RemoveCategory --> SuccessDelCat[Категория удалена]
    
    SaveCourse --> SuccessCourse[Курс сохранен]
    UpdateCourse --> SuccessCourse
    RemoveCourse --> SuccessDelCourse[Курс удален]
    
    SaveLesson --> SuccessLesson[Урок сохранен]
    UpdateLesson --> SuccessLesson
    RemoveLesson --> SuccessDelLesson[Урок удален]
    
    SuccessCat --> Dashboard
    SuccessDelCat --> Dashboard
    SuccessCourse --> Dashboard
    SuccessDelCourse --> Dashboard
    SuccessLesson --> Dashboard
    SuccessDelLesson --> Dashboard
    
    TestBuilder -.->|Возврат| Courses

{{< /mermaid >}}