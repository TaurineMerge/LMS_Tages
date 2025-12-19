package com.example.lms.test_attempt.domain.model;

import java.time.LocalDate;
import java.util.Objects;
import java.util.UUID;

/**
 * Domain Model: TestAttemptModel
 *
 * Представляет попытку прохождения теста студентом и соответствует строке таблицы test_attempt_b.
 *
 * Структура таблицы test_attempt_b:
 *  - id                UUID    (PK, not null)
 *  - student_id        UUID    (not null)
 *  - test_id           UUID    (FK → test_d.id, not null)
 *  - date_of_attempt   DATE    (может быть null)
 *  - point             INTEGER (может быть null)
 *  - certificate_id    UUID    (может быть null)
 *  - attempt_version   JSON    (может быть null)
 *  - attempt_snapshot  VARCHAR(256) (может быть null)
 *  - completed         BOOLEAN (может быть null)
 *  - UNIQUE(student_id, test_id, date_of_attempt)
 *
 * Доменные правила:
 *  - studentId, testId обязательны
 *  - dateOfAttempt может быть null (если попытка еще не начата)
 *  - point >= 0, если указано
 *  - completed по умолчанию false, если не указано
 */
public class TestAttemptModel {

    /** Уникальный ID попытки (PK). */
    private UUID id;

    /** ID студента, проходящего тест. */
    private UUID studentId;

    /** ID теста, который проходит студент. */
    private UUID testId;

    /** Дата попытки. */
    private LocalDate dateOfAttempt;

    /** Количество набранных баллов (может быть null). */
    private Integer point;

    /** ID сертификата (может быть null). */
    private UUID certificateId;

    /** Версия попытки в формате JSON (метаданные, версия теста). */
    private String attemptVersion;

    /** Снапшот теста на момент прохождения (JSON в виде строки). */
    private String attemptSnapshot;

    /** Флаг завершения попытки. */
    private Boolean completed;

    // ---------------------- КОНСТРУКТОРЫ ----------------------

    /**
     * 
     */
    public TestAttemptModel(UUID studentId, UUID testId) {
        this.studentId = Objects.requireNonNull(studentId, "Student ID cannot be null");
        this.testId = Objects.requireNonNull(testId, "Test ID cannot be null");
    }

    /**
     * Конструктор для создания новой попытки теста.
     */
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

    /**
     * Конструктор для загрузки попытки из базы данных.
     */
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

    // ---------------------- SETTERS С ВАЛИДАЦИЕЙ ----------------------

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
     * @param finalPoint итоговый балл
     */
    public void completeAttempt(int finalPoint) {
        if (finalPoint < 0) {
            throw new IllegalArgumentException("Final point cannot be negative");
        }
        this.point = finalPoint;
        this.completed = true;
        if (this.dateOfAttempt == null) {
            this.dateOfAttempt = LocalDate.now();
        }
    }

    /**
     * Обновляет снапшот теста.
     */
    public void updateSnapshot(String newSnapshot) {
        this.attemptSnapshot = newSnapshot;
    }

    /**
     * Проверка валидности попытки перед сохранением.
     */
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

    /**
     * Проверяет, является ли попытка успешной (если есть минимальный балл).
     */
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

    /**
     * Укороченное строковое представление попытки.
     */
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