package com.example.lms;

import io.javalin.Javalin;

import static io.javalin.apibuilder.ApiBuilder.*;

import com.example.lms.test_attempt.api.controller.TestAttemptController;
import com.example.lms.test.api.controller.TestController;

public class Main {

    public static void main(String[] args) {
        var app = Javalin.create(config -> {
            config.router.apiBuilder(() -> {
                registerTestRoutes();
                registerTestAttemptRoutes();
            });
        }).start(8085);
    }

    // Группировка эндпоинтов для Tests
    private static void registerTestRoutes() {
        path("/tests", () -> {
            get(TestController::getTests);
            post(TestController::createTest);

            path("/{id}", () -> {
                get(TestController::getTestById);
                put(TestController::updateTest);
                delete(TestController::deleteTest);
            });
        });
    }

    // Группировка эндпоинтов для Test Attempts
    private static void registerTestAttemptRoutes() {
        path("/test-attempts", () -> {
            get(TestAttemptController::getTestAttempts);
            post(TestAttemptController::createTestAttempt);

            path("/{id}", () -> {
                get(TestAttemptController::getTestAttemptById);
                delete(TestAttemptController::deleteTestAttempt);
            });
        });
    }
}