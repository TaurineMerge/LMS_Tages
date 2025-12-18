package com.example.lms.internal.api.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.UUID;

/**
 * DTO для статистики по конкретному тесту внутри UserStats.
 */
public class PerTestStats {

    @JsonProperty("test_id")
    private UUID testId;

    @JsonProperty("test_title")
    private String testTitle;

    @JsonProperty("attempts")
    private Integer attempts;

    @JsonProperty("best_score")
    private Integer bestScore;

    @JsonProperty("passed_count")
    private Integer passedCount;

    // Constructors
    public PerTestStats() {
    }

    public PerTestStats(UUID testId, String testTitle, Integer attempts,
            Integer bestScore, Integer passedCount) {
        this.testId = testId;
        this.testTitle = testTitle;
        this.attempts = attempts;
        this.bestScore = bestScore;
        this.passedCount = passedCount;
    }

    // Getters and Setters
    public UUID getTestId() {
        return testId;
    }

    public void setTestId(UUID testId) {
        this.testId = testId;
    }

    public String getTestTitle() {
        return testTitle;
    }

    public void setTestTitle(String testTitle) {
        this.testTitle = testTitle;
    }

    public Integer getAttempts() {
        return attempts;
    }

    public void setAttempts(Integer attempts) {
        this.attempts = attempts;
    }

    public Integer getBestScore() {
        return bestScore;
    }

    public void setBestScore(Integer bestScore) {
        this.bestScore = bestScore;
    }

    public Integer getPassedCount() {
        return passedCount;
    }

    public void setPassedCount(Integer passedCount) {
        this.passedCount = passedCount;
    }
}