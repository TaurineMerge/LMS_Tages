package com.example.lms.internal.api.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.List;
import java.util.UUID;

/**
 * Data Transfer Object (DTO) для представления агрегированной статистики пользователя по всем тестам.
 * <p>
 * Класс соответствует JSON-схеме UserStats и содержит обобщенные метрики по всем попыткам прохождения тестов,
 * включая как общую статистику, так и детализированные данные по каждому тесту.
 * <p>
 * Используется для предоставления сводной информации о прогрессе пользователя в системе тестирования.
 *
 * @author Команда разработки LMS
 * @version 1.0
 * @see PerTestStats
 * @see com.example.lms.internal.service.InternalApiService
 */
public class UserStats {

    /**
     * Уникальный идентификатор пользователя, для которого собрана статистика.
     */
    @JsonProperty("user_id")
    private UUID userId;

    /**
     * Общее количество попыток прохождения тестов пользователем.
     * Включает все попытки (как завершенные, так и незавершенные).
     * Значение всегда неотрицательное целое число.
     */
    @JsonProperty("attempts_total")
    private Integer attemptsTotal;

    /**
     * Количество успешно пройденных попыток тестов.
     * Попытка считается успешной, если тест пройден (набрано достаточное количество баллов).
     * Значение всегда неотрицательное целое число и не может превышать {@code attemptsTotal}.
     */
    @JsonProperty("attempts_passed")
    private Integer attemptsPassed;

    /**
     * Наивысший балл, набранный пользователем за все попытки прохождения тестов.
     * Может быть {@code null}, если пользователь еще не проходил тесты или ни одна попытка не оценена.
     */
    @JsonProperty("best_score")
    private Integer bestScore;

    /**
     * Дата и время последней завершенной попытки в формате ISO 8601.
     * Формат: "YYYY-MM-DD" или "YYYY-MM-DDTHH:MM:SSZ".
     * Может быть {@code null}, если у пользователя нет завершенных попыток.
     */
    @JsonProperty("last_attempt_at")
    private String lastAttemptAt;

    /**
     * Детализированная статистика по каждому тесту, который пользователь пытался пройти.
     * Содержит список объектов {@link PerTestStats} с агрегированными данными по каждому тесту.
     * Может быть пустым списком, если у пользователя нет попыток.
     */
    @JsonProperty("per_test")
    private List<PerTestStats> perTest;

    /**
     * Создает пустой экземпляр UserStats.
     * <p>
     * Используется фреймворками для десериализации JSON.
     */
    public UserStats() {
    }

    /**
     * Создает полностью инициализированный экземпляр UserStats.
     *
     * @param userId уникальный идентификатор пользователя
     * @param attemptsTotal общее количество попыток (не может быть отрицательным)
     * @param attemptsPassed количество успешных попыток (не может быть отрицательным и не может превышать attemptsTotal)
     * @param bestScore наивысший набранный балл (может быть {@code null})
     * @param lastAttemptAt дата последней завершенной попытки в формате ISO 8601 (может быть {@code null})
     * @param perTest список детализированной статистики по тестам (может быть пустым, но не {@code null})
     * @throws IllegalArgumentException если {@code attemptsTotal} или {@code attemptsPassed} отрицательные,
     *                                  или если {@code attemptsPassed} превышает {@code attemptsTotal}
     * @throws NullPointerException если {@code perTest} равен {@code null}
     */
    public UserStats(UUID userId, Integer attemptsTotal, Integer attemptsPassed,
            Integer bestScore, String lastAttemptAt, List<PerTestStats> perTest) {
        this.userId = userId;
        this.attemptsTotal = attemptsTotal;
        this.attemptsPassed = attemptsPassed;
        this.bestScore = bestScore;
        this.lastAttemptAt = lastAttemptAt;
        this.perTest = perTest;
    }

    /**
     * Возвращает уникальный идентификатор пользователя.
     *
     * @return идентификатор пользователя
     */
    public UUID getUserId() {
        return userId;
    }

    /**
     * Устанавливает уникальный идентификатор пользователя.
     *
     * @param userId идентификатор пользователя
     */
    public void setUserId(UUID userId) {
        this.userId = userId;
    }

    /**
     * Возвращает общее количество попыток прохождения тестов.
     *
     * @return общее количество попыток (неотрицательное целое число)
     */
    public Integer getAttemptsTotal() {
        return attemptsTotal;
    }

    /**
     * Устанавливает общее количество попыток прохождения тестов.
     *
     * @param attemptsTotal общее количество попыток (не может быть отрицательным)
     * @throws IllegalArgumentException если {@code attemptsTotal} отрицательное
     */
    public void setAttemptsTotal(Integer attemptsTotal) {
        this.attemptsTotal = attemptsTotal;
    }

    /**
     * Возвращает количество успешно пройденных попыток.
     *
     * @return количество успешных попыток (неотрицательное целое число)
     */
    public Integer getAttemptsPassed() {
        return attemptsPassed;
    }

    /**
     * Устанавливает количество успешно пройденных попыток.
     *
     * @param attemptsPassed количество успешных попыток (не может быть отрицательным)
     * @throws IllegalArgumentException если {@code attemptsPassed} отрицательное
     */
    public void setAttemptsPassed(Integer attemptsPassed) {
        this.attemptsPassed = attemptsPassed;
    }

    /**
     * Возвращает наивысший балл, набранный пользователем за все попытки.
     *
     * @return наивысший балл или {@code null}, если нет оцененных попыток
     */
    public Integer getBestScore() {
        return bestScore;
    }

    /**
     * Устанавливает наивысший балл, набранный пользователем за все попытки.
     *
     * @param bestScore наивысший балл (может быть {@code null})
     */
    public void setBestScore(Integer bestScore) {
        this.bestScore = bestScore;
    }

    /**
     * Возвращает дату и время последней завершенной попытки.
     *
     * @return дата последней попытки в формате ISO 8601 или {@code null}, если нет завершенных попыток
     */
    public String getLastAttemptAt() {
        return lastAttemptAt;
    }

    /**
     * Устанавливает дату и время последней завершенной попытки.
     *
     * @param lastAttemptAt дата последней попытки в формате ISO 8601 (может быть {@code null})
     */
    public void setLastAttemptAt(String lastAttemptAt) {
        this.lastAttemptAt = lastAttemptAt;
    }

    /**
     * Возвращает детализированную статистику по каждому тесту.
     *
     * @return список объектов {@link PerTestStats} (может быть пустым, но не {@code null})
     */
    public List<PerTestStats> getPerTest() {
        return perTest;
    }

    /**
     * Устанавливает детализированную статистику по каждому тесту.
     *
     * @param perTest список объектов {@link PerTestStats} (не может быть {@code null}, может быть пустым)
     * @throws NullPointerException если {@code perTest} равен {@code null}
     */
    public void setPerTest(List<PerTestStats> perTest) {
        this.perTest = perTest;
    }
}