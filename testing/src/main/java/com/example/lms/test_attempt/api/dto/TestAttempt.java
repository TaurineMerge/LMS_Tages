package com.example.lms.test_attempt.api.dto;

import java.io.Serializable;
import java.sql.Date;

/**
 * DTO для представления попытки прохождения теста.
 * <p>
 * Соответствует таблице test_attempt_b:
 * <ul>
 *     <li><b>id</b> – идентификатор попытки (формально PK по комбинации student_id + test_id + date_of_attempt)</li>
 *     <li><b>date_of_attempt</b> – дата прохождения теста</li>
 *     <li><b>point</b> – количество набранных баллов</li>
 *     <li><b>result</b> – текстовый результат (например, JSON или строка-статус)</li>
 * </ul>
 *
 * Данный DTO используется для отображения данных о попытках тестирования
 * в API, сервисах и контроллерах.
 */
public class TestAttempt implements Serializable {

    private static final long serialVersionUID = 1L;

    /**
     * Идентификатор попытки.
     * <p>
     * В таблице TEST_ATTEMPT_B PK составной,
     * однако здесь используется Long для удобства передачи данных.
     */
    private Long id;

    /**
     * Дата прохождения теста.
     */
    private Date date_of_attempt;

    /**
     * Количество баллов, набранных за тест.
     */
    private Integer point;

    /**
     * Результат прохождения теста в строковом виде.
     * Может быть JSON или статус.
     */
    private String result;

    /**
     * Пустой конструктор для сериализации/десериализации.
     */
    public TestAttempt() {
    }

    /**
     * Основной конструктор DTO.
     *
     * @param id              идентификатор попытки
     * @param date_of_attempt дата прохождения теста
     * @param point           количество баллов
     * @param result          текстовый результат
     */
    public TestAttempt(Long id, Date date_of_attempt, Integer point, String result) {
        this.id = id;
        this.date_of_attempt = date_of_attempt;
        this.point = point;
        this.result = result;
    }

    /** @return идентификатор попытки */
    public Long getId() {
        return id;
    }

    /** @param id новый идентификатор попытки */
    public void setId(Long id) {
        this.id = id;
    }

    /** @return дата попытки */
    public Date getDate_of_attempt() {
        return date_of_attempt;
    }

    /**
     * Устанавливает дату прохождения теста.
     *
     * @param date_of_attempt новая дата попытки
     */
    public void setDate_of_attempt(Date date_of_attempt) {
        this.date_of_attempt = date_of_attempt;
    }

    /** @return набранные баллы */
    public Integer getPoint() {
        return point;
    }

    /** @param point новые баллы */
    public void setPoint(Integer point) {
        this.point = point;
    }

    /** @return результат попытки */
    public String getResult() {
        return result;
    }

    /** @param result новый текст результата */
    public void setResult(String result) {
        this.result = result;
    }

    /**
     * Формирует строковое представление объекта.
     *
     * @return строка с полями объекта
     */
    @Override
    public String toString() {
        return "TestAttempt{" +
                "id=" + id +
                ", date_of_attempt=" + date_of_attempt +
                ", point=" + point +
                ", result='" + result + '\'' +
                '}';
    }
}