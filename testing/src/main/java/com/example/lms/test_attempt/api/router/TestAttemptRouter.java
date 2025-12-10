package com.example.lms.test_attempt.api.router;

import static io.javalin.apibuilder.ApiBuilder.*;

import com.example.lms.security.JwtHandler;
import com.example.lms.test_attempt.api.controller.TestAttemptController;

/**
 * Router для работы с попытками прохождения тестов.
 * <p>
 * Регистрирует REST-эндпоинты:
 * <ul>
 *     <li>GET  /test-attempts — список попыток</li>
 *     <li>POST /test-attempts — создание новой попытки</li>
 *     <li>GET  /test-attempts/{id} — получение попытки по ID</li>
 *     <li>DELETE /test-attempts/{id} — удаление попытки</li>
 * </ul>
 *
 * Все маршруты защищены JWT-аутентификацией через {@link JwtHandler#authenticate()}.
 * <p>
 * Метод {@link #register()} вызывается при инициализации приложения
 * (например, в {@code Main} классе через {@code config.router.apiBuilder(...)}).
 */
public class TestAttemptRouter {

    /**
     * Регистрирует маршруты для ресурса {@code /test-attempts}.
     * <p>
     * Структура:
     * <pre>
     * /test-attempts
     *    GET     — список попыток
     *    POST    — создание попытки
     *
     * /test-attempts/{id}
     *    GET     — получить попытку
     *    DELETE  — удалить попытку
     * </pre>
     */
    public static void register() {
        path("/test-attempts", () -> {

            // Глобальный фильтр для всех эндпоинтов /test-attempts — проверка JWT
            before(JwtHandler.authenticate());

            // Основные операции
            get(TestAttemptController::getTestAttempts);
            post(TestAttemptController::createTestAttempt);

            // Операции над конкретной попыткой
            path("/{id}", () -> {
                get(TestAttemptController::getTestAttemptById);
                delete(TestAttemptController::deleteTestAttempt);
            });
        });
    }
}