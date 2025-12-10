package com.example.lms.test_attempt.api.domain.model;

import java.time.LocalDate;
import java.util.Objects;
import java.util.UUID;

/**
 * Доменная модель: TestAttemptModel
 *
 * Представляет строку таблицы test_attempt_b.
 *
 * Поля:
 *  - id:              UUID  — первичный ключ попытки
 *  - student_id:      UUID  — идентификатор студента
 *  - test_id:         UUID  — идентификатор теста
 *  - date_of_attempt: Date  — дата прохождения теста
 *  - point:           int   — количество набранных баллов (может быть NULL)
 *  - certificate_id:  UUID  — сертификат, выданный за прохождение (NULL, если не выдан)
 *  - attempt_version: JSON  — структура попытки (хранится строкой)
 */
public class TestAttemptModel {

    /** Уникальный идентификатор попытки (PRIMARY KEY). */
    private UUID id;

    /** Идентификатор студента, проходившего тест. */
    private UUID studentId;

    /** Идентификатор теста. */
    private UUID testId;

    /** Дата прохождения теста. */
    private LocalDate dateOfAttempt;

    /** Количество баллов (null — если тест не завершён). */
    private Integer point;

    /** Идентификатор сертификата (если выдан). */
    private UUID certificateId;

    /** JSON-версия структуры попытки — ответы, шаги, всё, что нужно сохранить. */
    private String attemptVersion;


    // ─────────────────────────────────────────────
    //               КОНСТРУКТОРЫ
    // ─────────────────────────────────────────────

    /**
     * Конструктор для создания новой попытки (до сохранения в БД).
     */
    public TestAttemptModel(UUID studentId, UUID testId, String attemptVersion) {
        this.studentId = Objects.requireNonNull(studentId, "Student ID не может быть null");
        this.testId = Objects.requireNonNull(testId, "Test ID не может быть null");
        this.dateOfAttempt = LocalDate.now();
        this.attemptVersion = attemptVersion;
    }

    /**
     * Конструктор для загрузки попытки из БД.
     */
    public TestAttemptModel(
            UUID id,
            UUID studentId,
            UUID testId,
            LocalDate dateOfAttempt,
            Integer point,
            UUID certificateId,
            String attemptVersion
    ) {
        this.id = id;
        this.studentId = studentId;
        this.testId = testId;
        this.dateOfAttempt = dateOfAttempt;
        this.point = point;
        this.certificateId = certificateId;
        this.attemptVersion = attemptVersion;
    }


    // ─────────────────────────────────────────────
    //               БИЗНЕС-ЛОГИКА
    // ─────────────────────────────────────────────

    /**
     * Завершает попытку теста и устанавливает количество баллов.
     */
    public void complete(int points) {
        if (this.point != null) {
            throw new IllegalStateException("Попытка уже завершена");
        }
        if (points < 0) {
            throw new IllegalArgumentException("Баллы не могут быть отрицательными");
        }
        this.point = points;
    }

    /**
     * Привязывает сертификат к завершённой попытке.
     */
    public void attachCertificate(UUID certificateId) {
        if (this.point == null) {
            throw new IllegalStateException("Нельзя выдать сертификат незавершённой попытке");
        }
        this.certificateId = certificateId;
    }

    /**
     * Проверяет, завершена ли попытка.
     */
    public boolean isCompleted() {
        return point != null;
    }

    /**
     * Проверяет, пройден ли тест (если известно минимальное значение).
     */
    public boolean isPassed(Integer minPoints) {
        if (point == null) return false;     // попытка не завершена
        if (minPoints == null) return true;  // нет минимального балла — тест считается пройденным
        return point >= minPoints;
    }


    // ─────────────────────────────────────────────
    //               ГЕТТЕРЫ / СЕТТЕРЫ
    // ─────────────────────────────────────────────

    public UUID getId() { return id; }
    public UUID getStudentId() { return studentId; }
    public UUID getTestId() { return testId; }
    public LocalDate getDateOfAttempt() { return dateOfAttempt; }
    public Integer getPoint() { return point; }
    public UUID getCertificateId() { return certificateId; }
    public String getAttemptVersion() { return attemptVersion; }

    public void setId(UUID id) {
        this.id = Objects.requireNonNull(id, "ID не может быть null");
    }

    public void setPoint(Integer point) {
        if (point != null && point < 0) {
            throw new IllegalArgumentException("Баллы не могут быть отрицательными");
        }
        this.point = point;
    }

    public void setCertificateId(UUID certificateId) {
        this.certificateId = certificateId;
    }

    public void setAttemptVersion(String attemptVersion) {
        this.attemptVersion = attemptVersion;
    }


    // ─────────────────────────────────────────────
    //                  ВАЛИДАЦИЯ
    // ─────────────────────────────────────────────

    /**
     * Проверка корректности данных перед сохранением.
     */
    public void validate() {
        if (studentId == null) {
            throw new IllegalArgumentException("Student ID не может быть null");
        }
        if (testId == null) {
            throw new IllegalArgumentException("Test ID не может быть null");
        }
        if (dateOfAttempt == null) {
            throw new IllegalArgumentException("Дата попытки не может быть null");
        }
        if (dateOfAttempt.isAfter(LocalDate.now())) {
            throw new IllegalArgumentException("Дата попытки не может быть в будущем");
        }
        if (point != null && point < 0) {
            throw new IllegalArgumentException("Баллы не могут быть отрицательными");
        }
    }


    // ─────────────────────────────────────────────
    //                UTILS / LOGGING
    // ─────────────────────────────────────────────

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        TestAttemptModel that = (TestAttemptModel) o;
        return Objects.equals(id, that.id);
    }

    @Override
    public int hashCode() { return Objects.hash(id); }

    @Override
    public String toString() {
        return "TestAttemptModel{" +
                "id=" + id +
                ", studentId=" + studentId +
                ", testId=" + testId +
                ", dateOfAttempt=" + dateOfAttempt +
                ", point=" + point +
                ", completed=" + isCompleted() +
                '}';
    }
}