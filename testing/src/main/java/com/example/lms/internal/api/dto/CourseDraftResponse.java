package com.example.lms.internal.api.dto;

import java.io.Serializable;
import java.util.UUID;

/**
 * DTO для ответа на запрос черновика по courseId.
 */
public class CourseDraftResponse implements Serializable {

	private static final long serialVersionUID = 1L;

	/**
	 * Данные черновика.
	 */
	private DraftData data;

	/**
	 * ID курса.
	 */
	private UUID courseId;

	/**
	 * Статус ответа.
	 */
	private String status;

	public CourseDraftResponse() {
	}

	public CourseDraftResponse(DraftData data, UUID courseId, String status) {
		this.data = data;
		this.courseId = courseId;
		this.status = status;
	}

	public DraftData getData() {
		return data;
	}

	public void setData(DraftData data) {
		this.data = data;
	}

	public UUID getCourseId() {
		return courseId;
	}

	public void setCourseId(UUID courseId) {
		this.courseId = courseId;
	}

	public String getStatus() {
		return status;
	}

	public void setStatus(String status) {
		this.status = status;
	}

	/**
	 * Внутренний класс для данных черновика.
	 */
	public static class DraftData implements Serializable {

		private static final long serialVersionUID = 1L;

		private UUID id;
		private UUID testId;
		private UUID courseId;
		private String title;
		private Integer min_point;
		private String description;

		public DraftData() {
		}

		public DraftData(UUID id, UUID testId, UUID courseId, String title, Integer min_point, String description) {
			this.id = id;
			this.testId = testId;
			this.courseId = courseId;
			this.title = title;
			this.min_point = min_point;
			this.description = description;
		}

		public UUID getId() {
			return id;
		}

		public void setId(UUID id) {
			this.id = id;
		}

		public UUID getTestId() {
			return testId;
		}

		public void setTestId(UUID testId) {
			this.testId = testId;
		}

		public UUID getCourseId() {
			return courseId;
		}

		public void setCourseId(UUID courseId) {
			this.courseId = courseId;
		}

		public String getTitle() {
			return title;
		}

		public void setTitle(String title) {
			this.title = title;
		}

		public Integer getMin_point() {
			return min_point;
		}

		public void setMin_point(Integer min_point) {
			this.min_point = min_point;
		}

		public String getDescription() {
			return description;
		}

		public void setDescription(String description) {
			this.description = description;
		}
	}
}
