package com.example.api.dto;
import java.io.Serializable;

public class CreateTestRequestDto implements Serializable {

    private static final long serialVersionUID = 1L; // Для сериализации

    private Long id;
    private String title;
    private Integer min_point;
    private String description;

    public CreateTestRequestDto() {
    }

    public CreateTestRequestDto(Long id, String title, Integer min_point, String description) {
        this.id = id;
        this.title = title;
        this.min_point = min_point;
        this.description = description;
    }

    public Long getId() {
        return id;
    }

    public void setId(Long id) {
        this.id = id;
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

    @Override
    public String toString() {
        return "CreateTestRequestDto{" +
                "id=" + id +
                ", title='" + title + '\'' +
                ", lastName='" + min_point + '\'' +
                ", description='" + description + '\'' +
                '}';
    }
}