package com.example.lms.test.api.router;

import com.example.lms.test.api.controller.TestController;
import com.example.lms.security.JwtHandler;
import static io.javalin.apibuilder.ApiBuilder.*;

public class TestRouter {
    public static void register() {
        path("/tests", () -> {
            before(JwtHandler.authenticate());

            get(TestController::getTests);
            post(TestController::createTest);

            path("/{id}", () -> {
                get(TestController::getTestById);
                put(TestController::updateTest);
                delete(TestController::deleteTest);
            });
        });
    }
}