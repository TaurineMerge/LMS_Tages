---
title: "Flowchart-диаграмма прохождения теста"
date: 2025-12-15
layout: "single"
---

## Прохождение теста учеником

{{< mermaid >}}
flowchart TD
    STU_Start[Ученик завершил курс] --> STU_Click[Нажатие 'Пройти тест'];
    STU_Click --> STU_Redirect[Перенаправление в модуль тестирования];
    STU_Redirect --> STU_Load[Загрузка страницы теста<br>все вопросы на одной странице];
    
    STU_Load --> STU_Select{Выбор вопроса};
    STU_Select --> STU_View[Просмотр вопроса и вариантов];
    STU_View --> STU_Choose[Выбор ответа/ов];
    STU_Choose --> STU_Save[Нажатие 'Ответить'];
    
    STU_Save --> STU_SendToServer[Отправка ответа на сервер];
    STU_SendToServer --> STU_Confirmation[Подтверждение от сервера];
    STU_Confirmation --> STU_Saved[Ответ сохранен<br>редактирование заблокировано];
    
    STU_Saved --> STU_Check{Все вопросы отвечены?};
    
    STU_Check -->|Нет| STU_Continue{Продолжить отвечать?};
    STU_Continue -->|Да| STU_Select;
    STU_Continue -->|Нет| STU_Finish[Нажатие 'Завершить тест'];
    
    STU_Check -->|Да| STU_Finish;
    
    STU_Finish --> STU_Confirm{Все вопросы отвечены?};
    STU_Confirm -->|Нет| STU_Warning[Диалог подтверждения<br>'Уверены, что хотите закончить?'];
    STU_Warning --> STU_Decision{Подтверждение?};
    STU_Decision -->|Да| STU_Finalize;
    STU_Decision -->|Нет| STU_Select;
    
    STU_Confirm -->|Да| STU_Finalize[Отправка команды<br>'Завершить попытку'];
    
    STU_Finalize --> STU_Results[Перенаправление на страницу результатов];
    STU_Results --> STU_Show[Отображение результатов<br> на основе данных с сервера];
    STU_Show --> STU_PassCheck{Успешная сдача?<br>балл ≥ проходного};
    
    STU_PassCheck -->|Да| STU_Cert[Автоматическая генерация сертификата];
    STU_Cert --> STU_Send[Сертификат отправлен в ЛК];
    STU_Send --> STU_Congrats[Поздравление с успешным прохождением];
    
    STU_PassCheck -->|Нет| STU_Fail[Отображение сообщения о провале];
    STU_Fail --> STU_Retry[Кнопка 'Пройти еще раз'];
    STU_Retry --> STU_Return[Возврат к изучению материалов];

    style STU_Click fill:#e3f2fd
    style STU_Save fill:#fff3e0
    style STU_SendToServer fill:#e1f5fe
    style STU_Finish fill:#ffebee
    style STU_Finalize fill:#f3e5f5
    style STU_Cert fill:#e8f5e8
    style STU_Retry fill:#fce4ec
{{< /mermaid >}}

---

## Создание и редактирование теста учителем

{{< mermaid >}}
flowchart TD
    A[Админ в редакторе курса] --> B{Действие с тестом};
    B --> C[Создать новый тест];
    B --> D[Редактировать существующий];
    
    C --> E[Загрузка конструктора];
    D --> F{Тип теста};
    F --> G[Черновик];
    F --> H[Опубликованный тест];
    G --> E;
    H --> I[Проверка активных попыток];
    I --> J{Есть активные попытки?};
    J -->|Да| K[Предупреждение,<br>ограниченное редактирование];
    J -->|Нет| L[Полный доступ];
    K --> E;
    L --> E;
    
    E --> M[Интерфейс конструктора];
    M --> N{Добавить вопрос?};
    N -->|Да| O[Выбор типа вопроса];
    O --> P[Один ответ];
    O --> Q[Несколько ответов];
    P --> R[Формулировка вопроса];
    Q --> R;
    R --> S[Добавление вариантов];
    S --> T[Назначение баллов];
    T --> U[Авторасчет макс. баллов<br>за вопрос];
    U --> V[Изменение порядка];
    V --> W{Еще вопросы?};
    W -->|Да| N;
    W -->|Нет| X;
    
    X[Все вопросы готовы] --> Y[Авторасчет общего максимума];
    Y --> Z[Установка проходного балла];
    Z --> AA{Состояние теста};
    AA --> BB[Сохранить как черновик];
    AA --> CC[Опубликовать/обновить];
    
    BB --> DD[Тест - черновик];
    CC --> EE[Тест опубликован];
    DD --> FF[Возврат в редактор курса];
    EE --> FF;
    
    FF --> GG{Дальнейшие действия};
    GG --> HH[Продолжить курс];
    GG --> II[Выйти];
    GG --> JJ[Редактировать тест снова];
    JJ --> E;
    
    style C fill:#e3f2fd
    style D fill:#f3e5f5
    style BB fill:#fff3e0
    style CC fill:#e8f5e8
    style EE fill:#c8e6c9
    style G fill:#e0f7fa
    style H fill:#f1f8e9
    style K fill:#ffebee
    style L fill:#e8f5e8
    style JJ fill:#f3e5f5
{{< /mermaid >}}