package com.example.lms.test.domain.model;

import java.util.Objects;
import java.util.UUID;

/**
 * Domain Model: TestModel
 *
 * Представляет тест курса и соответствует строке таблицы test_d.
 *
 * Структура таблицы test_d:
 *  - id          UUID    (PK, not null)
 *  - course_id   UUID    (FK → course_d.id) — может быть null
 *  - title       VARCHAR — название теста
 *  - min_point   INT     — минимальный балл для прохождения (может быть null)
 *  - description TEXT     — описание теста (может быть null)
 *
 * Доменные правила:
 *  - title и courseId обязательны (доменные ограничения сильнее, чем ограничения БД)
 *  - minPoint >= 0, если указано
 *  - description и minPoint могут быть null
 */
public class TestModel {

    /** Уникальный ID теста (PK). */
    private UUID id;

    /** ID курса, которому принадлежит тест. */
    private UUID courseId;

    /** Название теста. */
    private String title;

    /** Минимальное количество баллов для прохождения (может быть null). */
    private Integer minPoint;

    /** Описание теста (может быть null). */
    private String description;

    // ---------------------- КОНСТРУКТОРЫ ----------------------

    /**
     * Конструктор для создания нового теста.
     */
    public TestModel(UUID courseId, String title, Integer minPoint, String description) {
        this.courseId = courseId;
        this.title = Objects.requireNonNull(title, "Title cannot be null");
        this.minPoint = minPoint;       // допускаем null
        this.description = description; // допускаем null
    }

    /**
     * Конструктор для загрузки теста из базы данных.
     */
    public TestModel(UUID id, UUID courseId, String title, Integer minPoint, String description) {
        this.id = id;
        this.courseId = courseId;
        this.title = title;
        this.minPoint = minPoint;
        this.description = description;
    }

    // ---------------------- GETTERS ----------------------

    public UUID getId() { return id; }
    public UUID getCourseId() { return courseId; }
    public String getTitle() { return title; }
    public Integer getMinPoint() { return minPoint; }
    public String getDescription() { return description; }

    // ---------------------- SETTERS С ВАЛИДАЦИЕЙ ----------------------

    public void setId(UUID id) {
        this.id = Objects.requireNonNull(id, "ID cannot be null");
    }

    public void setCourseId(UUID courseId) {
        this.courseId = Objects.requireNonNull(courseId, "Course ID cannot be null");
    }

    public void setTitle(String title) {
        if (title == null || title.trim().isEmpty()) {
            throw new IllegalArgumentException("Title cannot be null or empty");
        }
        this.title = title;
    }

    public void setMinPoint(Integer minPoint) {
        if (minPoint != null && minPoint < 0) {
            throw new IllegalArgumentException("minPoint cannot be negative");
        }
        this.minPoint = minPoint;
    }

    public void setDescription(String description) {
        this.description = description; // null OK
    }

    // ---------------------- ДОМЕННАЯ ЛОГИКА ----------------------

    /**
     * Проверяет, прошёл ли студент тест.
     * Если minPoint = null → считается, что пройти можно всегда.
     */
    public boolean isPassed(int studentScore) {
        return minPoint == null || studentScore >= minPoint;
    }

    /**
     * Проверка валидности теста перед сохранением.
     */
    public void validate() {
        if (courseId == null) {
            throw new IllegalArgumentException("Course ID cannot be null");
        }
        if (title == null || title.trim().isEmpty()) {
            throw new IllegalArgumentException("Title cannot be empty");
        }
        if (minPoint != null && minPoint < 0) {
            throw new IllegalArgumentException("Minimum points cannot be negative");
        }
    }

    // ---------------------- UTILS ----------------------

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        TestModel testModel = (TestModel) o;
        return Objects.equals(id, testModel.id);
    }

    @Override
    public int hashCode() {
        return Objects.hash(id);
    }

    /**
     * Укороченное строковое представление теста.
     */
    @Override
    public String toString() {
        return "TestModel{" +
                "id=" + id +
                ", title='" + title + '\'' +
                ", courseId=" + courseId +
                ", minPoint=" + minPoint +
                '}';
    }
}