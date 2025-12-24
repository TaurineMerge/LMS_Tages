package com.example.lms.draft.api.router;

import com.example.lms.draft.api.controller.DraftController;
import com.example.lms.security.JwtHandler;

import static io.javalin.apibuilder.ApiBuilder.*;

/**
 * Router для работы с черновиками тестов (draft).
 */
public class DraftRouter {

    public static void register(DraftController draftController) {
        path("/drafts", () -> {
            before(JwtHandler.authenticate());

            // ---------------------------
            // READ: список черновиков с фильтрацией по courseId
            // ---------------------------
            
            /**
             * GET /drafts
             * Поддерживает фильтрацию по courseId через query parameter
             * Если courseId не указан, возвращает все черновики
             */
            get(draftController::getDrafts);
            
            /**
             * GET /drafts/by-course/{courseId}
             * Получить черновики по ID курса (path parameter)
             */
            get("/by-course/{courseId}", draftController::getDraftsByCourseId);

            /**
             * GET /drafts/by-course/validate/{courseId}
             * Проверить существование черновиков для курса
             */
            get("/by-course/validate/{courseId}", draftController::validateDraftByCourseId);

            // ---------------------------
            // CREATE: создать черновик
            // ---------------------------
            post(draftController::createDraft);

            // ---------------------------
            // READ: получить по ID
            // ---------------------------
            get("/{id}", draftController::getDraftById);

            // ---------------------------
            // READ: получить по testId
            // ---------------------------
            get("/test/{testId}", draftController::getDraftByTestId);

            // ---------------------------
            // UPDATE: обновить черновик
            // ---------------------------
            put("/{id}", draftController::updateDraft);

            // ---------------------------
            // DELETE: удалить черновик
            // ---------------------------
            delete("/{id}", draftController::deleteDraft);
            
            // ---------------------------
            // DELETE: удалить по courseId
            // ---------------------------
            delete("/by-course/{courseId}", draftController::deleteDraftsByCourseId);
        });
    }
}