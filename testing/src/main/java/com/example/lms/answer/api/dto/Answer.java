package com.example.lms.answer.api.dto;

import java.util.UUID;

import com.fasterxml.jackson.annotation.JsonProperty;

/**
 * DTO для представления ответов на вопрос теста.
 * <p>
 * Соответствует таблице answer_d:
 * <ul>
 * <li><b>id</b> – первичный ключ (UUID), not null, unique</li>
 * <li><b>text</b> – текст ответа</li>
 * <li><b>question_id</b> – внешний ключ на QUESTION_D, not null</li>
 * <li><b>score</b> – балл за данный ответ, not null</li>
 * </ul>
 */
public class Answer {

    /**
     * Идентификатор ответа.
     * <p>
     * PK, not null, unique.
     * </p>
     */
    private UUID id;

    /**
     * Текст ответа, отображаемый пользователю.
     */
    private String text;

    /**
     * Идентификатор вопроса, к которому относится данный ответ.
     * <p>
     * FK на QUESTION_D, not null.
     * </p>
     */
    @JsonProperty("question_id")
    private UUID questionId;

    /**
     * Балл, начисляемый за выбор данного ответа.
     * <p>
     * Значение not null.
     * </p>
     */
    private Integer score;

    /**
     * Порядковый номер ответа в тесте.
     * <p>
     * Значение not null.
     * </p>
     */
    private Integer order;

    /**
     * Пустой конструктор для сериализации/десериализации.
     */
    public Answer() {
    }

    /**
     * Основной конструктор DTO.
     *
     * @param id         идентификатор ответа
     * @param text       текст ответа
     * @param questionId идентификатор связанного вопроса
     * @param score      балл за ответ
     * @param order      порядковый номер ответа
     */
    public Answer(UUID id, String text, UUID questionId, Integer score, Integer order) {
        this.id = id;
        this.text = text;
        this.questionId = questionId;
        this.score = score;
    }

    /** @return идентификатор ответа */
    public UUID getId() {
        return id;
    }

    /** @param id новый идентификатор ответа */
    public void setId(UUID id) {
        this.id = id;
    }

    /** @return текст ответа */
    public String getText() {
        return text;
    }

    /** @param text новый текст ответа */
    public void setText(String text) {
        this.text = text;
    }

    /** @return идентификатор вопроса */
    public UUID getQuestionId() {
        return questionId;
    }

    /** @param questionId новый идентификатор вопроса */
    public void setQuestionId(UUID questionId) {
        this.questionId = questionId;
    }

    /** @return балл за ответ */
    public Integer getScore() {
        return score;
    }

    /** @param score новый балл за ответ */
    public void setScore(Integer score) {
        this.score = score;
    }

    /** @return порядковый номер ответа */
    public Integer getOrder() {
        return order;
    }

    /** @param order новый порядковый номер ответа */
    public void setOrder(Integer order) {
        this.order = order;
    }
}