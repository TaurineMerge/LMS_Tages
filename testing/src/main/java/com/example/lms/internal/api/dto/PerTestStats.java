package com.example.lms.internal.api.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.UUID;

/**
 * Data Transfer Object (DTO) для представления статистики по конкретному тесту в рамках общей статистики пользователя.
 * <p>
 * Используется как часть объекта {@link UserStats} для предоставления детализированной информации
 * о результатах пользователя по каждому отдельному тесту. Класс содержит агрегированные данные
 * по всем попыткам прохождения конкретного теста.
 *
 * @author Команда разработки LMS
 * @version 1.0
 * @see UserStats
 * @see com.example.lms.internal.service.InternalApiService
 */
public class PerTestStats {

    /**
     * Уникальный идентификатор теста, для которого собрана статистика.
     */
    @JsonProperty("test_id")
    private UUID testId;

    /**
     * Название теста для удобства идентификации.
     * Пример: "Основы Java", "Алгоритмы и структуры данных"
     */
    @JsonProperty("test_title")
    private String testTitle;

    /**
     * Общее количество попыток прохождения данного теста пользователем.
     * Значение всегда неотрицательное целое число.
     */
    @JsonProperty("attempts")
    private Integer attempts;

    /**
     * Наивысший балл, набранный пользователем в данном тесте.
     * Может быть {@code null}, если тест еще не был пройден или не оценен.
     */
    @JsonProperty("best_score")
    private Integer bestScore;

    /**
     * Количество успешно пройденных попыток данного теста.
     * Попытка считается успешной, если набрано достаточное количество баллов
     * для прохождения теста. Значение всегда неотрицательное целое число.
     */
    @JsonProperty("passed_count")
    private Integer passedCount;

    /**
     * Создает пустой экземпляр PerTestStats.
     * <p>
     * Используется фреймворками для десериализации JSON.
     */
    public PerTestStats() {
    }

    /**
     * Создает полностью инициализированный экземпляр PerTestStats.
     *
     * @param testId уникальный идентификатор теста
     * @param testTitle название теста для отображения
     * @param attempts общее количество попыток (не может быть отрицательным)
     * @param bestScore наивысший набранный балл (может быть {@code null})
     * @param passedCount количество успешных попыток (не может быть отрицательным)
     * @throws IllegalArgumentException если {@code attempts} или {@code passedCount} отрицательные
     */
    public PerTestStats(UUID testId, String testTitle, Integer attempts,
            Integer bestScore, Integer passedCount) {
        this.testId = testId;
        this.testTitle = testTitle;
        this.attempts = attempts;
        this.bestScore = bestScore;
        this.passedCount = passedCount;
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
     * Возвращает название теста для отображения.
     *
     * @return название теста
     */
    public String getTestTitle() {
        return testTitle;
    }

    /**
     * Устанавливает название теста для отображения.
     *
     * @param testTitle название теста
     * @throws IllegalArgumentException если {@code testTitle} равен {@code null} или пустой строке
     */
    public void setTestTitle(String testTitle) {
        this.testTitle = testTitle;
    }

    /**
     * Возвращает общее количество попыток прохождения теста.
     *
     * @return количество попыток (неотрицательное целое число)
     */
    public Integer getAttempts() {
        return attempts;
    }

    /**
     * Устанавливает общее количество попыток прохождения теста.
     *
     * @param attempts количество попыток (не может быть отрицательным)
     * @throws IllegalArgumentException если {@code attempts} отрицательное
     */
    public void setAttempts(Integer attempts) {
        this.attempts = attempts;
    }

    /**
     * Возвращает наивысший балл, набранный пользователем в тесте.
     *
     * @return наивысший балл или {@code null}, если тест не был пройден или не оценен
     */
    public Integer getBestScore() {
        return bestScore;
    }

    /**
     * Устанавливает наивысший балл, набранный пользователем в тесте.
     *
     * @param bestScore наивысший балл (может быть {@code null})
     */
    public void setBestScore(Integer bestScore) {
        this.bestScore = bestScore;
    }

    /**
     * Возвращает количество успешно пройденных попыток теста.
     *
     * @return количество успешных попыток (неотрицательное целое число)
     */
    public Integer getPassedCount() {
        return passedCount;
    }

    /**
     * Устанавливает количество успешно пройденных попыток теста.
     *
     * @param passedCount количество успешных попыток (не может быть отрицательным)
     * @throws IllegalArgumentException если {@code passedCount} отрицательное
     */
    public void setPassedCount(Integer passedCount) {
        this.passedCount = passedCount;
    }
}