package com.example.lms.draft.api.dto;

import java.io.Serializable;
import java.util.UUID;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * DTO для представления черновика теста (draft).
 * <p>
 * Используется для передачи данных о черновике между слоями приложения.
 */
public class Draft implements Serializable {

    private static final long serialVersionUID = 1L;

    private UUID id;
    private String title;
    
    @JsonProperty("min_point")
    private Integer min_point;
    
    private String description;
    
    @JsonProperty("test_id")
    private UUID testId;    // может быть null
    
    @JsonProperty("course_id")
    private UUID courseId;  // может быть null

    /**
     * Пустой конструктор для сериализации/десериализации.
     */
    public Draft() {
    }

    /**
     * Основной конструктор.
     */
    public Draft(UUID id, String title, Integer min_point, String description, UUID testId, UUID courseId) {
        this.id = id;
        this.title = title;
        this.min_point = min_point;
        this.description = description;
        this.testId = testId;
        this.courseId = courseId;
    }

    // Геттеры и сеттеры
    public UUID getId() { return id; }
    public void setId(UUID id) { this.id = id; }
    
    public String getTitle() { return title; }
    public void setTitle(String title) { this.title = title; }
    
    public Integer getMin_point() { return min_point; }
    public void setMin_point(Integer min_point) { this.min_point = min_point; }
    
    public String getDescription() { return description; }
    public void setDescription(String description) { this.description = description; }
    
    public UUID getTestId() { return testId; }
    public void setTestId(UUID testId) { this.testId = testId; }
    
    public UUID getCourseId() { return courseId; }
    public void setCourseId(UUID courseId) { this.courseId = courseId; }

    @Override
    public String toString() {
        return "Draft{" +
                "id=" + id +
                ", title='" + title + '\'' +
                ", min_point=" + min_point +
                ", description='" + description + '\'' +
                ", testId=" + testId +
                ", courseId=" + courseId +
                '}';
    }
}