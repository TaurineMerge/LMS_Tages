---
title: "Flowchart-диаграмма действий гостя на публичном сайте"
date: 2025-12-15
layout: "single"
---

{{< mermaid >}}

flowchart TD
    Start([Гость заходит на сайт]) --> HomePage[Главная страница]
    
    %% Шапка сайта - навигация
    HomePage --> Header[Header с навигацией:<br/>Главная, Категории, Вход, Регистрация]
    
    %% Главная страница
    HomePage --> Welcome[Видит инструкцию по работе с сайтом]
    Welcome --> FeaturedCategories[Видит несколько категорий<br/>с курсами на главной]
    
    %% Выбор действий на главной
    FeaturedCategories --> HomeAction{Выбор действия на главной}
    HomeAction -->|Клик по категории| CategorySinglePage[Страница конкретной категории<br/>со всеми курсами]
    HomeAction -->|Клик по курсу| CourseSinglePage[Страница конкретного курса]
    HomeAction -->|Продолжить просмотр| ContinueHome
    
    ContinueHome --> HeaderNav{Выбор в хедере}
    
    %% Навигация в хедере
    Header --> HeaderNav
    
    HeaderNav -->|Категории| AllCategoriesPage[Страница 'Категории']
    HeaderNav -->|Вход| LoginPage[Страница входа<br/>через Keycloak]
    HeaderNav -->|Регистрация| RegisterPage[Страница регистрации<br/>через Keycloak]
    HeaderNav -->|Главная| HomePage
    
    %% Страница всех категорий
    AllCategoriesPage --> ViewAllCategories[Видит все категории<br/>Под каждой категорией до 5 курсов]
    
    ViewAllCategories --> CheckCourseCount{Курсов в категории > 5?}
    
    CheckCourseCount -->|≤ 5 курсов| ViewCourseFromCategory[Может просмотреть любой курс]
    CheckCourseCount -->|> 5 курсов| ShowMoreButton[Видит кнопку<br/>'Посмотреть еще']
    
    ShowMoreButton -->|Клик| CategorySinglePage
    ViewCourseFromCategory -->|Клик по курсу| CourseSinglePage
    
    %% Страница конкретной категории
    CategorySinglePage --> ViewAllCategoryCourses[Видит все курсы категории:<br/>- Название<br/>- Описание<br/>- Картинка]
    ViewAllCategoryCourses -->|Клик по курсу| CourseSinglePage
    
    %% Страница курса
    CourseSinglePage --> ViewCourseDetails[Видит курс:<br/>- Название<br/>- Описание<br/>- Картинка<br/>- Уровень сложности]
    ViewCourseDetails --> ViewLessons[Видит список уроков курса]
    
    ViewLessons -->|Клик по уроку| LessonPage[Страница урока]
    
    %% Страница урока
    LessonPage --> ViewLessonContent[Просматривает содержание урока]
    ViewLessonContent --> LessonNavigation{Навигация}
    
    LessonNavigation -->|← Вернуться к курсу| CourseSinglePage
    LessonNavigation -->|Следующий урок| NextLesson[Страница следующего урока]
    LessonNavigation -->|Предыдущий урок| PrevLesson[Страница предыдущего урока]
    
    NextLesson --> ViewLessonContent
    PrevLesson --> ViewLessonContent
    
    %% Кнопка теста (видна всегда, но активна только для авторизованных)
    CourseSinglePage --> TestButtonSection[Раздел с кнопкой 'Пройти тест']
    
    TestButtonSection --> AuthStatus{Пользователь авторизован?}
    
    AuthStatus -->|Нет| GuestTestButton[Кнопка 'Пройти тест'<br/>НЕАКТИВНА]
    
    AuthStatus -->|Да| AuthUserTestButton[Кнопка 'Пройти тест'<br/>АКТИВНА]
    
    GuestTestButton --> GuestAction{Действие гостя}
    GuestAction -->|Клик по неактивной кнопке| NoAction[Ничего не происходит]
    GuestAction -->|Решает авторизоваться| GoToLogin[Переход на страницу входа]
    
    GoToLogin --> LoginPage
    
    AuthUserTestButton -->|Клик| StartTest[Начать тест]
    StartTest --> ExternalTestingService[Внешний сервис тестирования]
    
    %% Процесс авторизации/регистрации
    LoginPage --> KeycloakAuth[Авторизация через Keycloak]
    RegisterPage --> KeycloakRegister[Регистрация через Keycloak]
    
    KeycloakAuth -->|Успех| ReturnToPage[Возврат на ту же страницу<br/>где был пользователь]
    KeycloakRegister -->|Успех| ReturnToPage
    
    ReturnToPage -->|Был на странице курса| CourseSinglePageWithAuth[Страница курса<br/>теперь с активной кнопкой теста]
    ReturnToPage -->|Был на другой странице| OriginalPage[Исходная страница]
    
    CourseSinglePageWithAuth --> AuthUserTestButton
    
    %% Стили для наглядности
    style Start fill:#f0f0f0,stroke:#333,stroke-width:2px
    style HomePage fill:#e3f2fd,stroke:#1565c0
    style CourseSinglePage fill:#f3e5f5,stroke:#7b1fa2
    style AuthUserTestButton fill:#c8e6c9,stroke:#2e7d32
    style GuestTestButton fill:#ffecb3,stroke:#ff8f00
    style ExternalTestingService fill:#f5f5f5,stroke:#616161,stroke-dasharray:5 5
    style KeycloakAuth fill:#ffecb3,stroke:#ff8f00

{{< /mermaid >}}