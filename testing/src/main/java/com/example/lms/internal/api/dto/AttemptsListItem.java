package com.example.lms.internal.api.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.UUID;

/**
 * Data Transfer Object (DTO) для представления элемента списка попыток прохождения тестов.
 * <p>
 * Класс соответствует JSON-схеме AttemptsList и используется для сериализации массива элементов
 * при передаче списка попыток между сервисами системы LMS.
 * <p>
 * Содержит сокращенный набор информации по сравнению с {@link AttemptDetail} для оптимизации
 * передачи данных при работе со списками.
 *
 * @author Команда разработки LMS
 * @version 1.0
 * @see AttemptDetail
 */
public class AttemptsListItem {

    /**
     * Уникальный идентификатор попытки прохождения теста.
     */
    @JsonProperty("attempt_id")
    private UUID attemptId;

    /**
     * Уникальный идентификатор теста, для которого выполнена попытка.
     */
    @JsonProperty("test_id")
    private UUID testId;

    /**
     * Дата выполнения попытки в формате ISO 8601 (YYYY-MM-DD).
     * Пример: "2024-01-15"
     */
    @JsonProperty("date_of_attempt")
    private String dateOfAttempt;

    /**
     * Количество баллов, набранных в попытке.
     * Может быть {@code null}, если попытка еще не оценена.
     */
    @JsonProperty("point")
    private Integer point;

    /**
     * Флаг завершенности попытки.
     * <ul>
     *   <li>{@code true} - попытка завершена</li>
     *   <li>{@code false} - попытка в процессе</li>
     *   <li>{@code null} - статус неизвестен</li>
     * </ul>
     */
    @JsonProperty("is_completed")
    private Boolean isCompleted;

    /**
     * Флаг успешного прохождения теста.
     * <ul>
     *   <li>{@code true} - тест пройден успешно</li>
     *   <li>{@code false} - тест не пройден</li>
     *   <li>{@code null} - результат не определен (попытка не завершена или данные неполные)</li>
     * </ul>
     */
    @JsonProperty("passed")
    private Boolean passed;

    /**
     * Уникальный идентификатор сертификата, выданного за успешное прохождение теста.
     * {@code null}, если сертификат не выдан или попытка не завершена успешно.
     */
    @JsonProperty("certificate_id")
    private UUID certificateId;

    /**
     * Ссылка на снимок (snapshot) попытки в хранилище S3 (MinIO).
     * Содержит полный URL или путь к файлу с детализированными данными попытки.
     */
    @JsonProperty("attempt_snapshot_s3")
    private String attemptSnapshot;

    /**
     * Создает пустой экземпляр AttemptsListItem.
     * <p>
     * Используется фреймворками для десериализации JSON.
     */
    public AttemptsListItem() {
    }

    /**
     * Создает полностью инициализированный экземпляр AttemptsListItem.
     *
     * @param attemptId уникальный идентификатор попытки
     * @param testId уникальный идентификатор теста
     * @param dateOfAttempt дата попытки в формате ISO 8601 (YYYY-MM-DD)
     * @param point количество набранных баллов (может быть {@code null})
     * @param isCompleted флаг завершенности попытки
     * @param passed флаг успешного прохождения (может быть {@code null})
     * @param certificateId идентификатор сертификата (может быть {@code null})
     * @param attemptSnapshot ссылка на снимок попытки в S3 (может быть {@code null})
     */
    public AttemptsListItem(UUID attemptId, UUID testId, String dateOfAttempt,
            Integer point, Boolean isCompleted, Boolean passed, UUID certificateId,
            String attemptSnapshot) {
        this.attemptId = attemptId;
        this.testId = testId;
        this.dateOfAttempt = dateOfAttempt;
        this.point = point;
        this.isCompleted = isCompleted;
        this.passed = passed;
        this.certificateId = certificateId;
        this.attemptSnapshot = attemptSnapshot;
    }

    /**
     * Возвращает уникальный идентификатор попытки.
     *
     * @return идентификатор попытки
     */
    public UUID getAttemptId() {
        return attemptId;
    }

    /**
     * Устанавливает уникальный идентификатор попытки.
     *
     * @param attemptId идентификатор попытки
     */
    public void setAttemptId(UUID attemptId) {
        this.attemptId = attemptId;
    }

    /**
     * Возвращает уникальный идентификатор теста.
     *
     * @return идентификатор теста
     */
    public UUID getTestId() {
        return testId;
    }

    /**
     * Устанавливает уникальный идентификатор теста.
     *
     * @param testId идентификатор теста
     */
    public void setTestId(UUID testId) {
        this.testId = testId;
    }

    /**
     * Возвращает дату выполнения попытки в формате ISO 8601.
     *
     * @return дата попытки в формате "YYYY-MM-DD"
     */
    public String getDateOfAttempt() {
        return dateOfAttempt;
    }

    /**
     * Устанавливает дату выполнения попытки.
     *
     * @param dateOfAttempt дата попытки в формате ISO 8601 (YYYY-MM-DD)
     * @throws IllegalArgumentException если формат даты не соответствует ожидаемому
     */
    public void setDateOfAttempt(String dateOfAttempt) {
        this.dateOfAttempt = dateOfAttempt;
    }

    /**
     * Возвращает количество баллов, набранных в попытке.
     *
     * @return количество баллов или {@code null}, если попытка еще не оценена
     */
    public Integer getPoint() {
        return point;
    }

    /**
     * Устанавливает количество баллов, набранных в попытке.
     *
     * @param point количество баллов (может быть {@code null})
     */
    public void setPoint(Integer point) {
        this.point = point;
    }

    /**
     * Возвращает флаг завершенности попытки.
     *
     * @return {@code true} если попытка завершена, {@code false} если в процессе,
     *         {@code null} если статус неизвестен
     */
    public Boolean getCompleted() {
        return isCompleted;
    }

    /**
     * Устанавливает флаг завершенности попытки.
     *
     * @param completed флаг завершенности (может быть {@code null})
     */
    public void setCompleted(Boolean completed) {
        this.isCompleted = completed;
    }

    /**
     * Возвращает флаг успешного прохождения теста.
     *
     * @return {@code true} если тест пройден успешно, {@code false} если не пройден,
     *         {@code null} если результат не определен
     */
    public Boolean getPassed() {
        return passed;
    }

    /**
     * Устанавливает флаг успешного прохождения теста.
     *
     * @param passed флаг успешного прохождения (может быть {@code null})
     */
    public void setPassed(Boolean passed) {
        this.passed = passed;
    }

    /**
     * Возвращает идентификатор сертификата, выданного за успешное прохождение.
     *
     * @return идентификатор сертификата или {@code null}, если сертификат не выдан
     */
    public UUID getCertificateId() {
        return certificateId;
    }

    /**
     * Устанавливает идентификатор сертификата.
     *
     * @param certificateId идентификатор сертификата (может быть {@code null})
     */
    public void setCertificateId(UUID certificateId) {
        this.certificateId = certificateId;
    }

    /**
     * Возвращает ссылку на снимок попытки в хранилище S3.
     *
     * @return URL или путь к файлу снимка попытки, или {@code null} если снимок отсутствует
     */
    public String getAttemptSnapshot() {
        return attemptSnapshot;
    }

    /**
     * Устанавливает ссылку на снимок попытки в хранилище S3.
     *
     * @param attemptSnapshot URL или путь к файлу снимка попытки (может быть {@code null})
     */
    public void setAttemptSnapshot(String attemptSnapshot) {
        this.attemptSnapshot = attemptSnapshot;
    }
}