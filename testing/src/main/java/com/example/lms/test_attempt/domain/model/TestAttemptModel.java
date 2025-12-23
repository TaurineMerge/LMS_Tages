package com.example.lms.test_attempt.domain.model;

import java.time.LocalDate;
import java.util.Objects;
import java.util.UUID;

/**
 * Domain Model: TestAttemptModel
 *
 * Представляет попытку прохождения теста студентом и соответствует строке таблицы test_attempt_b.
 *
 * ВАЖНО:
 * - id = PK
 * - date_of_attempt НЕ используем как идентификатор попытки (это не PK)
 * - если в БД стоит UNIQUE(student_id, test_id, date_of_attempt), то автопроставление даты
 *   при завершении может ломать "несколько попыток в один день".
 *   Поэтому в completeAttempt() мы НЕ трогаем dateOfAttempt автоматически.
 */
public class TestAttemptModel {

    private UUID id;
    private UUID studentId;
    private UUID testId;
    private LocalDate dateOfAttempt;
    private Integer point;
    private UUID certificateId;
    private String attemptVersion;
    private String attemptSnapshot;
    private Boolean completed;

    // ---------------------- КОНСТРУКТОРЫ ----------------------

    public TestAttemptModel(UUID studentId, UUID testId) {
        this.studentId = Objects.requireNonNull(studentId, "Student ID cannot be null");
        this.testId = Objects.requireNonNull(testId, "Test ID cannot be null");
        this.completed = false;
    }

    public TestAttemptModel(UUID studentId, UUID testId, LocalDate dateOfAttempt,
                            Integer point, UUID certificateId, String attemptVersion,
                            String attemptSnapshot, Boolean completed) {
        this.studentId = Objects.requireNonNull(studentId, "Student ID cannot be null");
        this.testId = Objects.requireNonNull(testId, "Test ID cannot be null");
        this.dateOfAttempt = dateOfAttempt;
        this.point = point;
        this.certificateId = certificateId;
        this.attemptVersion = attemptVersion;
        this.attemptSnapshot = attemptSnapshot;
        this.completed = completed != null ? completed : false;
    }

    public TestAttemptModel(UUID id, UUID studentId, UUID testId, LocalDate dateOfAttempt,
                            Integer point, UUID certificateId, String attemptVersion,
                            String attemptSnapshot, Boolean completed) {
        this.id = id;
        this.studentId = studentId;
        this.testId = testId;
        this.dateOfAttempt = dateOfAttempt;
        this.point = point;
        this.certificateId = certificateId;
        this.attemptVersion = attemptVersion;
        this.attemptSnapshot = attemptSnapshot;
        this.completed = completed != null ? completed : false;
    }

    // ---------------------- GETTERS ----------------------

    public UUID getId() { return id; }
    public UUID getStudentId() { return studentId; }
    public UUID getTestId() { return testId; }
    public LocalDate getDateOfAttempt() { return dateOfAttempt; }
    public Integer getPoint() { return point; }
    public UUID getCertificateId() { return certificateId; }
    public String getAttemptVersion() { return attemptVersion; }
    public String getAttemptSnapshot() { return attemptSnapshot; }
    public Boolean getCompleted() { return completed; }

    // ---------------------- SETTERS ----------------------

    public void setId(UUID id) {
        this.id = Objects.requireNonNull(id, "ID cannot be null");
    }

    public void setStudentId(UUID studentId) {
        this.studentId = Objects.requireNonNull(studentId, "Student ID cannot be null");
    }

    public void setTestId(UUID testId) {
        this.testId = Objects.requireNonNull(testId, "Test ID cannot be null");
    }

    public void setDateOfAttempt(LocalDate dateOfAttempt) {
        this.dateOfAttempt = dateOfAttempt; // null OK
    }

    public void setPoint(Integer point) {
        if (point != null && point < 0) {
            throw new IllegalArgumentException("Point cannot be negative");
        }
        this.point = point;
    }

    public void setCertificateId(UUID certificateId) {
        this.certificateId = certificateId; // null OK
    }

    public void setAttemptVersion(String attemptVersion) {
        this.attemptVersion = attemptVersion; // null OK
    }

    public void setAttemptSnapshot(String attemptSnapshot) {
        this.attemptSnapshot = attemptSnapshot; // null OK
    }

    public void setCompleted(Boolean completed) {
        this.completed = completed != null ? completed : false;
    }

    // ---------------------- ДОМЕННАЯ ЛОГИКА ----------------------

    /**
     * Помечает попытку как завершенную.
     *
     * ВАЖНО: дату не трогаем автоматически, потому что date_of_attempt не PK и
     * при UNIQUE(student_id,test_id,date_of_attempt) это может ломать множественные попытки в один день.
     */
    public void completeAttempt(int finalPoint) {
        if (finalPoint < 0) {
            throw new IllegalArgumentException("Final point cannot be negative");
        }
        this.point = finalPoint;
        this.completed = true;
    }

    public void updateSnapshot(String newSnapshot) {
        this.attemptSnapshot = newSnapshot;
    }

    public void validate() {
        if (studentId == null) {
            throw new IllegalArgumentException("Student ID cannot be null");
        }
        if (testId == null) {
            throw new IllegalArgumentException("Test ID cannot be null");
        }
        if (point != null && point < 0) {
            throw new IllegalArgumentException("Point cannot be negative");
        }
    }

    public boolean isSuccessful(int minPassingPoint) {
        return point != null && point >= minPassingPoint;
    }

    // ---------------------- UTILS ----------------------

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        TestAttemptModel that = (TestAttemptModel) o;
        return Objects.equals(id, that.id);
    }

    @Override
    public int hashCode() {
        return Objects.hash(id);
    }

    @Override
    public String toString() {
        return "TestAttemptModel{" +
                "id=" + id +
                ", studentId=" + studentId +
                ", testId=" + testId +
                ", dateOfAttempt=" + dateOfAttempt +
                ", point=" + point +
                ", completed=" + completed +
                '}';
    }
}