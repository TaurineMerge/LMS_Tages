package com.example.lms.internal.api.dto;

import com.fasterxml.jackson.annotation.JsonProperty;
import java.util.UUID;

/**
 * DTO для списка попыток.
 * Соответствует JSON-схеме AttemptsList (массив элементов).
 */
public class AttemptsListItem {

	@JsonProperty("attempt_id")
	private UUID attemptId;

	@JsonProperty("test_id")
	private UUID testId;

	@JsonProperty("date_of_attempt")
	private String dateOfAttempt;

	@JsonProperty("point")
	private Integer point;

	@JsonProperty("is_completed")
	private Boolean isCompleted;

	@JsonProperty("passed")
	private Boolean passed;

	@JsonProperty("certificate_id")
	private UUID certificateId;

	@JsonProperty("attempt_snapshot_s3")
	private String attemptSnapshot;

	// Constructors
	public AttemptsListItem() {
	}

	public AttemptsListItem(UUID attemptId, UUID testId, String dateOfAttempt,
			Integer point, Boolean isCompleted, Boolean passed, UUID certificateId,
			String attemptSnapshot) {
		this.attemptId = attemptId;
		this.testId = testId;
		this.dateOfAttempt = dateOfAttempt;
		this.point = point;
		this.isCompleted = isCompleted;
		this.passed = passed;
		this.certificateId = certificateId;
		this.attemptSnapshot = attemptSnapshot;
	}

	// Getters and Setters
	public UUID getAttemptId() {
		return attemptId;
	}

	public void setAttemptId(UUID attemptId) {
		this.attemptId = attemptId;
	}

	public UUID getTestId() {
		return testId;
	}

	public void setTestId(UUID testId) {
		this.testId = testId;
	}

	public String getDateOfAttempt() {
		return dateOfAttempt;
	}

	public void setDateOfAttempt(String dateOfAttempt) {
		this.dateOfAttempt = dateOfAttempt;
	}

	public Integer getPoint() {
		return point;
	}

	public void setPoint(Integer point) {
		this.point = point;
	}

	public Boolean getCompleted() {
		return isCompleted;
	}

	public void setCompleted(Boolean completed) {
		this.isCompleted = completed;
	}

	public Boolean getPassed() {
		return passed;
	}

	public void setPassed(Boolean passed) {
		this.passed = passed;
	}

	public UUID getCertificateId() {
		return null;
	}

	public void setCertificateId(UUID certificateId) {
		this.certificateId = certificateId;
	}

	public String getAttemptSnapshot() {
		return attemptSnapshot;
	}

	public void setAttemptSnapshot(String attemptSnapshot) {
		this.attemptSnapshot = attemptSnapshot;
	}
}
