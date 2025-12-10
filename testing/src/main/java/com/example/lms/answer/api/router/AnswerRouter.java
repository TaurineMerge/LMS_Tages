package com.example.lms.answer.api.router;

import com.example.lms.answer.api.controller.AnswerController;
import com.example.lms.tracing.SimpleTracer;
import com.example.lms.Main;
import com.example.lms.security.JwtHandler;
import static io.javalin.apibuilder.ApiBuilder.*;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Роутер для маршрутов управления ответами.
 * Определяет конечные точки REST API для операций с ответами.
 */
public class AnswerRouter {
    private static final Logger logger = LoggerFactory.getLogger(Main.class);

    /**
     * Регистрирует маршруты для управления ответами.
     *
     * @param answerController контроллер для обработки запросов
     */
    public static void register(AnswerController answerController) {
        path("/answers", () -> {
            before(JwtHandler.authenticate());

            before(ctx -> {
                logger.info("▶️  Request started: {} {} (traceId: {})",
                        ctx.method(),
                        ctx.path(),
                        SimpleTracer.getCurrentTraceId());
            });


            post(answerController::createAnswer);                   // POST /answers
            
            // Операции с конкретным ответом
            path("/{id}", () -> {
                get(answerController::getAnswerById);               // GET /answers/{id}
                put(answerController::updateAnswer);                // PUT /answers/{id}
                delete(answerController::deleteAnswer);             // DELETE /answers/{id}
            });

            // Операции с ответами по вопросам
            path("/question", () -> {
                get("/by-question", answerController::getAnswersByQuestionId);           // GET /answers/question/by-question?questionId={id}
                get("/correct", answerController::getCorrectAnswersByQuestionId);        // GET /answers/question/correct?questionId={id}
                delete("/delete-all", answerController::deleteAnswersByQuestionId);      // DELETE /answers/question/delete-all?questionId={id}
                get("/count", answerController::countAnswersByQuestionId);               // GET /answers/question/count?questionId={id}
                get("/count-correct", answerController::countCorrectAnswersByQuestionId); // GET /answers/question/count-correct?questionId={id}
            });

            after(ctx -> {
                logger.info("✅ Request completed: {} {} -> {} (traceId: {})",
                        ctx.method(),
                        ctx.path(),
                        ctx.status(),
                        SimpleTracer.getCurrentTraceId());
            });
        });
    }
}