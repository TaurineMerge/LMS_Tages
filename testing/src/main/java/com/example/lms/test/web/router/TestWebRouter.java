package com.example.lms.test.web.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.test.web.controller.TestFormController;
import com.example.lms.tracing.SimpleTracer;

import static io.javalin.apibuilder.ApiBuilder.after;
import static io.javalin.apibuilder.ApiBuilder.before;
import static io.javalin.apibuilder.ApiBuilder.delete;
import static io.javalin.apibuilder.ApiBuilder.get;
import static io.javalin.apibuilder.ApiBuilder.patch;
import static io.javalin.apibuilder.ApiBuilder.path;
import static io.javalin.apibuilder.ApiBuilder.post;
import static io.javalin.apibuilder.ApiBuilder.put;

/**
 * Роутер для веб-интерфейса (HTML формы) тестов
 */
public class TestWebRouter {
    private static final Logger logger = LoggerFactory.getLogger(TestWebRouter.class);
    
    public static void register(TestFormController controller) {
        path("/web", () -> {
            get("/", controller::showHomePage);
            
            path("/tests", () -> {
                before(ctx -> {
                    logger.info("Web request started: {} {} (traceId: {})",
                            ctx.method(),
                            ctx.path(),
                            SimpleTracer.getCurrentTraceId());
                });
                
                // ========== GET ЗАПРОСЫ (ОТОБРАЖЕНИЕ) ==========
                
                // 1. Форма создания нового теста
                get("/new", controller::showNewTestForm);
                
                // ========== POST ЗАПРОСЫ (СОЗДАНИЕ) ==========
                
                // 2. Создание нового теста (через универсальную форму)
                post("/save", controller::saveTestFromForm);
                
                // 3. Создание черновика нового теста
                post("/draft", controller::saveNewTestDraft);
                
                // ========== ОПЕРАЦИИ С КОНКРЕТНЫМ ТЕСТОМ/ЧЕРНОВИКОМ ==========
                
                path("/{id}", () -> {
                    // 4. Универсальная форма редактирования (теста или черновика)
                    get("/edit", controller::showEditTestForm);
                    
                    // 5. Предпросмотр теста
                    get("/preview", controller::previewTest);
                    
                    // 6. Обновление существующего теста/черновика (PUT - полное обновление)
                    put("/", ctx -> {
                        String id = ctx.pathParam("id");
                        if (id.startsWith("draft-")) {
                            controller.updateDraftFromForm(ctx);
                        } else {
                            controller.updateTestFromForm(ctx);
                        }
                    });
                    
                    // 7. Частичное обновление (PATCH - если нужно)
                    patch("/", ctx -> {
                        String id = ctx.pathParam("id");
                        if (id.startsWith("draft-")) {
                            // controller.partialUpdateDraft(ctx); // Пока нет реализации
                            ctx.status(501).result("Частичное обновление черновика не реализовано");
                        } else {
                            // controller.partialUpdateTest(ctx); // Пока нет реализации
                            ctx.status(501).result("Частичное обновление теста не реализовано");
                        }
                    });
                    
                    // 8. Удаление теста/черновика
                    delete("/", ctx -> {
                        String id = ctx.pathParam("id");
                        if (id.startsWith("draft-")) {
                            controller.deleteDraft(ctx);
                        } else {
                            controller.deleteTest(ctx);
                        }
                    });
                    post("/delete", ctx -> {
                        String id = ctx.pathParam("id");
                        if (id.startsWith("draft-")) {
                            controller.deleteDraft(ctx);
                        } else {
                            controller.deleteTest(ctx);
                        }
                    });
                    // ========== ДОПОЛНИТЕЛЬНЫЕ ОПЕРАЦИИ ==========
                    
                    // 9. Создание черновика из существующего теста
                    post("/create-draft", controller::createDraftFromTest);
                    
                    // 10. Сохранение черновика существующего теста
                    post("/draft", controller::saveExistingTestDraft);
                    
                    // 11. Публикация черновика
                    post("/publish", controller::publishDraft);
                });
                
                after(ctx -> {
                    logger.info("Web request completed: {} {} -> {} (traceId: {})",
                            ctx.method(),
                            ctx.path(),
                            ctx.status(),
                            SimpleTracer.getCurrentTraceId());
                });
            });
        });
    }
}