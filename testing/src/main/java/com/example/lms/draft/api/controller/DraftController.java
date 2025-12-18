package com.example.lms.draft.api.controller;

import com.example.lms.draft.api.dto.Draft;
import com.example.lms.draft.domain.service.DraftService;
import com.example.lms.tracing.SimpleTracer;

import io.javalin.http.Context;

import java.util.List;
import java.util.Map;
import java.util.NoSuchElementException;

/**
 * Контроллер для управления черновиками (draft).
 *
 * ВАЖНО: для удобства тестирования в Postman GET /drafts отдаёт JSON,
 * а не HTML-шаблон (в отличие от TestController).
 */
public class DraftController {

    /** Сервисный слой черновиков. */
    private final DraftService draftService;

    public DraftController(DraftService draftService) {
        this.draftService = draftService;
    }

    /**
     * GET /drafts
     * Возвращает список черновиков в JSON.
     */
    public void getDrafts(Context ctx) {
        SimpleTracer.runWithSpan("getDrafts", () -> {
            List<Draft> drafts = draftService.getAllDrafts();
            ctx.json(drafts);
        });
    }

    /**
     * POST /drafts
     * Создать черновик.
     */
    public void createDraft(Context ctx) {
        SimpleTracer.runWithSpan("createDraft", () -> {
            try {
                Draft dto = ctx.bodyAsClass(Draft.class);
                Draft created = draftService.createDraft(dto);
                ctx.status(201).json(created);
            } catch (IllegalArgumentException e) {
                // сюда обычно попадает "Invalid UUID string" и т.п.
                ctx.status(400).json(Map.of("error", e.getMessage()));
            }
        });
    }

    /**
     * GET /drafts/{id}
     * Получить черновик по ID.
     */
    public void getDraftById(Context ctx) {
        SimpleTracer.runWithSpan("getDraftById", () -> {
            try {
                String id = ctx.pathParam("id");
                Draft dto = draftService.getDraftById(id);
                ctx.json(dto);
            } catch (NoSuchElementException e) {
                ctx.status(404).json(Map.of("error", "Draft not found"));
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", e.getMessage()));
            }
        });
    }

    /**
     * GET /drafts/test/{testId}
     * Получить черновик по testId.
     */
    public void getDraftByTestId(Context ctx) {
        SimpleTracer.runWithSpan("getDraftByTestId", () -> {
            try {
                String testId = ctx.pathParam("testId");
                Draft dto = draftService.getDraftByTestId(testId);
                ctx.json(dto);
            } catch (NoSuchElementException e) {
                ctx.status(404).json(Map.of("error", "Draft not found"));
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", e.getMessage()));
            }
        });
    }

    /**
     * PUT /drafts/{id}
     * Обновить черновик по ID.
     */
    public void updateDraft(Context ctx) {
        SimpleTracer.runWithSpan("updateDraft", () -> {
            try {
                String id = ctx.pathParam("id");
                Draft dto = ctx.bodyAsClass(Draft.class);
                dto.setId(id);

                Draft updated = draftService.updateDraft(dto);
                ctx.json(updated);
            } catch (NoSuchElementException e) {
                ctx.status(404).json(Map.of("error", "Draft not found"));
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", e.getMessage()));
            }
        });
    }

    /**
     * DELETE /drafts/{id}
     * Удалить черновик по ID.
     */
    public void deleteDraft(Context ctx) {
        SimpleTracer.runWithSpan("deleteDraft", () -> {
            try {
                String id = ctx.pathParam("id");
                boolean deleted = draftService.deleteDraft(id);
                ctx.json(Map.of("deleted", deleted));
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", e.getMessage()));
            }
        });
    }
}