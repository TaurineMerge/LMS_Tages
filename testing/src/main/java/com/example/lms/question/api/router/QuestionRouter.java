package com.example.lms.question.api.router;

import static io.javalin.apibuilder.ApiBuilder.*;

import com.example.lms.Main;
import com.example.lms.question.api.controller.QuestionController;
import com.example.lms.security.JwtHandler;
import com.example.lms.tracing.SimpleTracer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Роутер Javalin для сущности {@code Question}.
 * <p>
 * Регистрирует HTTP-маршруты:
 * <ul>
 *     <li>GET    /questions          — список вопросов (опционально по testId)</li>
 *     <li>POST   /questions          — создание вопроса</li>
 *     <li>GET    /questions/{id}     — получение вопроса по ID</li>
 *     <li>PUT    /questions/{id}     — обновление вопроса</li>
 *     <li>DELETE /questions/{id}     — удаление вопроса</li>
 * </ul>
 *
 * Все запросы защищены JWT-аутентификацией через {@link JwtHandler}.
 */
public class QuestionRouter {

    private static final Logger logger = LoggerFactory.getLogger(Main.class);

    /**
     * Регистрирует маршруты для работы с вопросами.
     *
     * @param controller контроллер вопросов
     */
    public static void register(QuestionController controller) {
        path("/questions", () -> {
            // Глобальный before для аутентификации
            before(JwtHandler.authenticate());

            // Логирование начала запроса
            before(ctx -> {
                logger.info("Request started: {} {} (traceId: {})",
                        ctx.method(),
                        ctx.path(),
                        SimpleTracer.getCurrentTraceId());
            });

            // /questions
            get(controller::getQuestions);
            post(controller::createQuestion);

            // /questions/{id}
            path("/{id}", () -> {
                get(controller::getQuestionById);
                put(controller::updateQuestion);
                delete(controller::deleteQuestion);
            });

            // Логирование завершения запроса
            after(ctx -> {
                logger.info("Request completed: {} {} -> {} (traceId: {})",
                        ctx.method(),
                        ctx.path(),
                        ctx.status(),
                        SimpleTracer.getCurrentTraceId());
            });
        });
    }
}