package com.example.lms.test_attempt.api.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.List;
import java.util.UUID;

/**
 * DTO для снимка (snapshot) попытки прохождения теста.
 * <p>
 * Содержит полную информацию о попытке:
 * <ul>
 * <li>Структуру теста на момент прохождения</li>
 * <li>Все вопросы с контентом</li>
 * <li>Все ответы с контентом и выборами студента</li>
 * </ul>
 * 
 * <p>
 * Снимок сохраняется в MinIO в формате JSON при завершении попытки.
 * 
 * <h2>Пример JSON:</h2>
 * 
 * <pre>
 * {
 *   "attempt_id": "uuid",
 *   "questions": [
 *     {
 *       "question_id": "uuid",
 *       "content": [{"content_id": "uuid", "type": "text", "data": "..."}],
 *       "answers": [
 *         {
 *           "answer_id": "uuid",
 *           "content": {"content_id": "uuid", "type": "text", "data": "..."},
 *           "selected": true,
 *           "score": 5
 *         }
 *       ]
 *     }
 *   ]
 * }
 * </pre>
 * 
 * @see com.example.lms.test_attempt.infrastructure.storage.MinioStorageService
 */
public class AttemptSnapshotDto {

    /**
     * Идентификатор попытки, к которой относится снимок.
     */
    @JsonProperty("attempt_id")
    private UUID attemptId;

    /**
     * Список вопросов теста с ответами и выборами студента.
     */
    @JsonProperty("questions")
    private List<QuestionSnapshotDto> questions;

    // ======================================================================
    // CONSTRUCTORS
    // ======================================================================

    public AttemptSnapshotDto() {
    }

    public AttemptSnapshotDto(UUID attemptId, List<QuestionSnapshotDto> questions) {
        this.attemptId = attemptId;
        this.questions = questions;
    }

    // ======================================================================
    // GETTERS AND SETTERS
    // ======================================================================

    public UUID getAttemptId() {
        return attemptId;
    }

    public void setAttemptId(UUID attemptId) {
        this.attemptId = attemptId;
    }

    public List<QuestionSnapshotDto> getQuestions() {
        return questions;
    }

    public void setQuestions(List<QuestionSnapshotDto> questions) {
        this.questions = questions;
    }

    @Override
    public String toString() {
        return "AttemptSnapshotDto{" +
                "attemptId=" + attemptId +
                ", questions=" + (questions != null ? questions.size() : 0) +
                '}';
    }
}

/**
 * DTO для вопроса в снимке попытки.
 * <p>
 * Содержит идентификатор вопроса, контент и список ответов.
 */
class QuestionSnapshotDto {

    /**
     * Идентификатор вопроса.
     */
    @JsonProperty("question_id")
    private UUID questionId;

    /**
     * Список блоков контента вопроса (текст, изображения).
     * <p>
     * Каждый блок имеет свой ID, тип и данные.
     */
    @JsonProperty("content")
    private List<ContentBlockDto> content;

    /**
     * Список ответов на вопрос с информацией о выборе студента.
     */
    @JsonProperty("answers")
    private List<AnswerSnapshotDto> answers;

    // ======================================================================
    // CONSTRUCTORS
    // ======================================================================

    public QuestionSnapshotDto() {
    }

    public QuestionSnapshotDto(UUID questionId, List<ContentBlockDto> content,
            List<AnswerSnapshotDto> answers) {
        this.questionId = questionId;
        this.content = content;
        this.answers = answers;
    }

    // ======================================================================
    // GETTERS AND SETTERS
    // ======================================================================

    public UUID getQuestionId() {
        return questionId;
    }

    public void setQuestionId(UUID questionId) {
        this.questionId = questionId;
    }

    public List<ContentBlockDto> getContent() {
        return content;
    }

    public void setContent(List<ContentBlockDto> content) {
        this.content = content;
    }

    public List<AnswerSnapshotDto> getAnswers() {
        return answers;
    }

