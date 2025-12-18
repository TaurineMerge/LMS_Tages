package com.example.lms.content.api.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.content.api.controller.ContentController;
import com.example.lms.shared.router.RouterUtils;

import static io.javalin.apibuilder.ApiBuilder.*;
import static com.example.lms.shared.router.RouterUtils.*;

import java.util.Set;

/**
 * Роутер для управления элементами контента.
 * <p>
 * Регистрирует REST-эндпоинты:
 * <ul>
 * <li>GET /contents — получить список элементов контента (JSON)</li>
 * <li>GET /contents/html — получить список элементов контента (HTML)</li>
 * <li>POST /contents — создать элемент контента</li>
 * <li>GET /contents/{id} — получить элемент контента по ID</li>
 * <li>PUT /contents/{id} — обновить элемент контента</li>
 * <li>DELETE /contents/{id} — удалить элемент контента</li>
 * <li>GET /contents/search — поиск по содержимому</li>
 * <li>GET /contents/type — фильтрация по типу контента</li>
 * <li>GET /contents/question/{questionId} — поиск по вопросу</li>
 * <li>GET /contents/answer/{answerId} — поиск по ответу</li>
 * <li>GET /contents/validate/{id} — проверка существования</li>
 * <li>GET /contents/max-order — элементы с максимальным порядком</li>
 * </ul>
 * 
 * <h2>Политика доступа:</h2>
 * <ul>
 * <li><b>Студенты:</b> просмотр элементов контента</li>
 * <li><b>Преподаватели:</b> полный доступ к элементам контента</li>
 * </ul>
 * 
 * @see ContentController
 * @see RouterUtils
 */
public class ContentRouter {
    private static final Logger logger = LoggerFactory.getLogger(ContentRouter.class);

    /**
     * Регистрирует маршруты группы /contents и их подмаршрутов.
     * 
     * @param contentController контроллер, содержащий обработчики запросов
     * @throws IllegalArgumentException если contentController равен null
     */
    public static void register(ContentController contentController) {
        validateController(contentController, "ContentController");

        path("/contents", () -> {
            // Стандартные middleware
            applyStandardBeforeMiddleware(logger);

            // Получение списка элементов контента в JSON
            get(withRealm(READ_ACCESS_REALMS, contentController::getAllContents));
            
            // Получение списка элементов контента в HTML
            get("/html", withRealm(READ_ACCESS_REALMS, contentController::getContentsHtml));
            
            // Создание нового элемента контента
            post(withRealm(TEACHER_REALM, contentController::createContent));
            
            // Поиск по содержимому
            get("/search", withRealm(READ_ACCESS_REALMS, contentController::searchByContent));
            
            // Фильтрация по типу контента
            get("/type", withRealm(READ_ACCESS_REALMS, contentController::filterByType));
            
            // Получение элементов с максимальным порядком
            get("/max-order", withRealm(READ_ACCESS_REALMS, contentController::getWithMaxOrder));

            // Валидация существования
            get("/validate/{id}", withRealm(READ_ACCESS_REALMS, contentController::existsById));

            // Подмаршруты для конкретного элемента контента
            path("/{id}", () -> {
                // Просмотр элемента контента - для всех
                get(withRealm(READ_ACCESS_REALMS, contentController::getContentById));

                // Редактирование и удаление - только для преподавателей
                put(withRealm(TEACHER_REALM, contentController::updateContent));
                delete(withRealm(Set.of(TEACHER_REALM), contentController::deleteContent));
            });

            // Поиск элементов контента по вопросу
            get("/question/{questionId}", withRealm(READ_ACCESS_REALMS, contentController::getByQuestionId));
            
            // Поиск элементов контента по ответу
            get("/answer/{answerId}", withRealm(READ_ACCESS_REALMS, contentController::getByAnswerId));

            applyStandardAfterMiddleware(logger);
        });

        // Health check без аутентификации - полностью обходим аутентификацию
        get("/content/health", ctx -> {
            ctx.result("OK");
        });
    }
}