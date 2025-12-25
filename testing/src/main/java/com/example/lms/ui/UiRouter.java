package com.example.lms.ui;

import static io.javalin.apibuilder.ApiBuilder.get;
import static io.javalin.apibuilder.ApiBuilder.post;
import static io.javalin.apibuilder.ApiBuilder.path;

public class UiRouter {

    public static void register(UiTestController uiTestController) {
        path("/ui", () -> {

            // "умная" точка входа с курса
            // /ui/category/{categoryId}/course/{courseId}/test
            get("/category/{categoryId}/course/{courseId}/test", uiTestController::startOrResumeFromCourse);

            // открыть конкретную попытку с курса
            // /ui/category/{categoryId}/course/{courseId}/attempt/{attemptId}
            get("/category/{categoryId}/course/{courseId}/attempt/{attemptId}", uiTestController::openAttemptFromCourse);

            // attemptId идёт query-параметром:
            // /ui/tests/{testId}/take?attemptId=...
            get("/tests/{testId}/take", uiTestController::showTakePage);

            // submit: hidden attempt_id
            post("/tests/{testId}/questions/{questionId}/answer", uiTestController::submitAnswer);

            // finish/results тоже работают с attemptId:
            get("/tests/{testId}/finish", uiTestController::showFinishPage);
            post("/tests/{testId}/finish", uiTestController::finishAttempt);

            get("/tests/{testId}/results", uiTestController::showResultsPage);

            // "пройти ещё раз" — просто создаёт новую попытку и кидает на take
            post("/tests/{testId}/retry", uiTestController::startNewAttempt);
        });
    }
}
