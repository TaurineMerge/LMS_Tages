package com.example.lms.answer.api.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.Main;
import com.example.lms.answer.api.controller.AnswerController;
import com.example.lms.security.JwtHandler;
import com.example.lms.tracing.SimpleTracer;

import io.github.cdimascio.dotenv.Dotenv;

import static io.javalin.apibuilder.ApiBuilder.after;
import static io.javalin.apibuilder.ApiBuilder.before;
import static io.javalin.apibuilder.ApiBuilder.delete;
import static io.javalin.apibuilder.ApiBuilder.get;
import static io.javalin.apibuilder.ApiBuilder.path;
import static io.javalin.apibuilder.ApiBuilder.post;
import static io.javalin.apibuilder.ApiBuilder.put;

import java.util.Set;

/**
 * Роутер для маршрутов управления ответами.
 * Определяет конечные точки REST API для операций с ответами.
 */
public class AnswerRouter {
    private static final Logger logger = LoggerFactory.getLogger(Main.class);
    private static final Dotenv dotenv = Dotenv.load();

    private static final String KEYCLOAK_STUDENT_REALM = dotenv.get("KEYCLOAK_STUDENT_REALM");
    private static final String KEYCLOAK_TEACHER_REALM = dotenv.get("KEYCLOAK_TEACHER_REALM");

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
            post(ctx -> { // POST /answers
                JwtHandler.requireRealm(KEYCLOAK_TEACHER_REALM);
                answerController.createAnswer(ctx);
            });
            // Операции с ответами по вопросам
            path("/by-question", () -> {
                delete(ctx -> { // DELETE /answers/by-question?questionId={id}
                    JwtHandler.requireRealm(KEYCLOAK_TEACHER_REALM);
                    answerController.deleteAnswersByQuestionId(ctx);
                });
                get(ctx -> { // GET /answers/by-question?questionId={id}
                    JwtHandler.requireRealm(Set.of(KEYCLOAK_STUDENT_REALM, KEYCLOAK_TEACHER_REALM));
                    answerController.getAnswersByQuestionId(ctx);
                });
                path("/correct", () -> {
                    get(ctx -> { // GET /answers/by-question/correct?questionId={id}
                        JwtHandler.requireRealm(Set.of(KEYCLOAK_STUDENT_REALM, KEYCLOAK_TEACHER_REALM));
                        answerController.getCorrectAnswersByQuestionId(ctx);
                    });
                });
                path("/count", () -> { // GET /answers/by-question/count?questionId={id}
                    get(ctx -> {
                        JwtHandler.requireRealm(Set.of(KEYCLOAK_STUDENT_REALM, KEYCLOAK_TEACHER_REALM));
                        answerController.countAnswersByQuestionId(ctx);
                    });
                });
                path("/count-correct", () -> {
                    get(ctx -> { // GET /answers/by-question/count-correct?questionId={id}
                        JwtHandler.requireRealm(Set.of(KEYCLOAK_STUDENT_REALM, KEYCLOAK_TEACHER_REALM));
                        answerController.countCorrectAnswersByQuestionId(ctx);
                    });
                });
            });

            // Операции с конкретным ответом
            path("/{id}", () -> {
                get(ctx -> { // GET /answers/{id}
                    JwtHandler.requireRealm(Set.of(KEYCLOAK_STUDENT_REALM, KEYCLOAK_TEACHER_REALM));
                    answerController.getAnswerById(ctx);
                });
                put(ctx -> { // PUT /answers/{id}
                    JwtHandler.requireRealm(KEYCLOAK_TEACHER_REALM);
                    answerController.updateAnswer(ctx);
                });
                delete(ctx -> { // DELETE /answers/{id}
                    JwtHandler.requireRealm(KEYCLOAK_TEACHER_REALM);
                    answerController.deleteAnswer(ctx);
                });
                path("/correct", () -> {
                    get(ctx -> { // GET /answers/{id}/correct
                        JwtHandler.requireRealm(Set.of(KEYCLOAK_STUDENT_REALM, KEYCLOAK_TEACHER_REALM));
                        answerController.checkIfAnswerIsCorrect(ctx);
                    });
                });
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