package com.example.lms.test.api.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.test.api.controller.TestController;
import com.example.lms.shared.router.RouterUtils;

import static io.javalin.apibuilder.ApiBuilder.*;
import static com.example.lms.shared.router.RouterUtils.*;

import java.util.Set;

/**
 * Роутер для управления тестами.
 * <p>
 * Регистрирует REST-эндпоинты:
 * <ul>
 * <li>GET /tests — получить список тестов (HTML)</li>
 * <li>POST /tests — создать тест</li>
 * <li>GET /tests/{id} — получить тест по ID</li>
 * <li>PUT /tests/{id} — обновить тест</li>
 * <li>DELETE /tests/{id} — удалить тест</li>
 * </ul>
 * 
 * <h2>Политика доступа:</h2>
 * <ul>
 * <li><b>Студенты:</b> просмотр опубликованных тестов</li>
 * <li><b>Преподаватели:</b> полный доступ к своим тестам</li>
 * </ul>
 * 
 * @see TestController
 * @see RouterUtils
 */
public class TestRouter {

	private static final Logger logger = LoggerFactory.getLogger(TestRouter.class);

	/**
	 * Регистрирует маршруты группы /tests и их подмаршрутов.
	 * 
	 * @param testController контроллер, содержащий обработчики запросов
	 * @throws IllegalArgumentException если testController равен null
	 */
	public static void register(TestController testController) {
		validateController(testController, "TestController");

		path("/tests", () -> {
			// Стандартные middleware
			applyStandardBeforeMiddleware(logger);

			// Список и создание тестов
			get(withRealm(TEACHER_REALM, testController::getTests));

			// Создание теста — только преподаватель + schema validation
			post(withValidationAndRealm(
					"/schemas/test-schema.json",
					TEACHER_REALM,
					testController::createTest));

			path("/{id}", () -> {
				// Просмотр теста - для всех
				get(withRealm(READ_ACCESS_REALMS, testController::getTestById));

				// Редактирование — преподаватель + schema validation
				put(withValidationAndRealm(
						"/schemas/test-schema.json",
						TEACHER_REALM,
						testController::updateTest));

				// Удаление — только преподаватель
				delete(withRealm(Set.of(TEACHER_REALM), testController::deleteTest));
			});

			applyStandardAfterMiddleware(logger);
		});

		// Health check without auth - bypass authentication completely
		get("/tests/health", ctx -> {
			// Skip authentication for health check
			ctx.result("OK");
		});
	}
}