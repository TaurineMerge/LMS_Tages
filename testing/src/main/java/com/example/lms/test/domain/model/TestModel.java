package com.example.lms.test.domain.model;

import java.util.Objects;
import java.util.UUID;

/**
 * Модель теста - представляет строку таблицы test_d
 * Содержит только данные и базовую бизнес-логику
 */
public class TestModel {
    private UUID id;
    private UUID courseId;
    private String title;
    private Integer minPoint;
    private String description;
    
    // Конструктор для создания нового теста
    public TestModel(UUID courseId, String title, Integer minPoint, String description) {
        this.courseId = Objects.requireNonNull(courseId, "Course ID cannot be null");
        this.title = Objects.requireNonNull(title, "Title cannot be null");
        this.minPoint = minPoint;
        this.description = description;
    }
    
    // Конструктор для загрузки из БД
    public TestModel(UUID id, UUID courseId, String title, Integer minPoint, String description) {
        this.id = id;
        this.courseId = courseId;
        this.title = title;
        this.minPoint = minPoint;
        this.description = description;
    }
    
    // Геттеры
    public UUID getId() { return id; }
    public UUID getCourseId() { return courseId; }
    public String getTitle() { return title; }
    public Integer getMinPoint() { return minPoint; }
    public String getDescription() { return description; }
    
    // Сеттеры с валидацией
    public void setId(UUID id) { 
        this.id = Objects.requireNonNull(id, "ID cannot be null");
    }
    
    public void setTitle(String title) {
        this.title = Objects.requireNonNull(title, "Title cannot be null");
    }
    
    public void setMinPoint(Integer minPoint) {
        this.minPoint = minPoint; // Может быть null
    }
    
    public void setDescription(String description) {
        this.description = description; // Может быть null
    }
    
    /**
     * Проверить, прошел ли студент тест
     */
    public boolean isPassed(int studentScore) {
        return minPoint == null || studentScore >= minPoint;
    }
    
    /**
     * Валидация модели перед сохранением
     */
    public void validate() {
        if (title == null || title.trim().isEmpty()) {
            throw new IllegalArgumentException("Title cannot be empty");
        }
        if (courseId == null) {
            throw new IllegalArgumentException("Course ID cannot be null");
        }
        if (minPoint != null && minPoint < 0) {
            throw new IllegalArgumentException("Minimum points cannot be negative");
        }
    }
    
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