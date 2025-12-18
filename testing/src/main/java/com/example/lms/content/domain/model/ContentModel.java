package com.example.lms.content.domain.model;

import java.util.Objects;
import java.util.UUID;

/**
 * Domain Model: ContentModel
 *
 * Представляет элемент контента курса (вопрос/ответ) и соответствует строке таблицы content_d.
 *
 * Структура таблицы content_d:
 *  - id              UUID    (PK, not null)
 *  - order           INT     — порядковый номер в курсе
 *  - content         VARCHAR — текст контента/вопроса
 *  - type_of_content BOOLEAN — тип контента (true/false для различных типов)
 *  - question_id     UUID    — ссылка на вопрос (может быть null)
 *  - answer_id       UUID    — ссылка на ответ (может быть null)
 *
 * Доменные правила:
 *  - content обязателен
 *  - order должен быть неотрицательным
 *  - typeOfContent может быть null
 *  - questionId и answerId могут быть null
 */
public class ContentModel {

    /** Уникальный ID элемента контента (PK). */
    private UUID id;

    /** Порядковый номер элемента в курсе. */
    private Integer order;

    /** Текст контента/вопроса. */
    private String content;

    /** Тип контента. */
    private Boolean typeOfContent;

    /** ID вопроса. */
    private UUID questionId;

    /** ID ответа. */
    private UUID answerId;

    // ---------------------- КОНСТРУКТОРЫ ----------------------

    /**
     * Конструктор для создания нового элемента контента.
     */
    public ContentModel(Integer order, String content, Boolean typeOfContent, UUID questionId, UUID answerId) {
        setOrder(order);
        setContent(content);
        this.typeOfContent = typeOfContent;
        this.questionId = questionId;
        this.answerId = answerId;
    }

    /**
     * Конструктор для загрузки элемента контента из базы данных.
     */
    public ContentModel(UUID id, Integer order, String content, Boolean typeOfContent, UUID questionId, UUID answerId) {
        this.id = id;
        setOrder(order);
        setContent(content);
        this.typeOfContent = typeOfContent;
        this.questionId = questionId;
        this.answerId = answerId;
    }

    // ---------------------- GETTERS ----------------------

    public UUID getId() { return id; }
    public Integer getOrder() { return order; }
    public String getContent() { return content; }
    public Boolean getTypeOfContent() { return typeOfContent; }
    public UUID getQuestionId() { return questionId; }
    public UUID getAnswerId() { return answerId; }

    // ---------------------- SETTERS С ВАЛИДАЦИЕЙ ----------------------

    public void setId(UUID id) {
        this.id = Objects.requireNonNull(id, "ID cannot be null");
    }

    public void setOrder(Integer order) {
        if (order == null || order < 0) {
            throw new IllegalArgumentException("Order cannot be null and must be non-negative");
        }
        this.order = order;
    }

    public void setContent(String content) {
        if (content == null || content.trim().isEmpty()) {
            throw new IllegalArgumentException("Content cannot be null or empty");
        }
        this.content = content.trim();
    }

    public void setTypeOfContent(Boolean typeOfContent) {
        this.typeOfContent = typeOfContent;
    }

    public void setQuestionId(UUID questionId) {
        this.questionId = questionId;
    }

    public void setAnswerId(UUID answerId) {
        this.answerId = answerId;
    }

    // ---------------------- ДОМЕННАЯ ЛОГИКА ----------------------

    /**
     * Проверка валидности элемента контента перед сохранением.
     */
    public void validate() {
        if (order == null || order < 0) {
            throw new IllegalArgumentException("Order cannot be null and must be non-negative");
        }
        if (content == null || content.trim().isEmpty()) {
            throw new IllegalArgumentException("Content cannot be empty");
        }
    }

    // ---------------------- UTILS ----------------------

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        ContentModel that = (ContentModel) o;
        return Objects.equals(id, that.id);
    }

    @Override
    public int hashCode() {
        return Objects.hash(id);
    }

    @Override
    public String toString() {
        return "ContentModel{" +
                "id=" + id +
                ", order=" + order +
                ", content='" + content + '\'' +
                ", typeOfContent=" + typeOfContent +
                ", questionId=" + questionId +
                ", answerId=" + answerId +
                '}';
    }
}