package com.example.lms.test_attempt.api.dto;

import java.io.Serializable;

/**
 * DTO для представления попытки прохождения теста.
 */
public class TestAttempt implements Serializable {
    private static final long serialVersionUID = 1L;
    
    private String id;
    private String student_id;
    private String test_id;
    private String date_of_attempt; // ИЗМЕНЕНО: было LocalDate, стало String
    private Integer point;
    private String certificate_id;
    private String attempt_version;
    private String attempt_snapshot;
    private Boolean completed;

    public TestAttempt() {
    }

    public TestAttempt(String id, String student_id, String test_id, String date_of_attempt,
                      Integer point, String certificate_id, String attempt_version,
                      String attempt_snapshot, Boolean completed) {
        this.id = id;
        this.student_id = student_id;
        this.test_id = test_id;
        this.date_of_attempt = date_of_attempt;
        this.point = point;
        this.certificate_id = certificate_id;
        this.attempt_version = attempt_version;
        this.attempt_snapshot = attempt_snapshot;
        this.completed = completed != null ? completed : false;
    }

    // Getters and Setters
    public String getId() { return id; }
    public void setId(String id) { this.id = id; }
    
    public String getStudent_id() { return student_id; }
    public void setStudent_id(String student_id) { this.student_id = student_id; }
    
    public String getTest_id() { return test_id; }
    public void setTest_id(String test_id) { this.test_id = test_id; }
    
    public String getDate_of_attempt() { return date_of_attempt; }
    public void setDate_of_attempt(String date_of_attempt) { this.date_of_attempt = date_of_attempt; }
    
    public Integer getPoint() { return point; }
    public void setPoint(Integer point) { this.point = point; }
    
    public String getCertificate_id() { return certificate_id; }
    public void setCertificate_id(String certificate_id) { this.certificate_id = certificate_id; }
    
    public String getAttempt_version() { return attempt_version; }
    public void setAttempt_version(String attempt_version) { this.attempt_version = attempt_version; }
    
    public String getAttempt_snapshot() { return attempt_snapshot; }
    public void setAttempt_snapshot(String attempt_snapshot) { this.attempt_snapshot = attempt_snapshot; }
    
    public Boolean getCompleted() { return completed; }
    public void setCompleted(Boolean completed) { 
        this.completed = completed != null ? completed : false; 
    }

    @Override
    public String toString() {
        return "TestAttempt{" +
                "id='" + id + '\'' +
                ", student_id='" + student_id + '\'' +
                ", test_id='" + test_id + '\'' +
                ", date_of_attempt='" + date_of_attempt + '\'' +
                ", point=" + point +
                ", certificate_id='" + certificate_id + '\'' +
                ", completed=" + completed +
                '}';
    }
}