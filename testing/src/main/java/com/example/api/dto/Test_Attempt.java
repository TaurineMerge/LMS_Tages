package com.example.api.dto;
import java.io.Serializable;
import java.sql.Date;

//import javax.xml.crypto.Data;

public class Test_Attempt implements Serializable {

    private static final long serialVersionUID = 1L; // Для сериализации

    private Long id;
    private Date date_of_attempt;
    private Integer point;
    private String result;

    public Test_Attempt() {
    }

    public Test_Attempt(Long id, Date date_of_attempt, Integer point, String result) {
        this.id = id;
        this.date_of_attempt = date_of_attempt;
        this.point = point;
        this.result = result;
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
    }

    public Date getDate_of_attempt() {
        return date_of_attempt;
    }

    public void setTitle(Date date_of_attempt) {
        this.date_of_attempt = date_of_attempt;
    }

    public Integer getPoint() {
        return point;
    }

    public void setPoint(Integer point) {
        this.point = point;
    }

    public String getResult() {
        return result;
    }

    public void setResult(String result) {
        this.result = result;
    }

    @Override
    public String toString() {
        return "CreateTestRequestDto{" +
                "id =" + id +
                ", date of attempt ='" + date_of_attempt + '\'' +
                ", point ='" + point + '\'' +
                ", result ='" + result + '\'' +
                '}';
    }
}