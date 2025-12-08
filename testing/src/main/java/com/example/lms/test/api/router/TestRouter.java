package com.example.lms.test.api.router;

import com.example.lms.test.api.controller.TestController;
import com.example.lms.security.JwtHandler;
import static io.javalin.apibuilder.ApiBuilder.*;

public class TestRouter {
    public static void register(TestController testController) {
        path("/tests", () -> {
            before(JwtHandler.authenticate());

            get(testController::getTests);
            post(testController::createTest);

            path("/{id}", () -> {
                get(testController::getTestById);
                put(testController::updateTest);
                delete(testController::deleteTest);
            });
        });
    }
}