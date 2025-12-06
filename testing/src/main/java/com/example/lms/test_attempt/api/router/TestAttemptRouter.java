package com.example.lms.test_attempt.api.router;

import static io.javalin.apibuilder.ApiBuilder.*;

import com.example.lms.test_attempt.api.controller.TestAttemptController;

public class TestAttemptRouter {
    public static void register() {
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
