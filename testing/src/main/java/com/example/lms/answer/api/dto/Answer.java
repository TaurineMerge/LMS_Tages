package com.example.lms.answer.api.dto;
import java.util.UUID;

public class Answer {

    private UUID id;
    private String text;
    private UUID questionId;
    private Integer score;

    public Answer() {
    }


    public Answer(UUID id, String text, UUID questionId, Integer score) {
        this.id = id;
        this.text = text;
        this.questionId = questionId;
        this.score = score;
    }

    public UUID getId() {
        return id;
    }

    public void setId(UUID id) {
        this.id = id;
    }

    public String getText() {
        return text;
    }

    public void setText(String text) {
        this.text = text;
    }

    public UUID getQuestionId() {
        return questionId;
    }

    public void setQuestionId(UUID questionId) {
        this.questionId = questionId;
    }

    public Integer getScore() {
        return score;
    }

    public void setScore(Integer score) {
        this.score = score;
    }
}