    public void setAnswers(List<AnswerSnapshotDto> answers) {
        this.answers = answers;
    }

    @Override
    public String toString() {
        return "QuestionSnapshotDto{" +
                "questionId=" + questionId +
                ", content=" + (content != null ? content.size() : 0) +
                ", answers=" + (answers != null ? answers.size() : 0) +
                '}';
    }
}

/**
 * DTO для ответа в снимке попытки.
 * <p>
 * Содержит идентификатор ответа, контент, информацию о выборе и баллы.
 */
class AnswerSnapshotDto {

    /**
     * Идентификатор ответа.
     */
    @JsonProperty("answer_id")
    private UUID answerId;

    /**
     * Контент ответа (один блок: текст или изображение).
     */
    @JsonProperty("content")
    private ContentBlockDto content;

    /**
     * Флаг выбора студентом.
     * {@code true} - студент выбрал этот ответ.
     */
    @JsonProperty("selected")
    private Boolean selected;

    /**
     * Количество баллов за этот ответ.
     * <ul>
     * <li>0 - неправильный ответ</li>
     * <li>> 0 - правильный ответ с баллами</li>
     * </ul>
     */
    @JsonProperty("score")
    private Integer score;

    // ======================================================================
    // CONSTRUCTORS
    // ======================================================================

    public AnswerSnapshotDto() {
    }

    public AnswerSnapshotDto(UUID answerId, ContentBlockDto content,
            Boolean selected, Integer score) {
        this.answerId = answerId;
        this.content = content;
        this.selected = selected;
        this.score = score;
    }

    // ======================================================================
    // GETTERS AND SETTERS
    // ======================================================================

    public UUID getAnswerId() {
        return answerId;
    }

    public void setAnswerId(UUID answerId) {
        this.answerId = answerId;
    }

    public ContentBlockDto getContent() {
        return content;
    }

    public void setContent(ContentBlockDto content) {
        this.content = content;
    }

    public Boolean getSelected() {
        return selected;
    }

    public void setSelected(Boolean selected) {
        this.selected = selected;
    }

    public Integer getScore() {
        return score;
    }

    public void setScore(Integer score) {
        this.score = score;
    }

    @Override
    public String toString() {
        return "AnswerSnapshotDto{" +
                "answerId=" + answerId +
                ", selected=" + selected +
                ", score=" + score +
                '}';
    }
}

/**
 * DTO для блока контента (текст или изображение).
 * <p>
 * Используется в вопросах и ответах для хранения различных типов контента.
 */
class ContentBlockDto {

    /**
     * Идентификатор блока контента.
     */
    @JsonProperty("content_id")
    private UUID contentId;

    /**
     * Тип контента.
     * <ul>
     * <li>"text" - текстовый контент</li>
     * <li>"image" - изображение (URL или base64)</li>
     * </ul>
     */
    @JsonProperty("type")
    private String type;

    /**
     * Данные контента.
     * <p>
     * Для текста - сам текст.
     * Для изображения - URL или base64 строка.
     */
    @JsonProperty("data")
    private String data;

    // ======================================================================
    // CONSTRUCTORS
    // ======================================================================

    public ContentBlockDto() {
    }

    public ContentBlockDto(UUID contentId, String type, String data) {
        this.contentId = contentId;
        this.type = type;
        this.data = data;
    }

    // ======================================================================
    // GETTERS AND SETTERS
    // ======================================================================

    public UUID getContentId() {
        return contentId;
    }

    public void setContentId(UUID contentId) {
        this.contentId = contentId;
    }

    public String getType() {
        return type;
    }

    public void setType(String type) {
        this.type = type;
    }

    public String getData() {
        return data;
    }

    public void setData(String data) {
        this.data = data;
    }

    @Override
    public String toString() {
        return "ContentBlockDto{" +
                "contentId=" + contentId +
                ", type='" + type + '\'' +
                ", data='" + (data != null && data.length() > 50 ? data.substring(0, 50) + "..." : data) + '\'' +
                '}';
    }
}