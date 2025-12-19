package com.example.lms.answer.domain.model;

import java.util.Objects;
import java.util.UUID;

/**
 * Доменная модель: AnswerModel
 *
 * Соответствует строке таблицы answer_d.
 *
 * Поля:
 * - id: UUID — первичный ключ
 * - question_id: UUID — идентификатор вопроса (FK → question_d.id)
 * - text: String — текст ответа
 * - score: int — количество баллов за ответ (0 — неверный, >0 — верный)
 * - order: int — порядок ответа
 */
public class AnswerModel {

    /** Уникальный идентификатор ответа (PRIMARY KEY). */
    private UUID id;

    /** Идентификатор вопроса, которому принадлежит данный ответ. */
    private UUID questionId;

    /** Текст ответа, отображаемый пользователю. */
    private String text;

    /** Количество баллов за этот ответ (0 = неверный, >0 = верный). */
    private Integer score;

    /**
     * Конструктор для создания нового ответа (до сохранения в базу данных).
     *
     * @param questionId ID вопроса
     * @param text       текст ответа
     * @param score      количество баллов
     */
    public AnswerModel(UUID questionId, String text, Integer score, Integer order) {
        this.questionId = Objects.requireNonNull(questionId, "Question ID не может быть null");
        this.text = Objects.requireNonNull(text, "Текст ответа не может быть null");
        this.score = Objects.requireNonNull(score, "Score не может быть null");
    }

    /**
     * Конструктор для загрузки ответа из базы данных.
     *
     * @param id         идентификатор ответа
     * @param questionId идентификатор вопроса
     * @param text       текст ответа
     * @param score      количество баллов
     * @param order      порядковый номер
     */
    public AnswerModel(UUID id, UUID questionId, String text, Integer score) {
        this.id = id;
        this.questionId = questionId;
        this.text = text;
        this.score = score;
    }

    // ---------------------- GETTERS ----------------------

    public UUID getId() {
        return id;
    }

    public UUID getQuestionId() {
        return questionId;
    }

    public String getText() {
        return text;
    }

    public Integer getScore() {
        return score;
    }

    // ---------------------- SETTERS ----------------------

    public void setId(UUID id) {
        this.id = Objects.requireNonNull(id, "ID не может быть null");
    }

    public void setQuestionId(UUID questionId) {
        this.questionId = Objects.requireNonNull(questionId, "Question ID не может быть null");
    }

    public void setText(String text) {
        this.text = Objects.requireNonNull(text, "Текст ответа не может быть null");
    }

    public void setScore(Integer score) {
        this.score = Objects.requireNonNull(score, "Score не может быть null");
    }

    // ---------------------- ВАЛИДАЦИЯ ----------------------

    /**
     * Проверяет корректность данных перед сохранением в БД.
     */
    public void validate() {
        if (text == null || text.trim().isEmpty()) {
            throw new IllegalArgumentException("Текст ответа не может быть пустым");
        }
        if (score == null) {
            throw new IllegalArgumentException("Score не может быть null");
        }
        if (questionId == null) {
            throw new IllegalArgumentException("Question ID не может быть null");
        }
    }

    // ---------------------- УТИЛИТЫ ----------------------

    @Override
    public boolean equals(Object o) {
        if (this == o)
            return true;
        if (o == null || getClass() != o.getClass())
            return false;
        AnswerModel that = (AnswerModel) o;
        return Objects.equals(id, that.id);
    }

    @Override
    public int hashCode() {
        return Objects.hash(id);
    }

    /**
     * Краткое строковое представление объекта для удобного логирования.
     */
    @Override
    public String toString() {
        return "AnswerModel{" +
                "id=" + id +
                ", questionId=" + questionId +
                ", score=" + score +
                '}';
    }
}