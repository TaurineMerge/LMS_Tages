package com.example.lms.question.api.domain.model;

import java.util.Objects;
import java.util.UUID;

/**
 * Модель вопроса теста (QUESTION_D в БД)
 * Соответствует таблице QUESTION_D
 */
public class QuestionModel {
    private UUID id;
    private UUID testId;           // test_id в БД
    private String textOfQuestion; // text_of_question в БД
    private int order;             // order в БД (порядок в тесте)
    
    // Конструктор для создания нового вопроса
    public QuestionModel(UUID testId, String textOfQuestion, int order) {
        this.testId = Objects.requireNonNull(testId, "Test ID cannot be null");
        this.textOfQuestion = Objects.requireNonNull(textOfQuestion, "Question text cannot be null");
        this.order = order;
    }
    
    // Конструктор для загрузки из БД
    public QuestionModel(UUID id, UUID testId, String textOfQuestion, int order) {
        this.id = id;
        this.testId = testId;
        this.textOfQuestion = textOfQuestion;
        this.order = order;
    }
    
    // Доменные методы (бизнес-логика)
    
    /**
     * Проверить, является ли вопрос валидным
     * @throws IllegalArgumentException если вопрос невалиден
     */
    public void validate() {
        if (textOfQuestion == null || textOfQuestion.trim().isEmpty()) {
            throw new IllegalArgumentException("Question text cannot be empty");
        }
        if (testId == null) {
            throw new IllegalArgumentException("Question must belong to a test");
        }
        if (order < 0) {
            throw new IllegalArgumentException("Order cannot be negative");
        }
    }
    
    /**
     * Изменить порядок вопроса в тесте
     * @param newOrder новый порядковый номер
     */
    public void changeOrder(int newOrder) {
        if (newOrder < 0) {
            throw new IllegalArgumentException("Order cannot be negative");
        }
        this.order = newOrder;
    }
    
    /**
     * Обновить текст вопроса
     * @param newText новый текст вопроса
     */
    public void updateText(String newText) {
        if (newText == null || newText.trim().isEmpty()) {
            throw new IllegalArgumentException("Question text cannot be empty");
        }
        this.textOfQuestion = newText;
    }
    
    // Геттеры
    public UUID getId() { return id; }
    public UUID getTestId() { return testId; }
    public String getTextOfQuestion() { return textOfQuestion; }
    public int getOrder() { return order; }
    
    // Сеттеры с валидацией
    public void setId(UUID id) {
        this.id = Objects.requireNonNull(id, "ID cannot be null");
    }
    
    public void setTestId(UUID testId) {
        this.testId = Objects.requireNonNull(testId, "Test ID cannot be null");
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
    public int hashCode() {
        return Objects.hash(id);
    }
    
    @Override
    public String toString() {
        return "QuestionModel{" +
                "id=" + id +
                ", testId=" + testId +
                ", order=" + order +
                ", text='" + (textOfQuestion != null ? 
                    textOfQuestion.substring(0, Math.min(50, textOfQuestion.length())) + "..." : "null") +
                "'}";
    }
}
