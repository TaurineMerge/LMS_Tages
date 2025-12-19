package com.example.lms.test.web.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.test.web.controller.TestFormController;
import com.example.lms.tracing.SimpleTracer;

import static io.javalin.apibuilder.ApiBuilder.after;
import static io.javalin.apibuilder.ApiBuilder.before;
import static io.javalin.apibuilder.ApiBuilder.get;
import static io.javalin.apibuilder.ApiBuilder.path;
import static io.javalin.apibuilder.ApiBuilder.post;

/**
 * Роутер для веб-интерфейса (HTML формы) тестов
 * Этот роутер НЕ использует JWT аутентификацию для простоты
 * Можно добавить сессионную аутентификацию позже
 */
public class TestWebRouter {
    private static final Logger logger = LoggerFactory.getLogger(TestWebRouter.class);
    
    /**
     * Регистрирует маршруты для веб-интерфейса тестов
     * Все маршруты начинаются с /web/tests
     */
    public static void register(TestFormController controller) {
        // ВЕБ-ИНТЕРФЕЙС (HTML)
        path("/web", () -> {
            // Главная страница (перенаправление на создание теста)
            get("/", controller::showHomePage);
            
            // Тесты
            path("/tests", () -> {
                // Логирование начала запроса
                before(ctx -> {
                    logger.info("Web request started: {} {} (traceId: {})",
                            ctx.method(),
                            ctx.path(),
                            SimpleTracer.getCurrentTraceId());
                });
                
                // Создание нового теста
                get("/new", controller::showNewTestForm);
                
                // Сохранение теста
                post("/save", controller::saveTestFromForm);
                
                // Операции с конкретным тестом
                path("/{id}", () -> {
                    // Редактирование теста
                    get("/edit", controller::showEditTestForm);
                    
                    // Сохранение черновика
                    post("/draft", controller::saveTestDraft);
                    
                    // Предпросмотр теста
                    get("/preview", controller::previewTest);
                });
                
                // Логирование завершения запроса
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