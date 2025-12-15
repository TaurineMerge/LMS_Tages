package com.example.lms.test_attempt.api.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.security.JwtHandler;
import com.example.lms.test_attempt.api.controller.TestAttemptController;

import static io.javalin.apibuilder.ApiBuilder.before;
import static io.javalin.apibuilder.ApiBuilder.delete;
import static io.javalin.apibuilder.ApiBuilder.get;
import static io.javalin.apibuilder.ApiBuilder.path;
import static io.javalin.apibuilder.ApiBuilder.post;
import static io.javalin.apibuilder.ApiBuilder.put;

/**
 * Router для работы с попытками прохождения тестов.
 */
public class TestAttemptRouter {

    private static final Logger logger = LoggerFactory.getLogger(TestAttemptRouter.class);
    
    /**
     * Регистрирует маршруты для ресурса {@code /test-attempts}.
     */
    public static void register(TestAttemptController testAttemptController) {
        path("/test-attempts", () -> {

            // Глобальный фильтр для всех эндпоинтов /test-attempts — проверка JWT
            before(JwtHandler.authenticate());

            // Основные CRUD операции
            get(testAttemptController::getAllTestAttempts);
            post(testAttemptController::createTestAttempt);

            // Специальные запросы
            get("/completed", testAttemptController::getCompletedTestAttempts);
            get("/incomplete", testAttemptController::getIncompleteTestAttempts);

            // Операции над конкретной попыткой
            path("/{id}", () -> {
                get(testAttemptController::getTestAttemptById);
                put(testAttemptController::updateTestAttempt);
                delete(testAttemptController::deleteTestAttempt);
            });

            // Поиск по критериям
            path("/student/{studentId}", () -> {
                get(testAttemptController::getTestAttemptsByStudentId);
            });

            path("/test/{testId}", () -> {
                get(testAttemptController::getTestAttemptsByTestId);
            });

            path("/date/{date}", () -> {
                get(testAttemptController::getTestAttemptsByDate);
            });

            // Дополнительные операции для связки студент-тест
            path("/student/{studentId}/test/{testId}", () -> {
                get(testAttemptController::getAttemptsByStudentAndTest);
                get("/best", testAttemptController::getBestAttemptByStudentAndTest);
                get("/count", testAttemptController::countAttemptsByStudentAndTest);
            });

            // Проверка существования
            get("/exists/{id}", testAttemptController::existsById);
        });
        
        logger.info("Маршруты TestAttemptRouter успешно зарегистрированы");
    }
}