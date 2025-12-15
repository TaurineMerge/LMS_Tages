package com.example.lms.question.domain.model;

import java.util.Objects;
import java.util.UUID;

/**
 * Доменная модель вопроса теста.
 * <p>
 * Соответствует таблице {@code QUESTION_D} в БД:
 * <ul>
 *     <li>{@code id} – идентификатор вопроса</li>
 *     <li>{@code test_id} – идентификатор теста</li>
 *     <li>{@code text_of_question} – текст вопроса</li>
 *     <li>{@code order} – порядковый номер вопроса в тесте</li>
 * </ul>
 *
 * Модель включает базовую бизнес-логику:
 * <ul>
 *     <li>валидация состояния {@link #validate()}</li>
 *     <li>изменение порядка {@link #changeOrder(int)}</li>
 *     <li>обновление текста вопроса {@link #updateText(String)}</li>
 * </ul>
 *
 * Используется в сервисном слое и репозиториях как внутренняя структура данных.
 */
public class QuestionModel {

    /** Уникальный идентификатор вопроса (может быть null при создании нового). */
    private UUID id;

    /** Идентификатор теста, которому принадлежит вопрос. */
    private UUID testId;

    /** Текст вопроса, отображаемый студенту. */
    private String textOfQuestion;

    /** Порядковый номер вопроса в тесте. */
    private int order;

    /**
     * Конструктор для создания нового вопроса (без ID).
     *
     * @param testId         идентификатор теста
     * @param textOfQuestion текст вопроса
     * @param order          порядок вопроса в тесте
     * @throws NullPointerException     если testId или textOfQuestion = null
     * @throws IllegalArgumentException если текст пустой или order отрицательный
     */
    public QuestionModel(UUID testId, String textOfQuestion, int order) {
        this.testId = Objects.requireNonNull(testId, "Test ID cannot be null");
        this.textOfQuestion = Objects.requireNonNull(textOfQuestion, "Question text cannot be null");
        this.order = order;

        validate();
    }

    /**
     * Конструктор для загрузки модели из БД.
     *
     * @param id             идентификатор вопроса
     * @param testId         идентификатор теста
     * @param textOfQuestion текст вопроса
     * @param order          порядок вопроса
     */
    public QuestionModel(UUID id, UUID testId, String textOfQuestion, int order) {
        this.id = id;
        this.testId = testId;
        this.textOfQuestion = textOfQuestion;
        this.order = order;

        validate();
    }

    // ------------------------- Бизнес-логика -------------------------

    /**
     * Проверяет корректность модели.
     *
     * @throws IllegalArgumentException если текст пустой, testId = null или order < 0
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
     * Изменяет порядок вопроса в тесте.
     *
     * @param newOrder новый порядковый номер (>= 0)
     * @throws IllegalArgumentException если значение отрицательное
     */
    public void changeOrder(int newOrder) {
        if (newOrder < 0) {
            throw new IllegalArgumentException("Order cannot be negative");
        }
        this.order = newOrder;
    }

    /**
     * Обновляет текст вопроса.
     *
     * @param newText новый текст вопроса
     * @throws IllegalArgumentException если текст пустой
     */
    public void updateText(String newText) {
        if (newText == null || newText.trim().isEmpty()) {
            throw new IllegalArgumentException("Question text cannot be empty");
        }
        this.textOfQuestion = newText;
    }

    // ------------------------- Геттеры -------------------------

    /** @return идентификатор вопроса */
    public UUID getId() { return id; }

    /** @return идентификатор теста */
    public UUID getTestId() { return testId; }

    /** @return текст вопроса */
    public String getTextOfQuestion() { return textOfQuestion; }

    /** @return порядковый номер вопроса */
    public int getOrder() { return order; }

    // ------------------------- Сеттеры -------------------------

    /**
     * Устанавливает идентификатор вопроса.
     *
     * @param id новый ID
     * @throws NullPointerException если id = null
     */
    public void setId(UUID id) {
        this.id = Objects.requireNonNull(id, "ID cannot be null");
    }

    /**
     * Устанавливает идентификатор теста.
     *
     * @param testId новый testId
     * @throws NullPointerException если testId = null
     */
    public void setTestId(UUID testId) {
        this.testId = Objects.requireNonNull(testId, "Test ID cannot be null");
    }

    /**
     * Устанавливает текст вопроса с проверкой.
     *
     * @param textOfQuestion новый текст
     * @throws IllegalArgumentException если текст пустой
     */
    public void setTextOfQuestion(String textOfQuestion) {
        if (textOfQuestion == null || textOfQuestion.trim().isEmpty()) {
            throw new IllegalArgumentException("Question text cannot be empty");
        }
        this.textOfQuestion = textOfQuestion;
    }

    /**
     * Устанавливает порядковый номер вопроса.
     *
     * @param order новый номер
     * @throws IllegalArgumentException если число отрицательное
     */
    public void setOrder(int order) {
        if (order < 0) {
            throw new IllegalArgumentException("Order cannot be negative");
        }
        this.order = order;
    }

    // ------------------------- Equals / HashCode / ToString -------------------------

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
                ", order=" + order +
                ", text='" + (textOfQuestion != null ?
                    textOfQuestion.substring(0, Math.min(50, textOfQuestion.length())) + "..." : "null") +
                "'}";
    }
}