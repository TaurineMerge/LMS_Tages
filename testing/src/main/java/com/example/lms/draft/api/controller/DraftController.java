package com.example.lms.draft.api.controller;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;

import com.example.lms.draft.api.dto.Draft;
import com.example.lms.draft.domain.service.DraftService;
import com.example.lms.tracing.SimpleTracer;

import io.javalin.http.Context;

public class DraftController {
    private final DraftService draftService;

    public DraftController(DraftService draftService) {
        this.draftService = draftService;
    }

    /**
     * GET /drafts
     * Поддерживает фильтрацию по courseId через query parameter
     * Пример: /drafts?courseId=123e4567-e89b-12d3-a456-426614174000
     */
    public void getDrafts(Context ctx) {
        SimpleTracer.runWithSpan("getDrafts", () -> {
            try {
                String courseIdStr = ctx.queryParam("courseId");
                List<Draft> drafts;
                String status = "success";
                
                if (courseIdStr != null && !courseIdStr.trim().isEmpty()) {
                    UUID courseId = UUID.fromString(courseIdStr);
                    drafts = draftService.getDraftsByCourseId(courseId);
                } else {
                    drafts = draftService.getAllDrafts();
                }
                
                // Проверяем, есть ли черновики
                if (drafts == null || drafts.isEmpty()) {
                    status = "error";
                }
                
                Map<String, Object> response = new HashMap<>();
                response.put("data", drafts != null ? drafts : new ArrayList<>());
                response.put("courseId", courseIdStr);
                response.put("status", status);
                
                ctx.json(response);
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of(
                    "error", "Invalid UUID format for courseId",
                    "details", e.getMessage()
                ));
            } catch (Exception e) {
                ctx.status(500).json(Map.of(
                    "error", "Internal server error",
                    "details", e.getMessage()
                ));
            }
        });
    }

    /**
     * GET /drafts/by-course/{courseId}
     */
    public void getDraftsByCourseId(Context ctx) {
        SimpleTracer.runWithSpan("getDraftsByCourseId", () -> {
            try {
                String courseIdStr = ctx.pathParam("courseId");
                UUID courseId = UUID.fromString(courseIdStr);
                
                List<Draft> drafts = draftService.getDraftsByCourseId(courseId);
                ctx.json(drafts);
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", "Invalid UUID format"));
            }
        });
    }

    /**
     * GET /drafts/by-course/validate/{courseId}
     * Проверяет существование черновиков для курса
     */
    public void validateDraftByCourseId(Context ctx) {
        SimpleTracer.runWithSpan("validateDraftByCourseId", () -> {
            try {
                String courseIdStr = ctx.pathParam("courseId");
                UUID courseId = UUID.fromString(courseIdStr);
                
                List<Draft> drafts = draftService.getDraftsByCourseId(courseId);
                boolean exists = !drafts.isEmpty();
                
                ctx.json(Map.of("exists", exists));
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", "Invalid UUID format"));
            }
        });
    }

    /**
     * DELETE /drafts/by-course/{courseId}
     */
    public void deleteDraftsByCourseId(Context ctx) {
        SimpleTracer.runWithSpan("deleteDraftsByCourseId", () -> {
            try {
                String courseIdStr = ctx.pathParam("courseId");
                UUID courseId = UUID.fromString(courseIdStr);
                
                boolean deleted = draftService.deleteDraftsByCourseId(courseId);
                ctx.json(Map.of("deleted", deleted));
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", "Invalid UUID format"));
            }
        });
    }

    // Остальные существующие методы остаются без изменений
    public void createDraft(Context ctx) {
        SimpleTracer.runWithSpan("createDraft", () -> {
            Draft dto = ctx.bodyAsClass(Draft.class);
            if (dto.getId() != null) {
                ctx.status(400).json(Map.of("error", "ID should not be provided for new draft"));
                return;
            }
            
            Draft created = draftService.createDraft(dto);
            ctx.status(201).json(created);
        });
    }

    public void getDraftById(Context ctx) {
        SimpleTracer.runWithSpan("getDraftById", () -> {
            try {
                String idStr = ctx.pathParam("id");
                UUID id = UUID.fromString(idStr);
                
                Draft dto = draftService.getDraftById(id);
                if (dto == null) {
                    ctx.status(404).json(Map.of("error", "Draft not found"));
                } else {
                    ctx.json(dto);
                }
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", "Invalid UUID format"));
            }
        });
    }

    public void getDraftByTestId(Context ctx) {
        SimpleTracer.runWithSpan("getDraftByTestId", () -> {
            try {
                String testIdStr = ctx.pathParam("testId");
                UUID testId = UUID.fromString(testIdStr);
                
                Draft dto = draftService.getDraftByTestId(testId);
                if (dto == null) {
                    ctx.status(404).json(Map.of("error", "Draft not found for test"));
                } else {
                    ctx.json(dto);
                }
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", "Invalid UUID format"));
            }
        });
    }

    public void updateDraft(Context ctx) {
        SimpleTracer.runWithSpan("updateDraft", () -> {
            try {
                String idStr = ctx.pathParam("id");
                UUID id = UUID.fromString(idStr);
                
                Draft dto = ctx.bodyAsClass(Draft.class);
                dto.setId(id);
                
                Draft updated = draftService.updateDraft(dto);
                ctx.json(updated);
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", "Invalid UUID format"));
            }
        });
    }

    public void deleteDraft(Context ctx) {
        SimpleTracer.runWithSpan("deleteDraft", () -> {
            try {
                String idStr = ctx.pathParam("id");
                UUID id = UUID.fromString(idStr);
                
                boolean deleted = draftService.deleteDraft(id);
                ctx.json(Map.of("deleted", deleted));
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", "Invalid UUID format"));
            }
        });
    }
}