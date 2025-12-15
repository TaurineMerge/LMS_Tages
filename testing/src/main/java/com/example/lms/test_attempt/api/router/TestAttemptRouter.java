package com.example.lms.test_attempt.api.router;

import static io.javalin.apibuilder.ApiBuilder.*;

import com.example.lms.security.JwtHandler;
import com.example.lms.test_attempt.api.controller.TestAttemptController;
import com.example.lms.tracing.SimpleTracer;
import io.javalin.http.Context;

/**
 * Router для работы с ресурсом TestAttempt.
 * <p>
 * Регистрирует HTTP-маршруты для CRUD-операций над попытками прохождения тестов.
 * Все методы обёрнуты в обработку трассировки (OpenTelemetry) через {@link SimpleTracer},
 * а также используют JWT-аутентификацию через {@link JwtHandler}.
 *
 * Основные эндпоинты:
 * <ul>
 *     <li>GET /test-attempts — получить список попыток</li>
 *     <li>POST /test-attempts — создать новую попытку</li>
 *     <li>GET /test-attempts/{id} — получить попытку по ID</li>
 *     <li>DELETE /test-attempts/{id} — удалить попытку</li>
 * </ul>
 *
 * Каждый запрос автоматически:
 * <ul>
 *     <li>записывает атрибуты пользователя в Span</li>
 *     <li>логирует параметры запроса</li>
 *     <li>создаёт события трассировки (events)</li>
 * </ul>
 */
public class TestAttemptRouter {

    /**
     * Регистрирует все маршруты для работы с /test-attempts.
     * <p>
     * Вызывается из основного класса приложения (Main.java)
     * внутри config.router.apiBuilder(...).
     */
    public static void register() {
        path("/test-attempts", () -> {

            // Обязательная JWT-аутентификация перед каждым запросом
            before(JwtHandler.authenticate());

            /**
             * GET /test-attempts
             * Получение списка попыток с опциональными query-параметрами:
             * userId, testId, status.
             */
            get(ctx -> {
                SimpleTracer.runWithSpan("getTestAttempts", () -> {
                    captureUserAttributes(ctx);

                    // Логируем query-параметры, если они есть
                    String userId = ctx.queryParam("userId");
                    String testId = ctx.queryParam("testId");
                    String status = ctx.queryParam("status");

                    if (userId != null) SimpleTracer.addAttribute("query.userId", userId);
                    if (testId != null) SimpleTracer.addAttribute("query.testId", testId);
                    if (status != null) SimpleTracer.addAttribute("query.status", status);

                    // Контроллер-заглушка
                    TestAttemptController.getTestAttempts(ctx);

                    SimpleTracer.addEvent("test_attempts.retrieved");
                });
            });

            /**
             * POST /test-attempts
             * Создание новой попытки тестирования.
             * Логируются данные запроса: тип и размер контента.
             */
            post(ctx -> {
                SimpleTracer.runWithSpan("createTestAttempt", () -> {
                    captureUserAttributes(ctx);

                    SimpleTracer.addAttribute("content.type", ctx.contentType());
                    SimpleTracer.addAttribute("content.length", String.valueOf(ctx.contentLength()));

                    TestAttemptController.createTestAttempt(ctx);

                    SimpleTracer.addEvent("test_attempt.created");
                });
            });

            path("/{id}", () -> {

                /**
                 * GET /test-attempts/{id}
                 * Получение информации о попытке по ID.
                 */
                get(ctx -> {
                    SimpleTracer.runWithSpan("getTestAttemptById", () -> {
                        captureUserAttributes(ctx);

                        String attemptId = ctx.pathParam("id");
                        SimpleTracer.addAttribute("test_attempt.id", attemptId);

                        TestAttemptController.getTestAttemptById(ctx);

                        SimpleTracer.addEvent("test_attempt.retrieved.by.id");
                    });
                });

                /**
                 * DELETE /test-attempts/{id}
                 * Удаление попытки тестирования.
                 */
                delete(ctx -> {
                    SimpleTracer.runWithSpan("deleteTestAttempt", () -> {
                        captureUserAttributes(ctx);

                        String attemptId = ctx.pathParam("id");
                        SimpleTracer.addAttribute("test_attempt.id", attemptId);

                        TestAttemptController.deleteTestAttempt(ctx);

                        SimpleTracer.addEvent("test_attempt.deleted");
                    });
                });
            });
        });
    }

    /**
     * Извлекает пользовательские атрибуты из контекста Javalin
     * и добавляет их в текущий Span OpenTelemetry.
     * <p>
     * Атрибуты:
     * <ul>
     *     <li>user.id</li>
     *     <li>user.username</li>
     *     <li>user.email</li>
     *     <li>user.roles (если реализовано)</li>
     * </ul>
     *
     * @param ctx текущий HTTP-контекст Javalin
     */
    private static void captureUserAttributes(Context ctx) {
        Object userId = ctx.attribute("userId");
        Object username = ctx.attribute("username");
        Object email = ctx.attribute("email");
        Object roles = ctx.attribute("roles");

        if (userId != null) {
            SimpleTracer.addAttribute("user.id", userId.toString());
        }
        if (username != null) {
            SimpleTracer.addAttribute("user.username", username.toString());
        }
        if (email != null) {
            SimpleTracer.addAttribute("user.email", email.toString());
        }
        if (roles != null) {
            SimpleTracer.addAttribute("user.roles", roles.toString());
        }
    }
}