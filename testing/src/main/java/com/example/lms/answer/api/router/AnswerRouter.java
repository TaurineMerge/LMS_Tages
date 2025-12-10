package com.example.lms.answer.api.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.Main;
import com.example.lms.answer.api.controller.AnswerController;
import com.example.lms.security.JwtHandler;
import com.example.lms.tracing.SimpleTracer;

import static io.javalin.apibuilder.ApiBuilder.after;
import static io.javalin.apibuilder.ApiBuilder.before;
import static io.javalin.apibuilder.ApiBuilder.delete;
import static io.javalin.apibuilder.ApiBuilder.get;
import static io.javalin.apibuilder.ApiBuilder.path;
import static io.javalin.apibuilder.ApiBuilder.post;
import static io.javalin.apibuilder.ApiBuilder.put;

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
            // Аутентификация для всех маршрутов
            before(JwtHandler.authenticate());
            
            // Логирование начала запроса
            before(ctx -> {
                logger.info("▶️  Request started: {} {} (traceId: {})",
                        ctx.method(),
                        ctx.path(),
                        SimpleTracer.getCurrentTraceId());
            });

            // Базовые CRUD операции
            post(answerController::createAnswer);                   // POST /answers
            
            // Операции с конкретным ответом
            path("/{id}", () -> {
                get(answerController::getAnswerById);               // GET /answers/{id}
                put(answerController::updateAnswer);                // PUT /answers/{id}
                delete(answerController::deleteAnswer);             // DELETE /answers/{id}
                get("/correct", answerController::checkIfAnswerIsCorrect); // GET /answers/{id}/correct
            });

            // Операции с ответами по вопросам
            path("/by-question", () -> {
                get(answerController::getAnswersByQuestionId);           // GET /answers/by-question?questionId={id}
                get("/correct", answerController::getCorrectAnswersByQuestionId); // GET /answers/by-question/correct?questionId={id}
                delete(answerController::deleteAnswersByQuestionId);      // DELETE /answers/by-question?questionId={id}
                get("/count", answerController::countAnswersByQuestionId); // GET /answers/by-question/count?questionId={id}
                get("/count-correct", answerController::countCorrectAnswersByQuestionId); // GET /answers/by-question/count-correct?questionId={id}
            });

            // Логирование завершения запроса
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