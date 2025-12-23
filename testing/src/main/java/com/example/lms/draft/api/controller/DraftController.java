package com.example.lms.draft.api.controller;

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
     */
    public void getDrafts(Context ctx) {
        SimpleTracer.runWithSpan("getDrafts", () -> {
            List<Draft> drafts = draftService.getAllDrafts();
            ctx.json(drafts);
        });
    }

    /**
     * POST /drafts
     */
    public void createDraft(Context ctx) {
        SimpleTracer.runWithSpan("createDraft", () -> {
            Draft dto = ctx.bodyAsClass(Draft.class);
            // Проверяем, что ID не установлен (будет сгенерирован в БД)
            if (dto.getId() != null) {
                ctx.status(400).json(Map.of("error", "ID should not be provided for new draft"));
                return;
            }
            
            Draft created = draftService.createDraft(dto);
            ctx.status(201).json(created);
        });
    }

    /**
     * GET /drafts/{id}
     */
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

    /**
     * GET /drafts/test/{testId}
     */
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

    /**
     * PUT /drafts/{id}
     */
    public void updateDraft(Context ctx) {
        SimpleTracer.runWithSpan("updateDraft", () -> {
            try {
                String idStr = ctx.pathParam("id");
                UUID id = UUID.fromString(idStr);
                
                Draft dto = ctx.bodyAsClass(Draft.class);
                dto.setId(id); // Устанавливаем ID из пути
                
                Draft updated = draftService.updateDraft(dto);
                ctx.json(updated);
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", "Invalid UUID format"));
            }
        });
    }

    /**
     * DELETE /drafts/{id}
     */
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