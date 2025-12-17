package com.example.lms.test.api.dto;

import java.io.Serializable;

/**
 * DTO для представления теста.
 * <p>
 * Используется для передачи данных о тестах между слоями приложения.
 * Содержит информацию:
 * <ul>
 * <li><b>id</b> — идентификатор теста</li>
 * <li><b>title</b> — название теста</li>
 * <li><b>min_point</b> — минимальный проходной балл</li>
 * <li><b>description</b> — описание теста</li>
 * </ul>
 *
 * Реализует {@link Serializable}, что позволяет передавать объект
 * через сеть или сохранять в файлы, если это необходимо.
 */
public class Test implements Serializable {

	private static final long serialVersionUID = 1L;

	/**
	 * Идентификатор теста.
	 */
	private String id;

	/**
	 * Название теста.
	 */
	private String title;

	/**
	 * Минимальный балл, необходимый для прохождения теста.
	 */
	private Integer min_point;

	/**
	 * Текстовое описание теста.
	 */
	private String description;

	/**
	 * Пустой конструктор для сериализации/десериализации.
	 */
	public Test() {
	}

	/**
	 * Основной конструктор.
	 *
	 * @param id          идентификатор теста
	 * @param title       название теста
	 * @param min_point   минимальный проходной балл
	 * @param description описание теста
	 */
	public Test(String id, String title, Integer min_point, String description) {
		this.id = id;
		this.title = title;
		this.min_point = min_point;
		this.description = description;
	}

	/** @return идентификатор теста */
	public String getId() {
		return id;
	}

	/** @param id новый идентификатор теста */
	public void setId(String id) {
		this.id = id;
	}

	/** @return название теста */
	public String getTitle() {
		return title;
	}

	/** @param title новое название теста */
	public void setTitle(String title) {
		this.title = title;
	}

	/** @return минимальный проходной балл */
	public Integer getMin_point() {
		return min_point;
	}

	/** @param min_point новый минимальный проходной балл */
	public void setMin_point(Integer min_point) {
		this.min_point = min_point;
	}

	/** @return описание теста */
	public String getDescription() {
		return description;
	}

	/** @param description новое описание теста */
	public void setDescription(String description) {
		this.description = description;
	}

	/**
	 * Формирует строковое представление объекта Test.
	 *
	 * @return строка с полями объекта
	 */
	@Override
	public String toString() {
		return "Test{" +
				"id=" + id +
				", title='" + title + '\'' +
				", min_point=" + min_point +
				", description='" + description + '\'' +
				'}';
	}
}