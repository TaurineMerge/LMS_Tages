package com.example.lms.question.api.dto;

import java.util.UUID;

/**
 * DTO для представления вопроса теста.
 * <p>
 * Соответствует таблице question_d:
 * <ul>
 *     <li><b>id</b> – идентификатор вопроса, PK, not null, unique</li>
 *     <li><b>test_id</b> – внешний ключ на тест, not null</li>
 *     <li><b>text_of_question</b> – текст вопроса</li>
 *     <li><b>order</b> – порядок вопроса в тесте</li>
 * </ul>
 * Используется для передачи данных между слоями приложения и в API.
 */
public class Question {

    /**
     * Идентификатор вопроса.
     * PK, not null, unique.
     */
    private UUID id;

    /**
     * Идентификатор теста, к которому относится данный вопрос.
     * FK, not null.
     */
    private UUID testId;

    /**
     * Текст вопроса, отображаемый пользователю.
     */
    private String textOfQuestion;

    /**
     * Порядковый номер вопроса в тесте.
     */
    private Integer order;

    /**
     * Пустой конструктор для сериализации/десериализации.
     */
    public Question() {
    }

    /**
     * Основной конструктор.
     *
     * @param id             идентификатор вопроса
     * @param testId         идентификатор теста
     * @param textOfQuestion текст вопроса
     * @param order          порядковый номер вопроса
     */
    public Question(UUID id, UUID testId, String textOfQuestion, Integer order) {
        this.id = id;
        this.testId = testId;
        this.textOfQuestion = textOfQuestion;
        this.order = order;
    }

    /** @return идентификатор вопроса */
    public UUID getId() {
        return id;
    }

    /** @param id новый идентификатор вопроса */
    public void setId(UUID id) {
        this.id = id;
    }

    /** @return идентификатор теста */
    public UUID getTestId() {
        return testId;
    }

    /** @param testId новый идентификатор теста */
    public void setTestId(UUID testId) {
        this.testId = testId;
    }

    /** @return текст вопроса */
    public String getTextOfQuestion() {
        return textOfQuestion;
    }

    /** @param textOfQuestion новый текст вопроса */
    public void setTextOfQuestion(String textOfQuestion) {
        this.textOfQuestion = textOfQuestion;
    }

    /** @return порядковый номер вопроса */
    public Integer getOrder() {
        return order;
    }

    /** @param order новый порядковый номер вопроса */
    public void setOrder(Integer order) {
        this.order = order;
    }
}