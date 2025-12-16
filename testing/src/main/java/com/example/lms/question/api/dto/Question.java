package com.example.lms.question.api.dto;

import java.util.UUID;

/**
 * DTO для передачи данных о вопросе теста.
 * <p>
 * Используется на уровне API (контроллеры, сериализация/десериализация JSON)
 * и не содержит бизнес-логики — только структуру данных.
 * <br>
 * Соответствует сущности QUESTION_D в базе данных:
 * <ul>
 * <li>{@code id} — идентификатор вопроса</li>
 * <li>{@code testId} — идентификатор теста, которому принадлежит вопрос</li>
 * <li>{@code textOfQuestion} — текст вопроса</li>
 * <li>{@code order} — позиция вопроса внутри теста</li>
 * </ul>
 */
public class Question {

	/** Идентификатор вопроса. Может быть null при создании нового вопроса. */
	private UUID id;

	/** ID теста, к которому относится данный вопрос. */
	private UUID testId;

	/** Текст вопроса, отображаемый студенту. */
	private String textOfQuestion;

	/** Порядковый номер вопроса в тесте. */
	private Integer order;

	/**
	 * Пустой конструктор — требуется для сериализации/десериализации JSON.
	 */
	public Question() {
	}

	/**
	 * Основной конструктор DTO.
	 *
	 * @param id             идентификатор вопроса
	 * @param testId         идентификатор теста
	 * @param textOfQuestion текст вопроса
	 * @param order          порядковый номер вопроса
	 */
	public Question(UUID id, UUID testId, String textOfQuestion, Integer order) {
		this.id = id;
		this.testId = testId;
		this.textOfQuestion = textOfQuestion;
		this.order = order;
	}

	/** @return идентификатор вопроса */
	public UUID getId() {
		return id;
	}

	/** @param id новый идентификатор вопроса */
	public void setId(UUID id) {
		this.id = id;
	}

	/** @return идентификатор теста */
	public UUID getTestId() {
		return testId;
	}

	/** @param testId новый идентификатор теста */
	public void setTestId(UUID testId) {
		this.testId = testId;
	}

	/** @return текст вопроса */
	public String getTextOfQuestion() {
		return textOfQuestion;
	}

	/** @param textOfQuestion новый текст вопроса */
	public void setTextOfQuestion(String textOfQuestion) {
		this.textOfQuestion = textOfQuestion;
	}

	/** @return порядковый номер вопроса */
	public Integer getOrder() {
		return order;
	}

	/** @param order новый порядковый номер вопроса */
	public void setOrder(Integer order) {
		this.order = order;
	}
}