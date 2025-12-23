package com.example.lms.question.api.dto;

import java.util.UUID;

/**
 * DTO для передачи данных о вопросе теста или черновика.
 * <p>
 * Используется на уровне API (контроллеры, сериализация/десериализация JSON)
 * и не содержит бизнес-логики — только структуру данных.
 * <br>
 * Соответствует сущности QUESTION_D в базе данных:
 * <ul>
 * <li>{@code id} — идентификатор вопроса</li>
 * <li>{@code testId} — идентификатор теста, которому принадлежит вопрос (может быть null для черновиков)</li>
 * <li>{@code draftId} — идентификатор черновика, которому принадлежит вопрос (может быть null для тестов)</li>
 * <li>{@code textOfQuestion} — текст вопроса</li>
 * <li>{@code order} — позиция вопроса внутри теста/черновика</li>
 * </ul>
 */
public class Question {

	/** Идентификатор вопроса. Может быть null при создании нового вопроса. */
	private UUID id;

	/** ID теста, к которому относится данный вопрос. Может быть null для черновиков. */
	private UUID testId;

	/** ID черновика, к которому относится данный вопрос. Может быть null для тестов. */
	private UUID draftId;

	/** Текст вопроса, отображаемый студенту. */
	private String textOfQuestion;

	/** Порядковый номер вопроса в тесте/черновике. */
	private Integer order;

	/**
	 * Пустой конструктор — требуется для сериализации/десериализации JSON.
	 */
	public Question() {
	}

	/**
	 * Основной конструктор DTO (для обратной совместимости).
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

	/**
	 * Полный конструктор DTO с поддержкой черновиков.
	 *
	 * @param id             идентификатор вопроса
	 * @param testId         идентификатор теста
	 * @param draftId        идентификатор черновика
	 * @param textOfQuestion текст вопроса
	 * @param order          порядковый номер вопроса
	 */
	public Question(UUID id, UUID testId, UUID draftId, String textOfQuestion, Integer order) {
		this.id = id;
		this.testId = testId;
		this.draftId = draftId;
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

	/** @return идентификатор черновика */
	public UUID getDraftId() {
		return draftId;
	}

	/** @param draftId новый идентификатор черновика */
	public void setDraftId(UUID draftId) {
		this.draftId = draftId;
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

	/**
	 * Проверяет, принадлежит ли вопрос тесту.
	 *
	 * @return true если вопрос принадлежит тесту, false если черновику или ничему
	 */
	public boolean isForTest() {
		return testId != null;
	}

	/**
	 * Проверяет, принадлежит ли вопрос черновику.
	 *
	 * @return true если вопрос принадлежит черновику, false если тесту или ничему
	 */
	public boolean isForDraft() {
		return draftId != null;
	}

	/**
	 * Возвращает идентификатор родительской сущности (теста или черновика).
	 *
	 * @return UUID теста или черновика, или null если оба null
	 */
	public UUID getParentId() {
		if (testId != null) {
			return testId;
		}
		return draftId;
	}

	@Override
	public String toString() {
		return "Question{" +
				"id=" + id +
				", testId=" + testId +
				", draftId=" + draftId +
				", order=" + order +
				", text='" + (textOfQuestion != null ? 
					textOfQuestion.length() > 30 ? 
						textOfQuestion.substring(0, 27) + "..." : 
						textOfQuestion : "null") + '\'' +
				'}';
	}
}