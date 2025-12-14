package com.example.lms.answer.api.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.answer.api.controller.AnswerController;
import com.example.lms.shared.router.RouterUtils;

import static io.javalin.apibuilder.ApiBuilder.*;
import static com.example.lms.shared.router.RouterUtils.*;

/**
 * Роутер для управления ответами на вопросы тестов.
 * <p>
 * Определяет REST API эндпоинты для работы с ответами, включая:
 * <ul>
 * <li>Создание, чтение, обновление и удаление ответов</li>
 * <li>Получение ответов по ID вопроса</li>
 * <li>Получение правильных ответов</li>
 * <li>Подсчет ответов и правильных ответов</li>
 * </ul>
 * 
 * <h2>Структура маршрутов:</h2>
 * 
 * <pre>
 * POST   /answers                                    - создание ответа (teacher)
 * 
 * GET    /answers/by-question?questionId={id}        - получить все ответы вопроса (student, teacher)
 * DELETE /answers/by-question?questionId={id}        - удалить все ответы вопроса (teacher)
 * GET    /answers/by-question/correct?questionId={id} - получить правильные ответы (student, teacher)
 * GET    /answers/by-question/count?questionId={id}   - подсчет ответов (student, teacher)
 * GET    /answers/by-question/count-correct?questionId={id} - подсчет правильных ответов (student, teacher)
 * 
 * GET    /answers/{id}                               - получить ответ по ID (student, teacher)
 * PUT    /answers/{id}                               - обновить ответ (teacher)
 * DELETE /answers/{id}                               - удалить ответ (teacher)
 * GET    /answers/{id}/correct                       - проверить правильность ответа (student, teacher)
 * </pre>
 * 
 * <h2>Политика доступа:</h2>
 * <ul>
 * <li><b>Студенты:</b> могут только читать ответы и проверять их
 * правильность</li>
 * <li><b>Преподаватели:</b> полный доступ ко всем операциям</li>
 * </ul>
 * 
 * @see AnswerController
 * @see RouterUtils
 */
public class AnswerRouter {
    private static final Logger logger = LoggerFactory.getLogger(AnswerRouter.class);

    /**
     * Регистрирует все маршруты для управления ответами.
     * <p>
     * Использует стандартные middleware из {@link RouterUtils} для:
     * <ul>
     * <li>JWT аутентификации</li>
     * <li>Логирования запросов</li>
     * <li>Авторизации по realm</li>
     * </ul>
     *
     * @param answerController контроллер для обработки бизнес-логики запросов
     * @throws IllegalArgumentException если answerController равен null
     */
    public static void register(AnswerController answerController) {
        validateController(answerController, "AnswerController");

        path("/answers", () -> {

            // Стандартные middleware: аутентификация и логирование
            applyStandardBeforeMiddleware(logger);

            // Создание ответа
            post(withRealm(TEACHER_REALM, answerController::createAnswer));

            // Операции с ответами по ID вопроса
            path("/by-question", () -> {
                delete(withRealm(TEACHER_REALM, answerController::deleteAnswersByQuestionId));
                get(withRealm(READ_ACCESS_REALMS, answerController::getAnswersByQuestionId));

                path("/correct", () -> {
                    get(withRealm(READ_ACCESS_REALMS, answerController::getCorrectAnswersByQuestionId));
                });

                path("/count", () -> {
                    get(withRealm(READ_ACCESS_REALMS, answerController::countAnswersByQuestionId));
                });

                path("/count-correct", () -> {
                    get(withRealm(READ_ACCESS_REALMS, answerController::countCorrectAnswersByQuestionId));
                });
            });

            // Операции с конкретным ответом
            path("/{id}", () -> {
                get(withRealm(READ_ACCESS_REALMS, answerController::getAnswerById));
                put(withRealm(TEACHER_REALM, answerController::updateAnswer));
                delete(withRealm(TEACHER_REALM, answerController::deleteAnswer));

                path("/correct", () -> {
                    get(withRealm(READ_ACCESS_REALMS, answerController::checkIfAnswerIsCorrect));
                });
            });

            // Логирование завершения запроса
            applyStandardAfterMiddleware(logger);
        });
    }
}