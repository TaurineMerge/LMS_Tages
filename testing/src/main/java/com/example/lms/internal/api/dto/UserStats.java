package com.example.lms.internal.api.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.List;
import java.util.UUID;

/**
 * DTO для статистики пользователя.
 * Соответствует JSON-схеме UserStats.
 */
public class UserStats {

    @JsonProperty("user_id")
    private UUID userId;

    @JsonProperty("attempts_total")
    private Integer attemptsTotal;

    @JsonProperty("attempts_passed")
    private Integer attemptsPassed;

    @JsonProperty("best_score")
    private Integer bestScore;

    @JsonProperty("last_attempt_at")
    private String lastAttemptAt; // ISO 8601 date-time format

    @JsonProperty("per_test")
    private List<PerTestStats> perTest;

    // Constructors
    public UserStats() {
    }

    public UserStats(UUID userId, Integer attemptsTotal, Integer attemptsPassed,
            Integer bestScore, String lastAttemptAt, List<PerTestStats> perTest) {
        this.userId = userId;
        this.attemptsTotal = attemptsTotal;
        this.attemptsPassed = attemptsPassed;
        this.bestScore = bestScore;
        this.lastAttemptAt = lastAttemptAt;
        this.perTest = perTest;
    }

    // Getters and Setters
    public UUID getUserId() {
        return userId;
    }

    public void setUserId(UUID userId) {
        this.userId = userId;
    }

    public Integer getAttemptsTotal() {
        return attemptsTotal;
    }

    public void setAttemptsTotal(Integer attemptsTotal) {
        this.attemptsTotal = attemptsTotal;
    }

    public Integer getAttemptsPassed() {
        return attemptsPassed;
    }

    public void setAttemptsPassed(Integer attemptsPassed) {
        this.attemptsPassed = attemptsPassed;
    }

    public Integer getBestScore() {
        return bestScore;
    }

    public void setBestScore(Integer bestScore) {
        this.bestScore = bestScore;
    }

    public String getLastAttemptAt() {
        return lastAttemptAt;
    }

    public void setLastAttemptAt(String lastAttemptAt) {
        this.lastAttemptAt = lastAttemptAt;
    }

    public List<PerTestStats> getPerTest() {
        return perTest;
    }

    public void setPerTest(List<PerTestStats> perTest) {
        this.perTest = perTest;
    }
}