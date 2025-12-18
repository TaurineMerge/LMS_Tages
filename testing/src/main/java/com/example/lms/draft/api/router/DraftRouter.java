package com.example.lms.draft.api.router;

import com.example.lms.draft.api.controller.DraftController;
import com.example.lms.security.JwtHandler;

import static io.javalin.apibuilder.ApiBuilder.*;

/**
 * Router для работы с черновиками тестов (draft).
 *
 * Здесь мы:
 * 1) задаём базовый путь /drafts
 * 2) вешаем before()-middleware для проверки JWT на ВСЕ эндпоинты /drafts
 * 3) регистрируем CRUD-роуты и привязываем их к методам DraftController
 */
public class DraftRouter {

    /**
     * Регистрирует маршруты для ресурса /drafts.
     *
     * @param draftController контроллер, который содержит обработчики HTTP-запросов
     */
    public static void register(DraftController draftController) {

        // Группируем все маршруты под общим префиксом /drafts
        path("/drafts", () -> {

            /**
             * before(...) — фильтр (middleware), который выполнится ПЕРЕД любым эндпоинтом внутри path("/drafts", ...).
             *
             * JwtHandler.authenticate() возвращает Javalin Handler:
             * Handler = (Context ctx) -> { ... }
             *
             * Этот handler:
             * - проверит токен (JWT)
             * - если токен невалидный / отсутствует — выбросит ошибку (обычно 401)
             * - если всё ок — пропустит запрос дальше к конкретному endpoint handler-у
             */
            before(JwtHandler.authenticate());

            // ---------------------------
            // READ: список черновиков
            // ---------------------------

            /**
             * GET /drafts
             * В твоём контроллере сейчас метод getDrafts рендерит HTML (если ты оставляешь Handlebars).
             * Если тебе нужен только JSON — можно заменить контроллерный метод.
             */
            get(draftController::getDrafts);

            // ---------------------------
            // CREATE: создать черновик
            // ---------------------------

            /**
             * POST /drafts
             * Ожидаем JSON в body → DraftController.createDraft
             */
            post(draftController::createDraft);

            // ---------------------------
            // READ: получить по ID
            // ---------------------------

            /**
             * GET /drafts/{id}
             * pathParam "id" должен быть UUID в виде строки
             */
            get("/{id}", draftController::getDraftById);

            // ---------------------------
            // READ: получить по testId
            // ---------------------------

            /**
             * GET /drafts/test/{testId}
             * Удобный эндпоинт: найти черновик по id теста
             */
            get("/test/{testId}", draftController::getDraftByTestId);

            // ---------------------------
            // UPDATE: обновить черновик
            // ---------------------------

            /**
             * PUT /drafts/{id}
             * ID берём из pathParam, а поля — из JSON body
             */
            put("/{id}", draftController::updateDraft);

            // ---------------------------
            // DELETE: удалить черновик
            // ---------------------------

            /**
             * DELETE /drafts/{id}
             * Возвращаем {"deleted": true/false}
             */
            delete("/{id}", draftController::deleteDraft);
        });
    }
}