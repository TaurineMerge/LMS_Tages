package com.example.lms.test_attempt.api.dto;

import java.io.Serializable;
import java.sql.Date;
import java.util.UUID;

/**
 * DTO для передачи данных о попытке прохождения теста через API.
 * <p>
 * Используется для:
 * <ul>
 * <li>ответов контроллера</li>
 * <li>получения данных от клиента</li>
 * </ul>
 *
 * Это упрощённое представление доменной модели
 * {@code TestAttemptModel}, предназначенное только для коммуникации
 * по сети. DTO не содержит бизнес-логики.
 */
public class TestAttempt implements Serializable {

	private static final long serialVersionUID = 1L;

	/** Уникальный идентификатор попытки. */
	private UUID id;

	/** Дата прохождения попытки (java.sql.Date для совместимости с БД). */
	private Date date_of_attempt;

	/** Количество набранных баллов; может быть null, если попытка не завершена. */
	private Integer point;

	/**
	 * Строковое представление результата (например: "passed", "failed",
	 * "incomplete").
	 */
	private String result;

	/**
	 * Пустой конструктор необходим для сериализации и десериализации
	 * (Jackson/Javalin).
	 */
	public TestAttempt() {
	}

	/**
	 * Полный конструктор DTO.
	 *
	 * @param id              идентификатор попытки
	 * @param date_of_attempt дата прохождения
	 * @param point           набранные баллы
	 * @param result          строковый статус результата
	 */
	public TestAttempt(UUID id, Date date_of_attempt, Integer point, String result) {
		this.id = id;
		this.date_of_attempt = date_of_attempt;
		this.point = point;
		this.result = result;
	}

	// --------------------------------------------------------------------
	// GETTERS / SETTERS
	// --------------------------------------------------------------------

	public UUID getId() {
		return id;
	}

	public void setId(UUID id) {
		this.id = id;
	}

	/** @return дата попытки */
	public Date getDate_of_attempt() {
		return date_of_attempt;
	}

	/**
	 * Устанавливает дату прохождения попытки.
	 *
	 * @param date дата попытки
	 */
	public void setDate_of_attempt(Date date) {
		this.date_of_attempt = date;
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

	// --------------------------------------------------------------------
	// toString()
	// --------------------------------------------------------------------

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