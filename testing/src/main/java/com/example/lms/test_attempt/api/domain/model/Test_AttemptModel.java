package com.example.lms.test_attempt.api.domain.model;

import java.time.LocalDate;
import java.util.Objects;
import java.util.UUID;

/**
 * Модель попытки прохождения теста (test_attempt_b в БД)
 * Соответствует таблице test_attempt_b
 */
public class Test_AttemptModel {
    private UUID id;
    private UUID studentId;          // student_id в БД
    private UUID testId;             // test_id в БД
    private LocalDate dateOfAttempt; // date_of_attempt в БД
    private Integer point;           // может быть null если тест не завершен
    private UUID certificateId;      // certificate_id в БД (может быть null)
    private String attemptVersion;   // attempt_version как JSON в БД, у нас String
    
    // Конструктор для создания новой попытки
    public Test_AttemptModel(UUID studentId, UUID testId, String attemptVersion) {
        this.studentId = Objects.requireNonNull(studentId, "Student ID cannot be null");
        this.testId = Objects.requireNonNull(testId, "Test ID cannot be null");
        this.dateOfAttempt = LocalDate.now();
        this.attemptVersion = attemptVersion;
        // point и certificateId могут быть null изначально
    }
    
    // Конструктор для загрузки из БД
    public Test_AttemptModel(UUID id, UUID studentId, UUID testId, 
                           LocalDate dateOfAttempt, Integer point, 
                           UUID certificateId, String attemptVersion) {
        this.id = id;
        this.studentId = studentId;
        this.testId = testId;
        this.dateOfAttempt = dateOfAttempt;
        this.point = point;
        this.certificateId = certificateId;
        this.attemptVersion = attemptVersion;
    }
    
    // Доменные методы (бизнес-логика)
    
    /**
     * Завершить попытку теста с указанным количеством баллов
     * @param points количество набранных баллов
     * @throws IllegalStateException если попытка уже завершена
     */
    public void complete(int points) {
        if (this.point != null) {
            throw new IllegalStateException("Test attempt is already completed");
        }
        if (points < 0) {
            throw new IllegalArgumentException("Points cannot be negative");
        }
        this.point = points;
    }
    
    /**
     * Привязать сертификат к попытке
     * @param certificateId ID выданного сертификата
     * @throws IllegalStateException если попытка не завершена
     */
    public void attachCertificate(UUID certificateId) {
        if (this.point == null) {
            throw new IllegalStateException("Cannot attach certificate to incomplete attempt");
        }
        this.certificateId = certificateId;
    }
    
    /**
     * Проверить, завершена ли попытка
     */
    public boolean isCompleted() {
        return point != null;
    }
    
    /**
     * Проверить, пройден ли тест (если установлен минимальный балл в тесте)
     * @param minPoints минимальный балл для прохождения (может быть null)
     * @return true если тест пройден
     */
    public boolean isPassed(Integer minPoints) {
        if (point == null) {
            return false; // Тест не завершен
        }
        if (minPoints == null) {
            return true; // Если нет минимального балла, любой результат = пройден
        }
        return point >= minPoints;
    }
    
    // Геттеры
    public UUID getId() { return id; }
    public UUID getStudentId() { return studentId; }
    public UUID getTestId() { return testId; }
    public LocalDate getDateOfAttempt() { return dateOfAttempt; }
    public Integer getPoint() { return point; }
    public UUID getCertificateId() { return certificateId; }
    public String getAttemptVersion() { return attemptVersion; }
    
    // Сеттеры с валидацией
    public void setId(UUID id) {
        this.id = Objects.requireNonNull(id, "ID cannot be null");
    }
    
    public void setPoint(Integer point) {
        if (point != null && point < 0) {
            throw new IllegalArgumentException("Points cannot be negative");
        }
        this.point = point;
    }
    
    public void setCertificateId(UUID certificateId) {
        this.certificateId = certificateId;
    }
    
    public void setAttemptVersion(String attemptVersion) {
        this.attemptVersion = attemptVersion;
    }
    
    /**
     * Валидация модели перед сохранением
     */
    public void validate() {
        if (studentId == null) {
            throw new IllegalArgumentException("Student ID cannot be null");
        }
        if (testId == null) {
            throw new IllegalArgumentException("Test ID cannot be null");
        }
        if (dateOfAttempt == null) {
            throw new IllegalArgumentException("Date of attempt cannot be null");
        }
        if (dateOfAttempt.isAfter(LocalDate.now())) {
            throw new IllegalArgumentException("Date of attempt cannot be in the future");
        }
        if (point != null && point < 0) {
            throw new IllegalArgumentException("Points cannot be negative");
        }
    }
    
    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        Test_AttemptModel that = (Test_AttemptModel) o;
        return Objects.equals(id, that.id);
    }
    
    @Override
    public int hashCode() {
        return Objects.hash(id);
    }
    
    @Override
    public String toString() {
        return "Test_AttemptModel{" +
                "id=" + id +
                ", studentId=" + studentId +
                ", testId=" + testId +
                ", dateOfAttempt=" + dateOfAttempt +
                ", point=" + point +
                ", completed=" + isCompleted() +
                '}';
    }
}