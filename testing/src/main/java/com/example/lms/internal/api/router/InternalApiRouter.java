package com.example.lms.internal.api.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.internal.api.controller.InternalApiController;

import static io.javalin.apibuilder.ApiBuilder.*;
import static com.example.lms.shared.router.RouterUtils.*;

/**
 * Роутер для внутреннего API, предоставляющего данные другим сервисам системы LMS.
 * <p>
 * Определяет все маршруты, доступные для внутренних интеграций, включая:
 * <ul>
 *   <li>Получение статистики пользователей</li>
 *   <li>Доступ к попыткам прохождения тестов</li>
 *   <li>Получение информации о тестах и черновиках по курсам</li>
 *   <li>Health check эндпоинты для мониторинга</li>
 * </ul>
 * <p>
 * Роутер применяет стандартную middleware для аутентификации, авторизации и обработки ошибок.
 *
 * @author Команда разработки LMS
 * @version 1.0
 * @see InternalApiController
 * @see com.example.lms.shared.router.RouterUtils
 */
public class InternalApiRouter {
    
    private static final Logger logger = LoggerFactory.getLogger(InternalApiRouter.class);
    private InternalApiController internalApiController;

    /**
     * Создает новый экземпляр роутера для внутреннего API.
     *
     * @param internalApiController контроллер, обрабатывающий запросы внутреннего API
     * @throws NullPointerException если {@code internalApiController} равен {@code null}
     */
    public InternalApiRouter(InternalApiController internalApiController) {
        this.internalApiController = internalApiController;
    }

    /**
     * Регистрирует все маршруты внутреннего API в приложении Javalin.
     * <p>
     * Метод определяет следующие группы маршрутов:
     * <ol>
     *   <li>{@code /internal/health} - health check эндпоинт</li>
     *   <li>{@code /internal/categories/{categoryId}/courses/{courseId}/test} - получение теста курса</li>
     *   <li>{@code /internal/users/{userId}/attempts} - попытки пользователя</li>
     *   <li>{@code /internal/users/{userId}/stats} - статистика пользователя</li>
     *   <li>{@code /internal/attempts/{attemptId}} - детали попытки</li>
     *   <li>{@code /internal/categories/{categoryId}/courses/{courseId}/draft} - черновик теста курса</li>
     * </ol>
     * <p>
     * Маршруты, требующие авторизации, оборачиваются с помощью {@code withRealm} для проверки прав доступа.
     *
     * @param internalApiController контроллер, содержащий обработчики запросов внутреннего API
     * @throws IllegalArgumentException если {@code internalApiController} равен {@code null}
     * @throws IllegalStateException если не удалось применить middleware или зарегистрировать маршруты
     */
    public static void register(InternalApiController internalApiController) {
        validateController(internalApiController, "InternalApiController");

        path("/internal", () -> {
            // Публичный health check эндпоинт (доступен без авторизации)
            get("/health", ctx -> ctx.result("OK"));
            
            // Публичный эндпоинт для получения теста курса
            get("/categories/{categoryId}/courses/{courseId}/test", 
                internalApiController::getCourseTest);

            // Маршруты для работы с пользователями
            path("/users/{userId}", () -> {
                path("/attempts", () -> {
                    get(withRealm(READ_ACCESS_REALMS, internalApiController::getUserAttempts));
                });
                
                path("/stats", () -> {
                    get(withRealm(READ_ACCESS_REALMS, internalApiController::getUserStats));
                });
            });

            // Маршрут для получения деталей конкретной попытки
            path("/attempts/{attemptId}", () -> {
                get(withRealm(READ_ACCESS_REALMS, internalApiController::getAttemptDetail));
            });

            // Маршрут для получения черновика теста курса
            path("/categories/{categoryId}/courses/{courseId}/draft", () -> {
                get(withRealm(READ_ACCESS_REALMS, internalApiController::getCourseDraft));
            });

            // Применение стандартной after-middleware (логирование, обработка ошибок)
            applyStandardAfterMiddleware(logger);
        });
    }
}