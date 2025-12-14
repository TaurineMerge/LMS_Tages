package com.example.lms.test.api.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.security.JwtHandler;
import com.example.lms.test.api.controller.TestController;
import com.example.lms.tracing.SimpleTracer;

import io.github.cdimascio.dotenv.Dotenv;

import static io.javalin.apibuilder.ApiBuilder.*;

import java.util.Set;

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
    private static final Dotenv dotenv = Dotenv.load();

    private static final String KEYCLOAK_STUDENT_REALM = dotenv.get("KEYCLOAK_STUDENT_REALM");
    private static final String KEYCLOAK_TEACHER_REALM = dotenv.get("KEYCLOAK_TEACHER_REALM");

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
            before(JwtHandler.authenticate());

            // ---- MIDDLEWARE: Логирование начала запроса ----
            before(ctx -> {
                logger.info("Request started: {} {} (traceId: {})",
                        ctx.method(),
                        ctx.path(),
                        SimpleTracer.getCurrentTraceId());
            });

            // ---- CRUD маршруты ----
            get(ctx -> { // GET /tests
                JwtHandler.requireRealm(Set.of(KEYCLOAK_TEACHER_REALM, KEYCLOAK_STUDENT_REALM));
                testController.getTests(ctx);
            });
            post(ctx -> { // POST /tests
                JwtHandler.requireRealm(KEYCLOAK_TEACHER_REALM);
                testController.createTest(ctx);
            });

            path("/{id}", () -> {
                get(ctx -> { // GET /tests/{id}
                    JwtHandler.requireRealm(Set.of(KEYCLOAK_TEACHER_REALM, KEYCLOAK_STUDENT_REALM));
                    testController.getTestById(ctx);
                });
                put(ctx -> { // PUT /tests/{id}
                    JwtHandler.requireRealm(KEYCLOAK_TEACHER_REALM);
                    testController.updateTest(ctx);
                });
                delete(ctx -> { // DELETE /tests/{id}
                    JwtHandler.requireRealm(KEYCLOAK_TEACHER_REALM);
                    testController.deleteTest(ctx);
                });
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