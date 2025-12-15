package com.example.lms.question.api.domain.model;

import java.util.Objects;
import java.util.UUID;

/**
 * Domain Model: QuestionModel
 *
 * Представляет вопрос теста и соответствует строке в таблице question_d.
 *
 * Структура таблицы question_d:
 *  - id              UUID  (PK, not null)
 *  - test_id         UUID  (FK → test_d.id, not null)
 *  - text_of_question TEXT (может быть null с точки зрения БД, но в доменной модели считаем обязательным)
 *  - order           INT   (порядковый номер вопроса в тесте, может быть null)
 *
 * Договорённости доменной модели:
 *  - Вопрос всегда должен иметь:
 *      - testId (входит в конкретный тест),
 *      - непустой текст (textOfQuestion).
 *  - Порядок (order) может быть не задан (null), но если задан, то не может быть отрицательным.
 */
public class QuestionModel {

    /** Уникальный идентификатор вопроса (PRIMARY KEY). */
    private UUID id;

    /** Идентификатор теста, к которому относится вопрос (FOREIGN KEY). */
    private UUID testId;

    /** Текст вопроса (question_d.text_of_question). */
    private String textOfQuestion;

    /**
     * Порядок вопроса в тесте.
     *
     * Может быть:
     *  - null — если порядок ещё не определён;
     *  - 0, 1, 2, ... — если порядок задан явно.
     */
    private Integer order;

    // ---------------------- КОНСТРУКТОРЫ ----------------------

    /**
     * Конструктор для создания нового вопроса (до сохранения в БД).
     *
     * @param testId         идентификатор теста (not null)
     * @param textOfQuestion текст вопроса (not null, не пустой)
     * @param order          порядок вопроса (может быть null)
     */
    public QuestionModel(UUID testId, String textOfQuestion, Integer order) {
        this.testId = Objects.requireNonNull(testId, "Test ID cannot be null");
        if (textOfQuestion == null || textOfQuestion.trim().isEmpty()) {
            throw new IllegalArgumentException("Question text cannot be null or empty");
        }
        this.textOfQuestion = textOfQuestion;
        this.order = order;
    }

    /**
     * Конструктор для загрузки вопроса из БД.
     *
     * @param id             идентификатор вопроса
     * @param testId         идентификатор теста
     * @param textOfQuestion текст вопроса
     * @param order          порядок вопроса (может быть null)
     */
    public QuestionModel(UUID id, UUID testId, String textOfQuestion, Integer order) {
        this.id = id;
        this.testId = testId;
        this.textOfQuestion = textOfQuestion;
        this.order = order;
    }

    // ---------------------- ДОМЕННАЯ ЛОГИКА ----------------------

    /**
     * Проверить, является ли вопрос валидным с точки зрения домена.
     *
     * @throws IllegalArgumentException если вопрос невалиден
     */
    public void validate() {
        if (textOfQuestion == null || textOfQuestion.trim().isEmpty()) {
            throw new IllegalArgumentException("Question text cannot be empty");
        }
        if (testId == null) {
            throw new IllegalArgumentException("Question must belong to a test");
        }
        if (order != null && order < 0) {
            throw new IllegalArgumentException("Order cannot be negative");
        }
    }

    /**
     * Изменить порядок вопроса в тесте.
     *
     * @param newOrder новый порядковый номер (должен быть >= 0)
     */
    public void changeOrder(int newOrder) {
        if (newOrder < 0) {
            throw new IllegalArgumentException("Order cannot be negative");
        }
        this.order = newOrder;
    }

    /**
     * Обновить текст вопроса.
     *
     * @param newText новый текст вопроса (не может быть пустым)
     */
    public void updateText(String newText) {
        if (newText == null || newText.trim().isEmpty()) {
            throw new IllegalArgumentException("Question text cannot be empty");
        }
        this.textOfQuestion = newText;
    }

    // ---------------------- GETTERS ----------------------

    public UUID getId() {
        return id;
    }

    public UUID getTestId() {
        return testId;
    }

    public String getTextOfQuestion() {
        return textOfQuestion;
    }

    public Integer getOrder() {
        return order;
    }

    // ---------------------- SETTERS ----------------------

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

    public void setOrder(Integer order) {
        if (order != null && order < 0) {
            throw new IllegalArgumentException("Order cannot be negative");
        }
        this.order = order;
    }

    // ---------------------- EQUALS / HASHCODE / TO_STRING ----------------------

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

    /**
     * Краткое строковое представление вопроса.
     * Текст вопроса обрезается до 50 символов для удобства логирования.
     */
    @Override
    public String toString() {
        String shortText = null;
        if (textOfQuestion != null) {
            shortText = textOfQuestion.substring(0, Math.min(50, textOfQuestion.length()));
        }

        return "QuestionModel{" +
                "id=" + id +
                ", testId=" + testId +
                ", order=" + order +
                ", text='" + (shortText != null ? shortText + "..." : "null") +
                "'}";
    }
}
