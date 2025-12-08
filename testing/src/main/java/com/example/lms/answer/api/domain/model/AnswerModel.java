package com.example.lms.answer.api.domain.model;

import java.util.Objects;
import java.util.UUID;

/**
 * Модель ответа - представляет строку таблицы answer_d
 */
public class AnswerModel {
    private UUID id;
    private UUID questionId;
    private String answerText;
    private Boolean isCorrect;
    private Integer displayOrder;
    private String explanation;
    
    // Конструктор для создания нового ответа
    public AnswerModel(UUID questionId, String answerText, Boolean isCorrect, 
                      Integer displayOrder, String explanation) {
        this.questionId = Objects.requireNonNull(questionId, "Question ID cannot be null");
        this.answerText = Objects.requireNonNull(answerText, "Answer text cannot be null");
        this.isCorrect = Objects.requireNonNull(isCorrect, "IsCorrect flag cannot be null");
        this.displayOrder = displayOrder;
        this.explanation = explanation;
    }
    
    // Конструктор для загрузки из БД
    public AnswerModel(UUID id, UUID questionId, String answerText, Boolean isCorrect, 
                      Integer displayOrder, String explanation) {
        this.id = id;
        this.questionId = questionId;
        this.answerText = answerText;
        this.isCorrect = isCorrect;
        this.displayOrder = displayOrder;
        this.explanation = explanation;
    }
    
    // Геттеры
    public UUID getId() { return id; }
    public UUID getQuestionId() { return questionId; }
    public String getAnswerText() { return answerText; }
    public Boolean getIsCorrect() { return isCorrect; }
    public Integer getDisplayOrder() { return displayOrder; }
    public String getExplanation() { return explanation; }
    
    // Сеттеры с валидацией
    public void setId(UUID id) { 
        this.id = Objects.requireNonNull(id, "ID cannot be null");
    }
    
    public void setQuestionId(UUID questionId) {
        this.questionId = Objects.requireNonNull(questionId, "Question ID cannot be null");
    }
    
    public void setAnswerText(String answerText) {
        this.answerText = Objects.requireNonNull(answerText, "Answer text cannot be null");
    }
    
    public void setIsCorrect(Boolean isCorrect) {
        this.isCorrect = Objects.requireNonNull(isCorrect, "IsCorrect flag cannot be null");
    }
    
    public void setDisplayOrder(Integer displayOrder) {
        this.displayOrder = displayOrder; // Может быть null
    }
    
    public void setExplanation(String explanation) {
        this.explanation = explanation; // Может быть null
    }
    
    /**
     * Проверить, является ли ответ правильным
     */
    public boolean isCorrectAnswer() {
        return Boolean.TRUE.equals(isCorrect);
    }
    
    /**
     * Валидация модели перед сохранением
     */
    public void validate() {
        if (answerText == null || answerText.trim().isEmpty()) {
            throw new IllegalArgumentException("Answer text cannot be empty");
        }
        if (questionId == null) {
            throw new IllegalArgumentException("Question ID cannot be null");
        }
        if (isCorrect == null) {
            throw new IllegalArgumentException("IsCorrect flag cannot be null");
        }
        if (displayOrder != null && displayOrder < 0) {
            throw new IllegalArgumentException("Display order cannot be negative");
        }
    }
    
    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        AnswerModel that = (AnswerModel) o;
        return Objects.equals(id, that.id);
    }
    
    @Override
    public int hashCode() {
        return Objects.hash(id);
    }
    
    @Override
    public String toString() {
        return "AnswerModel{" +
                "id=" + id +
                ", questionId=" + questionId +
                ", isCorrect=" + isCorrect +
                ", displayOrder=" + displayOrder +
                '}';
    }
}
