package com.example.lms.domain.model.test;

import java.time.LocalDateTime;
import java.util.UUID;

public class Test {
    private UUID id;
    private UUID courseId;
    private String title;
    private Integer minPoint;
    private String description;
    private LocalDateTime createdAt;
    private LocalDateTime updatedAt;

    public Test(UUID courseId, String title, Integer minPoint, String description) {
        this.courseId = courseId;
        this.title = title;
        this.minPoint = minPoint != null ? minPoint : 0;
        this.description = description;
    }
}