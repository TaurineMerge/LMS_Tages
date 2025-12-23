package com.example.lms.draft.api.dto;

import java.io.Serializable;
import java.util.UUID;

/**
 * DTO для представления черновика теста (draft).
 * <p>
 * Используется для передачи данных о черновике между слоями приложения.
 * Содержит информацию:
 * <ul>
 * <li><b>id</b> — идентификатор черновика</li>
 * <li><b>title</b> — название черновика (обычно название теста)</li>
 * <li><b>min_point</b> — минимальный проходной балл</li>
 * <li><b>description</b> — описание</li>
 * <li><b>testId</b> — идентификатор теста, к которому относится черновик (может быть null)</li>
 * <li><b>courseId</b> — идентификатор курса, к которому относится черновик (может быть null)</li>
 * </ul>
 *
 * Реализует {@link Serializable}, что позволяет передавать объект
 * через сеть или сохранять в файлы, если это необходимо.
 */
public class Draft implements Serializable {

    private static final long serialVersionUID = 1L;

    private UUID id;
    private String title;
    private Integer min_point;
    private String description;
    private UUID testId;    // может быть null
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