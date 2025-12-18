package com.example.lms.draft.domain.model;

import java.util.Objects;
import java.util.UUID;

/**
 * Domain Model: DraftModel
 *
 * Представляет черновик теста и соответствует строке таблицы draft_b.
 *
 * Структура таблицы draft_b:
 *  - id          UUID    (PK, not null, unique)
 *  - title       VARCHAR — название (может быть null на уровне БД, но в домене считаем обязательным)
 *  - min_point   INT     — минимальный балл для прохождения (может быть null)
 *  - description TEXT     — описание (может быть null)
 *  - test_id     UUID    (FK/ссылка на test_d, not null)
 *
 * Доменные правила:
 *  - testId обязателен
 *  - title обязателен и не пустой
 *  - minPoint >= 0, если указан
 *  - description может быть null
 */
public class DraftModel {

    /** Уникальный ID черновика (PK). */
    private UUID id;

    /** ID теста, к которому относится черновик (обязателен). */
    private UUID testId;

    /** Название черновика. */
    private String title;

    /** Минимальное количество баллов для прохождения (может быть null). */
    private Integer minPoint;

    /** Описание черновика (может быть null). */
    private String description;

    // ---------------------- КОНСТРУКТОРЫ ----------------------

    /**
     * Конструктор для создания нового черновика.
     */
    public DraftModel(UUID testId, String title, Integer minPoint, String description) {
        this.testId = Objects.requireNonNull(testId, "Test ID cannot be null");
        this.title = Objects.requireNonNull(title, "Title cannot be null");
        this.minPoint = minPoint;       // допускаем null
        this.description = description; // допускаем null
    }

    /**
     * Конструктор для загрузки черновика из базы данных.
     */
    public DraftModel(UUID id, UUID testId, String title, Integer minPoint, String description) {
        this.id = id;
        this.testId = testId;
        this.title = title;
        this.minPoint = minPoint;
        this.description = description;
    }

    // ---------------------- GETTERS ----------------------

    public UUID getId() { return id; }
    public UUID getTestId() { return testId; }
    public String getTitle() { return title; }
    public Integer getMinPoint() { return minPoint; }
    public String getDescription() { return description; }

    // ---------------------- SETTERS С ВАЛИДАЦИЕЙ ----------------------

    public void setId(UUID id) {
        this.id = Objects.requireNonNull(id, "ID cannot be null");
    }

    public void setTestId(UUID testId) {
        this.testId = Objects.requireNonNull(testId, "Test ID cannot be null");
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
     * Проверяет, прошёл ли студент тест по этому черновику.
     * Если minPoint = null → считается, что пройти можно всегда.
     */
    public boolean isPassed(int studentScore) {
        return minPoint == null || studentScore >= minPoint;
    }

    /**
     * Проверка валидности черновика перед сохранением.
     */
    public void validate() {
        if (testId == null) {
            throw new IllegalArgumentException("Test ID cannot be null");
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
        DraftModel draftModel = (DraftModel) o;
        return Objects.equals(id, draftModel.id);
    }

    @Override
    public int hashCode() {
        return Objects.hash(id);
    }

    /**
     * Укороченное строковое представление черновика.
     */
    @Override
    public String toString() {
        return "DraftModel{" +
                "id=" + id +
                ", title='" + title + '\'' +
                ", testId=" + testId +
                ", minPoint=" + minPoint +
                '}';
    }
}