package com.example.lms.question.api.dto;
import java.util.UUID;

public class Question {

    private UUID id;
    private UUID testId;
    private String textOfQuestion;
    private Integer order;

    public Question() {
    }

    public Question(UUID id, UUID testId, String textOfQuestion, Integer order) {
        this.id = id;
        this.testId = testId;
        this.textOfQuestion = textOfQuestion;
        this.order = order;
    }

    public UUID getId() {
        return id;
    }

    public void setId(UUID id) {
        this.id = id;
    }

    public UUID getTestId() {
        return testId;
    }

    public void setTestId(UUID testId) {
        this.testId = testId;
    }

    public String getTextOfQuestion() {
        return textOfQuestion;
    }

    public void setTextOfQuestion(String textOfQuestion) {
        this.textOfQuestion = textOfQuestion;
    }

    public Integer getOrder() {
        return order;
    }

    public void setOrder(Integer order) {
        this.order = order;
    }
}