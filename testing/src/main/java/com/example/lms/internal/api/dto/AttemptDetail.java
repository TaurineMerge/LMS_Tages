package com.example.lms.internal.api.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.UUID;

/**
 * DTO для детальной информации о попытке теста.
 * Соответствует JSON-схеме AttemptDetail.
 */
public class AttemptDetail {

    @JsonProperty("attempt_id")
    private UUID attemptId;

    @JsonProperty("student_id")
    private UUID studentId;

    @JsonProperty("test_id")
    private UUID testId;

    @JsonProperty("date_of_attempt")
    private String dateOfAttempt; // ISO date format: "YYYY-MM-DD"

    @JsonProperty("point")
    private Integer point;

    @JsonProperty("completed")
    private Boolean completed;

    @JsonProperty("passed")
    private Boolean passed;

    @JsonProperty("certificate_id")
    private UUID certificateId;

    @JsonProperty("attempt_version")
    private Object attemptVersion; // Пока null

    @JsonProperty("attempt_snapshot_s3")
    private String attemptSnapshot;

    @JsonProperty("meta")
    private Object meta;

    // Constructors
    public AttemptDetail() {
    }

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

    // Getters and Setters
    public UUID getAttemptId() {
        return attemptId;
    }

    public void setAttemptId(UUID attemptId) {
        this.attemptId = attemptId;
    }

    public UUID getStudentId() {
        return studentId;
    }

    public void setStudentId(UUID studentId) {
        this.studentId = studentId;
    }

    public UUID getTestId() {
        return testId;
    }

    public void setTestId(UUID testId) {
        this.testId = testId;
    }

    public String getDateOfAttempt() {
        return dateOfAttempt;
    }

    public void setDateOfAttempt(String dateOfAttempt) {
        this.dateOfAttempt = dateOfAttempt;
    }

    public Integer getPoint() {
        return point;
    }

    public void setPoint(Integer point) {
        this.point = point;
    }

    public Boolean getCompleted() {
        return completed;
    }

    public void setCompleted(Boolean completed) {
        this.completed = completed;
    }

    public Boolean getPassed() {
        return passed;
    }

    public void setPassed(Boolean passed) {
        this.passed = passed;
    }

    public UUID getCertificateId() {
        return certificateId;
    }

    public void setCertificateId(UUID certificateId) {
        this.certificateId = certificateId;
    }

    public Object getAttemptVersion() {
        return attemptVersion;
    }

    public void setAttemptVersion(Object attemptVersion) {
        this.attemptVersion = attemptVersion;
    }

    public String getAttemptSnapshot() {
        return attemptSnapshot;
    }

    public void setAttemptSnapshot(String attemptSnapshot) {
        this.attemptSnapshot = attemptSnapshot;
    }

    public Object getMeta() {
        return meta;
    }

    public void setMeta(Object meta) {
        this.meta = meta;
    }
}
