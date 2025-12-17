package com.example.lms.question.api.router;

import static io.javalin.apibuilder.ApiBuilder.*;

import com.example.lms.question.api.controller.QuestionController;
import com.example.lms.shared.router.RouterUtils;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import static com.example.lms.shared.router.RouterUtils.*;

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
 * Все запросы защищены JWT-аутентификацией и проверкой realm через {@link RouterUtils}.
 */
public class QuestionRouter {

    private static final Logger logger = LoggerFactory.getLogger(QuestionRouter.class);
    /**
     * Регистрирует маршруты для работы с вопросами.
     *
     * @param controller контроллер вопросов
     */
    public static void register(QuestionController controller) {
        validateController(controller, "QuestionController");

        path("/questions", () -> {

            // Стандартные middleware: аутентификация, traceId и логирование старта
            applyStandardBeforeMiddleware(logger);

            // /questions
            get(withRealm(READ_ACCESS_REALMS, controller::getQuestions));

            // создание вопроса: teacher + schema validation
            post(withValidationAndRealm(
                    "/schemas/question-schema.json",
                    TEACHER_REALM,
                    controller::createQuestion));

            // /questions/{id}
            path("/{id}", () -> {
                get(withRealm(READ_ACCESS_REALMS, controller::getQuestionById));

                put(withValidationAndRealm(
                        "/schemas/question-schema.json",
                        TEACHER_REALM,
                        controller::updateQuestion));

                delete(withRealm(TEACHER_REALM, controller::deleteQuestion));
            });

            // Логирование завершения запроса
            applyStandardAfterMiddleware(logger);
        });
    }
}