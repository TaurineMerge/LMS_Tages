package com.example.lms.internal.api.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.internal.api.controller.InternalApiController;

import static io.javalin.apibuilder.ApiBuilder.*;
import static com.example.lms.shared.router.RouterUtils.*;

/**
 * Роутер для Internal API эндпоинтов.
 * Определяет маршруты для взаимодействия с Python-сервисом.
 */
public class InternalApiRouter {
    private InternalApiController internalApiController;
    private static final Logger logger = LoggerFactory.getLogger(InternalApiRouter.class);

    public InternalApiRouter(InternalApiController internalApiController) {
        this.internalApiController = internalApiController;
    }

    /**
     * Регистрирует маршруты группы /tests и их подмаршрутов.
     * 
     * @param internalApiController контроллер, содержащий обработчики запросов
     * @throws IllegalArgumentException если testController равен null
     */
    public static void register(InternalApiController internalApiController) {
        validateController(internalApiController, "TestController");

        path("/internal", () -> {
            applyStandardBeforeMiddleware(logger);

            path("/users/{userId}", () -> {
                path("/attempts", () -> {
                    get(withRealm(READ_ACCESS_REALMS, internalApiController::getUserAttempts));
                });
                path("/stats", () -> {
                    get(withRealm(READ_ACCESS_REALMS, internalApiController::getUserStats));
                });
            });

            path("/attempts/{attemptId}", () -> {
                get(withRealm(READ_ACCESS_REALMS, internalApiController::getAttemptDetail));
            });

            applyStandardAfterMiddleware(logger);
        });

        get("/internal/health", ctx -> {
            ctx.result("OK");
        });
    }
}