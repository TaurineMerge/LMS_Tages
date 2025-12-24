package com.example.lms.test.api.controller;

import java.io.StringWriter;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.test.api.dto.Test;
import com.example.lms.test.domain.service.TestService;
import com.example.lms.tracing.SimpleTracer;
import com.github.jknack.handlebars.Handlebars;
import com.github.jknack.handlebars.Template;

import io.javalin.http.Context;
import io.minio.messages.ErrorResponse;

/**
 * Контроллер для управления тестами.
 * <p>
 * Обрабатывает HTTP-запросы к ресурсам:
 * <ul>
 * <li>GET /tests — получение списка тестов в HTML-шаблоне</li>
 * <li>POST /tests — создание нового теста</li>
 * <li>GET /tests/{id} — получение теста по ID</li>
 * <li>PUT /tests/{id} — обновление теста</li>
 * <li>DELETE /tests/{id} — удаление теста</li>
 * </ul>
 *
 * Контроллер использует:
 * <ul>
 * <li>{@link TestService} — бизнес-логику работы с тестами</li>
 * <li>{@link Handlebars} — генерацию HTML-шаблонов</li>
 * <li>{@link SimpleTracer} — OpenTelemetry-трейсинг</li>
 * </ul>
 */
public class TestController {
	private static final Logger logger = LoggerFactory.getLogger(TestController.class);

	/** Сервисный слой, содержащий бизнес-логику работы с тестами. */
	private final TestService testService;

	/** Шаблонизатор Handlebars для генерации HTML. */
	private Handlebars handlebars;

	/**
	 * Создаёт контроллер тестов.
	 *
	 * @param testService сервис управления тестами
	 */
	public TestController(TestService testService, Handlebars handlebars) {
		this.testService = testService;
		this.handlebars = handlebars;
	}

	// Метод для рендеринга шаблонов
	private void renderTemplate(Context ctx, String templatePath, Map<String, Object> model) {
		try {
			Template template = handlebars.compile(templatePath);
			StringWriter writer = new StringWriter();
			template.apply(model, writer);
			ctx.contentType("text/html; charset=utf-8"); // добавлено для правильной кодировки
			ctx.result(writer.toString());
		} catch (Exception e) {
			e.printStackTrace();
			ctx.status(500).result("Ошибка рендеринга шаблона");
		}
	}

	/**
	 * Обработчик GET /tests.
	 * 
	 * Поддерживает query parameter: courseId для фильтрации тестов по курсу.
	 *
	 * @param ctx контекст HTTP-запроса
	 */
	public void getTests(Context ctx) {
		SimpleTracer.runWithSpan("getTests", () -> {
			String courseId = ctx.queryParam("courseId");
			List<Test> tests;
			String status = "success"; // Переименовал success в status для ясности
			
			if (courseId != null && !courseId.trim().isEmpty()) {
				// Фильтрация по courseId
				tests = testService.getTestsByCourseId(courseId);
			} else {
				// Все тесты
				tests = testService.getAllTests();
			}
			
			// Проверяем, есть ли тесты
			if (tests == null || tests.isEmpty()) {
				status = "error";
			}
			
			Map<String, Object> response = new HashMap<>();
			response.put("data", tests != null ? tests : new ArrayList<>());
			response.put("courseId", courseId);
			response.put("status", status);

			ctx.json(response);
		});
	}

	/**
	 * Обработчик POST /tests.
	 * <p>
	 * Получает JSON с данными теста, создаёт новый тест и возвращает его.
	 *
	 * @param ctx HTTP-контекст
	 */
	public void createTest(Context ctx) {
		SimpleTracer.runWithSpan("createTest", () -> {
			try {
				// 1. Получаем валидированный JSON (валидация структуры уже прошла)
				// Можно использовать JsonSchemaValidator.getValidatedJson(ctx) если нужно,
				// но проще использовать bodyValidator для бизнес-правил
				
				// 2. Бизнес-валидация через bodyValidator
				Test test = ctx.bodyValidator(Test.class)
						.check(t -> t.getTitle() != null && !t.getTitle().trim().isEmpty(), 
							"Название теста обязательно")
						.check(t -> t.getTitle().length() <= 200, 
							"Название теста не должно превышать 200 символов")
						.check(t -> t.getMin_point() != null && t.getMin_point() >= 0, 
							"Минимальный балл должен быть >= 0")
						.check(t -> t.getDescription() == null || t.getDescription().length() <= 1000, 
							"Описание не должно превышать 1000 символов")
						.get();
				
				// 3. Дополнительные бизнес-правила (если нужно)
				// Например, проверка уникальности названия в рамках курса
				// if (testService.existsByTitleAndCourseId(test.getTitle(), test.getCourseId())) {
				//     ctx.status(409).json(new ErrorResponse("Тест с таким названием уже существует в этом курсе"));
				//     return;
				// }
				
				// 4. Сохраняем тест
				Test createdTest = testService.createTest(test);
				logger.info("Создан новый тест с ID: {}", createdTest.getId());
				ctx.status(201).json(createdTest);
				
			} catch (IllegalArgumentException e) {
				// Ошибка валидации через bodyValidator
				logger.error("Ошибка валидации при создании теста", e);
				ctx.status(400).json(new ErrorResponse("Ошибка валидации: " + e.getMessage()));
			} catch (Exception e) {
				logger.error("Ошибка при создании теста", e);
				ctx.status(500).json(new ErrorResponse("Ошибка при создании теста"));
			}
		});
	}

