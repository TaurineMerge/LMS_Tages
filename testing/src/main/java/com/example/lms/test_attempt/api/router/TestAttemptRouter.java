package com.example.lms.test_attempt.api.router;

import com.example.lms.test_attempt.api.controller.TestAttemptController;
import com.example.lms.security.JwtHandler;
import static io.javalin.apibuilder.ApiBuilder.*;

public class TestAttemptRouter {
    public static void register() {
        path("/test-attempts", () -> {
            before(JwtHandler.authenticate()); // JWT проверка

            get(TestAttemptController::getAttempts);
            post(TestAttemptController::createAttempt);
            // остальные роуты...
        });
    }
}