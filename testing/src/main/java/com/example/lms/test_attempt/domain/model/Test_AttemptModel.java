package com.example.lms.test_attempt.domain.model;

import java.time.LocalDate;
import java.util.Objects;
import java.util.UUID;

/**
 * Доменная модель попытки прохождения теста.
 * <p>
 * Соответствует таблице <b>test_attempt_b</b> в базе данных.
 * Хранит информацию о:
 * <ul>
 *     <li>ID попытки</li>
 *     <li>ID студента</li>
 *     <li>ID теста</li>
 *     <li>дате прохождения</li>
 *     <li>результате (баллы)</li>
 * </ul>
 *
 * Модель содержит бизнес-методы для:
 * <ul>
 *     <li>завершения попытки</li>
 *     <li>проверки завершённости и успешности теста</li>
 *     <li>валидации перед сохранением</li>
 * </ul>
 */
public class Test_AttemptModel {

    private UUID id;
    private UUID studentId;
    private UUID testId;
    private LocalDate dateOfAttempt;
    private Integer point;

    /**
     * Конструктор для создания новой попытки.
     * <p>
     * Используется на уровне сервиса при создании новой записи.
     * Значения:
     * <ul>
     *     <li>{@code dateOfAttempt} устанавливается как текущая дата</li>
     *     <li>{@code point} изначально равен null</li>
     * </ul>
     *
     * @param studentId идентификатор студента
     * @param testId    идентификатор теста
     */
    public Test_AttemptModel(UUID studentId, UUID testId) {
        this.studentId = Objects.requireNonNull(studentId, "Student ID cannot be null");
        this.testId = Objects.requireNonNull(testId, "Test ID cannot be null");
        this.dateOfAttempt = LocalDate.now();
        this.point = null;
    }

    public Test_AttemptModel() {
        // Пустой конструктор для десериализации
    }

    /**
     * Конструктор для загрузки данных из БД.
     *
     * @param id            идентификатор попытки
     * @param studentId     идентификатор студента
     * @param testId        идентификатор теста
     * @param dateOfAttempt дата прохождения
     * @param point         набранные баллы (null, если не завершено)
     */
    public Test_AttemptModel(UUID id, UUID studentId, UUID testId,
                             LocalDate dateOfAttempt, Integer point) {
        this.id = id;
        this.studentId = studentId;
        this.testId = testId;
        this.dateOfAttempt = dateOfAttempt;
        this.point = point;
    }

    // ----------------------------------------------------------------------
    //                       BUSINESS LOGIC METHODS
    // ----------------------------------------------------------------------

    /**
     * Завершает попытку теста, устанавливая количество набранных баллов.
     *
     * @param points баллы за прохождение теста
     * @throws IllegalStateException    если попытка уже завершена
     * @throws IllegalArgumentException если points отрицательные
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
     * Проверяет, завершена ли попытка.
     *
     * @return true — если есть значение point, иначе false
     */
    public boolean isCompleted() {
        return point != null;
    }

    /**
     * Проверяет, пройден ли тест.
     *
     * @param minPoints минимальный балл, требуемый для прохождения теста (nullable)
     * @return true — если тест завершён и баллов достаточно
     */
    public boolean isPassed(Integer minPoints) {
        if (point == null) {
            return false;
        }
        if (minPoints == null) {
            return true;
        }
        return point >= minPoints;
    }

    // ----------------------------------------------------------------------
    //                              GETTERS
    // ----------------------------------------------------------------------

    public UUID getId() { return id; }
    public UUID getStudentId() { return studentId; }
    public UUID getTestId() { return testId; }
    public LocalDate getDateOfAttempt() { return dateOfAttempt; }
    public Integer getPoint() { return point; }

    // ----------------------------------------------------------------------
    //                              SETTERS
    // ----------------------------------------------------------------------

    public void setId(UUID id) {
        this.id = Objects.requireNonNull(id, "ID cannot be null");
    }

    public void setPoint(Integer point) {
        if (point != null && point < 0) {
            throw new IllegalArgumentException("Points cannot be negative");
        }
        this.point = point;
    }

    /**
     * Устанавливает дату попытки.
     * 
     * @param dateOfAttempt дата попытки
     */
    public void setDateOfAttempt(LocalDate dateOfAttempt) {
        this.dateOfAttempt = dateOfAttempt;
    }

    // ----------------------------------------------------------------------
    //                              VALIDATION
    // ----------------------------------------------------------------------

    /**
     * Проверяет корректность данных модели перед сохранением.
     *
     * @throws IllegalArgumentException если модель содержит недопустимые значения
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

    // ----------------------------------------------------------------------
    //                          OVERRIDDEN METHODS
    // ----------------------------------------------------------------------

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