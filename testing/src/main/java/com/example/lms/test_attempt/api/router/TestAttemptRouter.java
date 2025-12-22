package com.example.lms.test_attempt.api.router;

import java.util.Set;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.shared.router.RouterUtils;
import static com.example.lms.shared.router.RouterUtils.READ_ACCESS_REALMS;
import static com.example.lms.shared.router.RouterUtils.TEACHER_REALM;
import static com.example.lms.shared.router.RouterUtils.STUDENT_REALM;
import static com.example.lms.shared.router.RouterUtils.applyStandardAfterMiddleware;
import static com.example.lms.shared.router.RouterUtils.applyStandardBeforeMiddleware;
import static com.example.lms.shared.router.RouterUtils.validateController;
import static com.example.lms.shared.router.RouterUtils.withRealm;
import com.example.lms.test_attempt.api.controller.TestAttemptController;

import static io.javalin.apibuilder.ApiBuilder.delete;
import static io.javalin.apibuilder.ApiBuilder.get;
import static io.javalin.apibuilder.ApiBuilder.path;
import static io.javalin.apibuilder.ApiBuilder.post;
import static io.javalin.apibuilder.ApiBuilder.put;

/**
 * Роутер для управления попытками прохождения тестов.
 * <p>
 * Регистрирует REST-эндпоинты:
 * <ul>
 * <li>GET /test-attempts — получить список попыток (HTML)</li>
 * <li>POST /test-attempts — создать новую попытку</li>
 * <li>GET /test-attempts/{id} — получить попытку по ID</li>
 * <li>PUT /test-attempts/{id} — обновить попытку</li>
 * <li>DELETE /test-attempts/{id} — удалить попытку</li>
 * <li>POST /test-attempts/{id}/complete — завершить попытку</li>
 * <li>PUT /test-attempts/{id}/snapshot — обновить снапшот</li>
 * <li>GET /test-attempts/student/{studentId} — получить попытки студента</li>
 * <li>GET /test-attempts/test/{testId} — получить попытки по тесту</li>
 * </ul>
 * 
 * <h2>Политика доступа:</h2>
 * <ul>
 * <li><b>Студенты:</b> доступ только к своим попыткам</li>
 * <li><b>Преподаватели:</b> полный доступ ко всем попыткам</li>
 * </ul>
 * 
 * @see TestAttemptController
 * @see RouterUtils
 */
public class TestAttemptRouter {
    private static final Logger logger = LoggerFactory.getLogger(TestAttemptRouter.class);

    /**
     * Регистрирует маршруты группы /test-attempts и их подмаршрутов.
     * 
     * @param testAttemptController контроллер, содержащий обработчики запросов
     * @throws IllegalArgumentException если testAttemptController равен null
     */
    public static void register(TestAttemptController testAttemptController) {
        validateController(testAttemptController, "TestAttemptController");

        path("/test-attempts", () -> {
            // Стандартные middleware
            applyStandardBeforeMiddleware(logger);

            // Список и создание попыток
            get(withRealm(Set.of(TEACHER_REALM, STUDENT_REALM), testAttemptController::getTestAttempts));
            post(withRealm(Set.of(TEACHER_REALM, STUDENT_REALM), testAttemptController::createTestAttempt));

            path("/{id}", () -> {
                // Просмотр попытки
                get(withRealm(Set.of(TEACHER_REALM, STUDENT_REALM), testAttemptController::getTestAttemptById));

                // Редактирование и удаление
                put(withRealm(TEACHER_REALM, testAttemptController::updateTestAttempt));
                delete(withRealm(TEACHER_REALM, testAttemptController::deleteTestAttempt));

                // Специальные операции
                post("/complete",
                        withRealm(Set.of(TEACHER_REALM, STUDENT_REALM), testAttemptController::completeTestAttempt));
                put("/snapshot",
                        withRealm(Set.of(TEACHER_REALM, STUDENT_REALM), testAttemptController::updateSnapshot));
            });

            // Поиск по студенту и тесту
            path("/student/{studentId}", () -> {
                get(withRealm(Set.of(TEACHER_REALM, STUDENT_REALM), testAttemptController::getTestAttemptsByStudentId));
            });

            path("/test/{testId}", () -> {
                get(withRealm(Set.of(TEACHER_REALM), testAttemptController::getTestAttemptsByTestId));
            });

            applyStandardAfterMiddleware(logger);
        });

        // Health check without auth - bypass authentication completely
        get("/test-attempts/health", ctx -> {
            // Skip authentication for health check
            ctx.result("OK");
        });
    }
}