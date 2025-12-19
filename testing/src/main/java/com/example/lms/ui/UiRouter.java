package com.example.lms.ui;

import static io.javalin.apibuilder.ApiBuilder.get;
import static io.javalin.apibuilder.ApiBuilder.post;
import static io.javalin.apibuilder.ApiBuilder.path;

public class UiRouter {

    public static void register(UiTestController uiTestController) {
        path("/ui", () -> {

            get("/tests/{testId}/take", uiTestController::showTakePage);

            post("/tests/{testId}/questions/{questionId}/answer", uiTestController::submitAnswer);

            get("/tests/{testId}/finish", uiTestController::showFinishPage);
            post("/tests/{testId}/finish", uiTestController::finishAttempt);

            get("/tests/{testId}/results", uiTestController::showResultsPage);
        });
    }
}