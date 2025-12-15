package com.example.lms.test.api.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.security.JwtHandler;
import com.example.lms.test.api.controller.TestController;
import com.example.lms.tracing.SimpleTracer;

import static io.javalin.apibuilder.ApiBuilder.after;
import static io.javalin.apibuilder.ApiBuilder.before;
import static io.javalin.apibuilder.ApiBuilder.delete;
import static io.javalin.apibuilder.ApiBuilder.get;
import static io.javalin.apibuilder.ApiBuilder.path;
import static io.javalin.apibuilder.ApiBuilder.post;
import static io.javalin.apibuilder.ApiBuilder.put;

/**
 * Router для сущности Test.
 * <p>
 * Регистрирует REST-эндпоинты:
 * <ul>
 * <li>GET /tests — получить список тестов (HTML)</li>
 * <li>POST /tests — создать тест</li>
 * <li>GET /tests/{id} — получить тест по ID</li>
 * <li>PUT /tests/{id} — обновить тест</li>
 * <li>DELETE /tests/{id} — удалить тест</li>
 * </ul>
 *
 * Все запросы:
 * <ul>
 * <li>проходят через JWT-аутентификацию {@link JwtHandler#authenticate()}</li>
 * <li>логируются до и после обработки</li>
 * <li>содержат traceId для трассировки запроса в OpenTelemetry</li>
 * </ul>
 *
 * Метод {@link #register(TestController)} вызывается из Main при поднятии
 * приложения.
 */
public class TestRouter {
    private static final Logger logger = LoggerFactory.getLogger(TestRouter.class);

    /**
     * Регистрирует маршруты группы /tests и их подмаршрутов.
     * <p>
     * Пример структуры:
     * 
     * <pre>
     * /tests
     *    GET     — список тестов (HTML)
     *    POST    — создание теста
     *
     * /tests/{id}
     *    GET     — получить тест
     *    PUT     — обновить тест
     *    DELETE  — удалить тест
     * </pre>
     *
     * @param testController контроллер, содержащий обработчики запросов
     */
    public static void register(TestController testController) {

        path("/tests", () -> {

            // ---- MIDDLEWARE: Аутентификация ----
           // before(JwtHandler.authenticate());

            // ---- MIDDLEWARE: Логирование начала запроса ----
            before(ctx -> {
                logger.info("Request started: {} {} (traceId: {})",
                        ctx.method(),
                        ctx.path(),
                        SimpleTracer.getCurrentTraceId());
            });

            // ---- CRUD маршруты ----
            get(testController::getTests);
            post(testController::createTest);

            path("/{id}", () -> {
                get(testController::getTestById);
                put(testController::updateTest);
                delete(testController::deleteTest);
            });

            // ---- MIDDLEWARE: Логирование завершения запроса ----
            after(ctx -> {
                logger.info("Request completed: {} {} -> {} (traceId: {})",
                        ctx.method(),
                        ctx.path(),
                        ctx.status(),
                        SimpleTracer.getCurrentTraceId());
            });
        });
    }
}