package com.example.lms.question.domain.model;

import java.util.Objects;
import java.util.UUID;

public class QuestionModel {
    private UUID id;
    private UUID testId;      // может быть null для черновиков
    private UUID draftId;     // может быть null для тестов
    private String textOfQuestion;
    private int order;

    /**
     * Конструктор для создания нового вопроса.
     */
    public QuestionModel(UUID testId, UUID draftId, String textOfQuestion, int order) {
        this.testId = testId;
        this.draftId = draftId;
        this.textOfQuestion = Objects.requireNonNull(textOfQuestion, "Question text cannot be null");
        this.order = order;
        validate();
    }

    /**
     * Конструктор для загрузки модели из БД.
     */
    public QuestionModel(UUID id, UUID testId, UUID draftId, String textOfQuestion, int order) {
        this.id = id;
        this.testId = testId;
        this.draftId = draftId;
        this.textOfQuestion = textOfQuestion;
        this.order = order;
        validate();
    }

    public void validate() {
        if (textOfQuestion == null || textOfQuestion.trim().isEmpty()) {
            throw new IllegalArgumentException("Question text cannot be empty");
        }
        // Разрешаем либо testId, либо draftId, но не оба одновременно
        if (testId == null && draftId == null) {
            throw new IllegalArgumentException("Question must belong to either a test or a draft");
        }
        if (testId != null && draftId != null) {
            throw new IllegalArgumentException("Question cannot belong to both a test and a draft");
        }
        if (order < 0) {
            throw new IllegalArgumentException("Order cannot be negative");
        }
    }

    public void changeOrder(int newOrder) {
        if (newOrder < 0) {
            throw new IllegalArgumentException("Order cannot be negative");
        }
        this.order = newOrder;
    }

    public void updateText(String newText) {
        if (newText == null || newText.trim().isEmpty()) {
            throw new IllegalArgumentException("Question text cannot be empty");
        }
        this.textOfQuestion = newText;
    }

    // Геттеры
    public UUID getId() { return id; }
    public UUID getTestId() { return testId; }
    public UUID getDraftId() { return draftId; }
    public String getTextOfQuestion() { return textOfQuestion; }
    public int getOrder() { return order; }

    // Сеттеры
    public void setId(UUID id) {
        this.id = Objects.requireNonNull(id, "ID cannot be null");
    }
    
    public void setTestId(UUID testId) {
        this.testId = testId;
        validate();
    }
    
    public void setDraftId(UUID draftId) {
        this.draftId = draftId;
        validate();
    }
    
    public void setTextOfQuestion(String textOfQuestion) {
        if (textOfQuestion == null || textOfQuestion.trim().isEmpty()) {
            throw new IllegalArgumentException("Question text cannot be empty");
        }
        this.textOfQuestion = textOfQuestion;
    }
    
    public void setOrder(int order) {
        if (order < 0) {
            throw new IllegalArgumentException("Order cannot be negative");
        }
        this.order = order;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        QuestionModel that = (QuestionModel) o;
        return Objects.equals(id, that.id);
    }

    @Override
    public int hashCode() { return Objects.hash(id); }

    @Override
    public String toString() {
        return "QuestionModel{" +
                "id=" + id +
                ", testId=" + testId +
                ", draftId=" + draftId +
                ", order=" + order +
                ", text='" + (textOfQuestion != null ?
                    textOfQuestion.substring(0, Math.min(50, textOfQuestion.length())) + "..." : "null") +
                "'}";
    }
}