package com.example.lms.content.api.dto;

import java.io.Serializable;

/**
 * DTO для представления элемента контента.
 * <p>
 * Используется для передачи данных о элементах контента между слоями приложения.
 * Содержит информацию:
 * <ul>
 * <li><b>id</b> — идентификатор элемента контента</li>
 * <li><b>order</b> — порядковый номер элемента</li>
 * <li><b>content</b> — текст/содержимое элемента</li>
 * <li><b>typeOfContent</b> — тип контента</li>
 * <li><b>questionId</b> — идентификатор вопроса (если привязан)</li>
 * <li><b>answerId</b> — идентификатор ответа (если привязан)</li>
 * </ul>
 *
 * Реализует {@link Serializable}, что позволяет передавать объект
 * через сеть или сохранять в файлы, если это необходимо.
 */
public class Content implements Serializable {

    private static final long serialVersionUID = 1L;

    /**
     * Идентификатор элемента контента.
     */
    private String id;

    /**
     * Порядковый номер элемента контента.
     */
    private Integer order;

    /**
     * Содержимое элемента контента.
     */
    private String content;

    /**
     * Тип контента.
     */
    private Boolean typeOfContent;

    /**
     * Идентификатор вопроса (если привязан к вопросу).
     */
    private String questionId;

    /**
     * Идентификатор ответа (если привязан к ответу).
     */
    private String answerId;

    /**
     * Пустой конструктор для сериализации/десериализации.
     */
    public Content() {
    }

    /**
     * Основной конструктор.
     *
     * @param id             идентификатор элемента контента
     * @param order          порядковый номер
     * @param content        содержимое элемента
     * @param typeOfContent  тип контента
     * @param questionId     идентификатор вопроса
     * @param answerId       идентификатор ответа
     */
    public Content(String id, Integer order, String content, Boolean typeOfContent, 
                   String questionId, String answerId) {
        this.id = id;
        this.order = order;
        this.content = content;
        this.typeOfContent = typeOfContent;
        this.questionId = questionId;
        this.answerId = answerId;
    }

    /** 
     * @return идентификатор элемента контента 
     */
    public String getId() {
        return id;
    }

    /** 
     * @param id новый идентификатор элемента контента 
     */
    public void setId(String id) {
        this.id = id;
    }

    /** 
     * @return порядковый номер элемента контента 
     */
    public Integer getOrder() {
        return order;
    }

    /** 
     * @param order новый порядковый номер элемента контента 
     */
    public void setOrder(Integer order) {
        this.order = order;
    }

    /** 
     * @return содержимое элемента контента 
     */
    public String getContent() {
        return content;
    }

    /** 
     * @param content новое содержимое элемента контента 
     */
    public void setContent(String content) {
        this.content = content;
    }

    /** 
     * @return тип контента 
     */
    public Boolean getTypeOfContent() {
        return typeOfContent;
    }

    /** 
     * @param typeOfContent новый тип контента 
     */
    public void setTypeOfContent(Boolean typeOfContent) {
        this.typeOfContent = typeOfContent;
    }

    /** 
     * @return идентификатор вопроса 
     */
    public String getQuestionId() {
        return questionId;
    }

    /** 
     * @param questionId новый идентификатор вопроса 
     */
    public void setQuestionId(String questionId) {
        this.questionId = questionId;
    }

    /** 
     * @return идентификатор ответа 
     */
    public String getAnswerId() {
        return answerId;
    }

    /** 
     * @param answerId новый идентификатор ответа 
     */
    public void setAnswerId(String answerId) {
        this.answerId = answerId;
    }

    /**
     * Формирует строковое представление объекта Content.
     *
     * @return строка с полями объекта
     */
    @Override
    public String toString() {
        return "Content{" +
                "id='" + id + '\'' +
                ", order=" + order +
                ", content='" + content + '\'' +
                ", typeOfContent=" + typeOfContent +
                ", questionId='" + questionId + '\'' +
                ", answerId='" + answerId + '\'' +
                '}';
    }
}