	/**
	 * Обработчик GET /tests/{id}.
	 * <p>
	 * Возвращает тест по его идентификатору.
	 *
	 * @param ctx HTTP-контекст
	 */
	public void getTestById(Context ctx) {
		SimpleTracer.runWithSpan("getTestById", () -> {
			try {
				String idParam = ctx.pathParam("id");
				UUID id = UUID.fromString(idParam);
				
				Test test = testService.getTestById(id.toString());
				if (test == null) {
					ctx.status(404).json(new ErrorResponse("Тест с ID " + id + " не найден"));
					return;
				}
				ctx.json(test);
			} catch (IllegalArgumentException e) {
				logger.error("Неверный формат UUID", e);
				ctx.status(400).json(new ErrorResponse("Неверный формат идентификатора"));
			} catch (Exception e) {
				logger.error("Ошибка при получении теста по ID", e);
				ctx.status(500).json(new ErrorResponse("Ошибка при получении теста"));
			}
		});
	}

	/**
	 * Обработчик PUT /tests/{id}.
	 * <p>
	 * Обновляет тест: данные принимаются в JSON, а ID берётся из path parameter.
	 *
	 * @param ctx HTTP-контекст
	 */
	public void updateTest(Context ctx) {
		SimpleTracer.runWithSpan("updateTest", () -> {
			try {
				// 1. Валидация UUID параметра
				String idParam = ctx.pathParam("id");
				UUID id;
				try {
					id = UUID.fromString(idParam);
				} catch (IllegalArgumentException e) {
					ctx.status(400).json(new ErrorResponse("Неверный формат UUID"));
					return;
				}
				
				// 2. Проверка существования теста
				Test existingTest = testService.getTestById(id.toString());
				if (existingTest == null) {
					ctx.status(404).json(new ErrorResponse("Тест с ID " + id + " не найден"));
					return;
				}
				
				// 3. Бизнес-валидация через bodyValidator
				Test test = ctx.bodyValidator(Test.class)
						.check(t -> t.getTitle() != null && !t.getTitle().trim().isEmpty(), 
							"Название теста обязательно")
						.check(t -> t.getTitle().length() <= 200, 
							"Название теста не должно превышать 200 символов")
						.check(t -> t.getMin_point() != null && t.getMin_point() >= 0, 
							"Минимальный балл должен быть >= 0")
						.check(t -> t.getDescription() == null || t.getDescription().length() <= 1000, 
							"Описание не должно превышать 1000 символов")
						.get();
				
				test.setId(id.toString());
				
				// 4. Обновляем тест
				Test updatedTest = testService.updateTest(test);
				logger.info("Обновлен тест с ID: {}", id);
				ctx.json(updatedTest);
            
			} catch (IllegalArgumentException e) {
				logger.error("Ошибка валидации при обновлении теста", e);
				ctx.status(400).json(new ErrorResponse("Ошибка валидации: " + e.getMessage()));
			} catch (Exception e) {
				logger.error("Ошибка при обновлении теста", e);
				ctx.status(500).json(new ErrorResponse("Ошибка при обновлении теста"));
			}
		});
	}

	/**
	 * Обработчик DELETE /tests/{id}.
	 * <p>
	 * Удаляет тест по идентификатору и возвращает JSON вида:
	 * 
	 * <pre>
	 * {"deleted": true}
	 * </pre>
	 *
	 * @param ctx HTTP-контекст
	 */
	public void deleteTest(Context ctx) {
		SimpleTracer.runWithSpan("deleteTest", () -> {
			String id = ctx.pathParam("id");
			boolean deleted = testService.deleteTest(id);
			ctx.json(Map.of("deleted", deleted));
		});
	}

	/**
	 * Обработчик GET /tests/by-course/validate/{course-id}.
	 * <p>
	 * Проверяет, существует ли тест для данного курса по courseId и возвращает JSON вида:
	 * 
	 * <pre>
	 * {"exists": true}
	 * </pre>
	 * 
	 * @param ctx HTTP-контекст
	 */
	public void existsByCourseId(Context ctx) {
		SimpleTracer.runWithSpan("testExistsByCourseId", () -> {
			String course_id = ctx.pathParam("courseId");
			boolean exist = testService.existsByCourseId(course_id);
			ctx.json(Map.of("exists", exist));
		});
	}

	/**
	 * Обработчик GET /tests/by-course/{courseId}.
	 * <p>
	 * Возвращает тест по идентификатору курса.
	 *
	 * @param ctx HTTP-контекст
	 */
	public void getTestByCourseId(Context ctx) {
		SimpleTracer.runWithSpan("getTestByCourseId", () -> {
			String courseId = ctx.pathParam("courseId");
			List<Test> tests = testService.getTestsByCourseId(courseId);
			ctx.json(tests);
		});
	}

	/**
     * Класс для передачи сообщений об ошибках
     */
    private static class ErrorResponse {
        private final String error;

        public ErrorResponse(String error) {
            this.error = error;
        }

        public String getError() {
            return error;
        }
    }
}