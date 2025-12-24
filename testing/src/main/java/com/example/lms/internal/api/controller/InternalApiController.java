package com.example.lms.internal.api.controller;

import java.util.List;
import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.internal.api.dto.AttemptDetail;
import com.example.lms.internal.api.dto.AttemptsListItem;
import com.example.lms.internal.api.dto.CourseDraftResponse;
import com.example.lms.internal.api.dto.CourseTestResponse;
import com.example.lms.internal.api.dto.UserStats;
import com.example.lms.internal.service.InternalApiService;

import io.javalin.http.Context;
import io.javalin.http.HttpStatus;

/**
 * Контроллер для Internal API.
 * Предоставляет эндпоинты для взаимодействия с Python-сервисом personal
 * account.
 */
public class InternalApiController {

	private static final Logger logger = LoggerFactory.getLogger(InternalApiController.class);

	private final InternalApiService internalApiService;

	public InternalApiController(InternalApiService internalApiService) {
		this.internalApiService = internalApiService;
	}

	/**
	 * Получить детальную информацию о попытке.
	 * GET /internal/attempts/{attempt_id}
	 *
	 * @param ctx HTTP контекст
	 */
	public void getAttemptDetail(Context ctx) {
		try {
			String attemptIdParam = ctx.pathParam("attemptId");
			logger.info("Запрос деталей попытки: {}", attemptIdParam);

			UUID attemptId = UUID.fromString(attemptIdParam);
			AttemptDetail detail = internalApiService.getAttemptDetail(attemptId);

			if (detail == null) {
				ctx.status(HttpStatus.NOT_FOUND)
						.json(createErrorResponse("Попытка не найдена"));
				logger.warn("Попытка {} не найдена", attemptId);
				return;
			}

			ctx.json(detail);
			ctx.status(HttpStatus.OK);
			logger.info("Возвращены детали попытки: {}", attemptId);

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат UUID: {}", e.getMessage());
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат attempt_id"));
		} catch (Exception e) {
			logger.error("Ошибка при получении деталей попытки", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Внутренняя ошибка сервера"));
		}
	}

	/**
	 * Получить список всех попыток пользователя.
	 * GET /internal/users/{user_id}/attempts
	 *
	 * @param ctx HTTP контекст
	 */
	public void getUserAttempts(Context ctx) {
		try {
			String userIdParam = ctx.pathParam("userId");
			logger.info("Запрос попыток пользователя: {}", userIdParam);

			UUID userId = UUID.fromString(userIdParam);
			List<AttemptsListItem> attempts = internalApiService.getUserAttempts(userId);

			ctx.json(attempts);
			ctx.status(HttpStatus.OK);
			logger.info("Возвращено {} попыток для пользователя {}", attempts.size(), userId);

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат UUID: {}", e.getMessage());
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат user_id"));
		} catch (Exception e) {
			logger.error("Ошибка при получении попыток пользователя", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Внутренняя ошибка сервера"));
		}
	}

	/**
	 * Получить статистику пользователя по всем попыткам.
	 * GET /internal/users/{user_id}/stats
	 *
	 * @param ctx HTTP контекст
	 */
	public void getUserStats(Context ctx) {
		try {
			String userIdParam = ctx.pathParam("userId");
			logger.info("Запрос статистики пользователя: {}", userIdParam);

			UUID userId = UUID.fromString(userIdParam);
			UserStats stats = internalApiService.getUserStats(userId);

			ctx.json(stats);
			ctx.status(HttpStatus.OK);
			logger.info("Возвращена статистика для пользователя {}", userId);

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат UUID: {}", e.getMessage());
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат user_id"));
		} catch (Exception e) {
			logger.error("Ошибка при получении статистики пользователя", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Внутренняя ошибка сервера"));
		}
	}

	/**
	 * Получить тест для курса.
	 * GET /internal/course/{courseId}/test
	 *
	 * @param ctx HTTP контекст
	 */
	public void getCourseTest(Context ctx) {
		try {
			String courseIdParam = ctx.pathParam("courseId");
			logger.info("Запрос теста для курса: {}", courseIdParam);

			UUID courseId = UUID.fromString(courseIdParam);
			CourseTestResponse response = internalApiService.getTestByCourseId(courseId);

			ctx.json(response);
			ctx.status(HttpStatus.OK);
			logger.info("Возвращен тест для курса {}", courseId);

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат UUID: {}", e.getMessage());
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат courseId"));
		} catch (Exception e) {
			logger.error("Ошибка при получении теста для курса", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Внутренняя ошибка сервера"));
		}
	}

	/**
	 * Получить черновик для курса.
	 * GET /internal/course/{courseId}/draft
	 *
	 * @param ctx HTTP контекст
	 */
	public void getCourseDraft(Context ctx) {
		try {
			String courseIdParam = ctx.pathParam("courseId");
			logger.info("Запрос черновика для курса: {}", courseIdParam);

			UUID courseId = UUID.fromString(courseIdParam);
			CourseDraftResponse response = internalApiService.getDraftByCourseId(courseId);

			ctx.json(response);
			ctx.status(HttpStatus.OK);
			logger.info("Возвращен черновик для курса {}", courseId);

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат UUID: {}", e.getMessage());
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат courseId"));
		} catch (Exception e) {
			logger.error("Ошибка при получении черновика для курса", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Внутренняя ошибка сервера"));
		}
	}

	// ------------------------------------------------------------------
	// ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ
	// ------------------------------------------------------------------

	/**
	 * Создает объект ошибки для ответа.
	 */
	private ErrorResponse createErrorResponse(String message) {
		return new ErrorResponse(message);
	}

	/**
	 * Класс для ответа с ошибкой.
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
