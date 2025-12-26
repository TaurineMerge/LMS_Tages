package com.example.lms.internal.api.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.UUID;

/**
 * Data Transfer Object (DTO) для передачи детальной информации о попытке прохождения теста.
 * <p>
 * Класс соответствует JSON-схеме AttemptDetail и используется для сериализации/десериализации
 * данных при взаимодействии с другими сервисами системы LMS.
 * <p>
 * Все поля аннотированы {@link JsonProperty} для корректного маппинга между JSON и объектом.
 * Формат даты соответствует стандарту ISO 8601: "YYYY-MM-DD".
 *
 * @author Команда разработки LMS
 * @version 1.0
 * @see <a href="https://en.wikipedia.org/wiki/ISO_8601">ISO 8601 Date Format</a>
 */
public class AttemptDetail {

    /**
     * Уникальный идентификатор попытки прохождения теста.
     */
    @JsonProperty("attempt_id")
    private UUID attemptId;

    /**
     * Уникальный идентификатор студента, выполнившего попытку.
     */
    @JsonProperty("student_id")
    private UUID studentId;

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
     * {@code true} - попытка завершена, {@code false} - попытка в процессе.
     */
    @JsonProperty("completed")
    private Boolean completed;

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
     * {@code null}, если сертификат не выдан.
     */
    @JsonProperty("certificate_id")
    private UUID certificateId;

    /**
     * Версия попытки для отслеживания изменений.
     * В текущей реализации всегда {@code null}.
     */
    @JsonProperty("attempt_version")
    private Object attemptVersion;

    /**
     * Ссылка на снимок (snapshot) попытки в хранилище S3 (MinIO).
     * Содержит полный URL или путь к файлу с детализированными данными попытки.
     */
    @JsonProperty("attempt_snapshot_s3")
    private String attemptSnapshot;

    /**
     * Метаданные попытки для расширенной информации.
     * Может содержать дополнительные поля в формате ключ-значение.
     */
    @JsonProperty("meta")
    private Object meta;

    /**
     * Создает пустой экземпляр AttemptDetail.
     * <p>
     * Используется фреймворками для десериализации JSON.
     */
    public AttemptDetail() {
    }

    /**
     * Создает полностью инициализированный экземпляр AttemptDetail.
     *
     * @param attemptId уникальный идентификатор попытки
     * @param studentId уникальный идентификатор студента
     * @param testId уникальный идентификатор теста
     * @param dateOfAttempt дата попытки в формате ISO 8601 (YYYY-MM-DD)
     * @param point количество набранных баллов (может быть {@code null})
     * @param completed флаг завершенности попытки
     * @param passed флаг успешного прохождения (может быть {@code null})
     * @param certificateId идентификатор сертификата (может быть {@code null})
     * @param attemptVersion версия попытки (в текущей реализации {@code null})
     * @param attemptSnapshot ссылка на снимок попытки в S3
     * @param meta метаданные попытки (может быть {@code null})
     */
    public AttemptDetail(UUID attemptId, UUID studentId, UUID testId, String dateOfAttempt,
            Integer point, Boolean completed, Boolean passed, UUID certificateId,
            Object attemptVersion, String attemptSnapshot, Object meta) {
        this.attemptId = attemptId;
        this.studentId = studentId;
        this.testId = testId;
        this.dateOfAttempt = dateOfAttempt;
        this.point = point;
        this.completed = completed;
        this.passed = passed;
        this.certificateId = certificateId;
        this.attemptVersion = attemptVersion;
        this.attemptSnapshot = attemptSnapshot;
        this.meta = meta;
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
     * Возвращает уникальный идентификатор студента.
     *
     * @return идентификатор студента
     */
    public UUID getStudentId() {
        return studentId;
    }

    /**
     * Устанавливает уникальный идентификатор студента.
     *
     * @param studentId идентификатор студента
     */
    public void setStudentId(UUID studentId) {
        this.studentId = studentId;
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
     * @return {@code true} если попытка завершена, {@code false} если в процессе
     */
    public Boolean getCompleted() {
        return completed;
    }

    /**
     * Устанавливает флаг завершенности попытки.
     *
     * @param completed флаг завершенности
     */
    public void setCompleted(Boolean completed) {
        this.completed = completed;
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
     * Возвращает версию попытки.
     * <p>
     * В текущей реализации всегда возвращает {@code null}.
     *
     * @return версия попытки (в текущей реализации {@code null})
     */
    public Object getAttemptVersion() {
        return attemptVersion;
    }

    /**
     * Устанавливает версию попытки.
     *
     * @param attemptVersion версия попытки
     */
    public void setAttemptVersion(Object attemptVersion) {
        this.attemptVersion = attemptVersion;
    }

    /**
     * Возвращает ссылку на снимок попытки в хранилище S3.
     *
     * @return URL или путь к файлу снимка попытки
     */
    public String getAttemptSnapshot() {
        return attemptSnapshot;
    }

    /**
     * Устанавливает ссылку на снимок попытки в хранилище S3.
     *
     * @param attemptSnapshot URL или путь к файлу снимка попытки
     */
    public void setAttemptSnapshot(String attemptSnapshot) {
        this.attemptSnapshot = attemptSnapshot;
    }

    /**
     * Возвращает метаданные попытки.
     *
     * @return объект метаданных или {@code null}, если метаданные отсутствуют
     */
    public Object getMeta() {
        return meta;
    }

    /**
     * Устанавливает метаданные попытки.
     *
     * @param meta объект метаданных (может быть {@code null})
     */
    public void setMeta(Object meta) {
        this.meta = meta;
    }
}