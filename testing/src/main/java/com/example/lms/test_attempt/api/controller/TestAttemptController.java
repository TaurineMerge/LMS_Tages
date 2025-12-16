package com.example.lms.test_attempt.api.controller;

import java.sql.Date;
import java.time.LocalDate;
import java.util.List;
import java.util.UUID;
import java.util.stream.Collectors;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.test_attempt.api.dto.TestAttempt;
import com.example.lms.test_attempt.domain.model.TestAttemptModel;
import com.example.lms.test_attempt.domain.service.TestAttemptService;

import io.javalin.http.Context;
import io.javalin.http.HttpStatus;

/**
 * Контроллер для управления попытками прохождения тестов через REST API.
 * <p>
 * Работает с сервисным слоем {@link TestAttemptService}.
 */
public class TestAttemptController {

	private static final Logger logger = LoggerFactory.getLogger(TestAttemptController.class);

	private final TestAttemptService testAttemptService;

	/**
	 * Создает контроллер попыток тестов.
	 *
	 * @param testAttemptService сервис для работы с попытками
	 */
	public TestAttemptController(TestAttemptService testAttemptService) {
		this.testAttemptService = testAttemptService;
	}

	/**
	 * Получить все попытки тестов.
	 * GET /test-attempts
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void getAllTestAttempts(Context ctx) {
		try {
			logger.info("Получение всех попыток тестов");
			List<TestAttemptModel> attempts = testAttemptService.getAllTestAttempts();
			List<TestAttempt> dtos = attempts.stream()
					.map(this::toDTO)
					.collect(Collectors.toList());

			ctx.json(dtos);
			ctx.status(HttpStatus.OK);
			logger.info("Возвращено {} попыток", dtos.size());

		} catch (Exception e) {
			logger.error("Ошибка при получении всех попыток", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при получении попыток"));
		}
	}

	/**
	 * Получить попытку теста по ID.
	 * GET /test-attempts/{id}
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void getTestAttemptById(Context ctx) {
		try {
			String idParam = ctx.pathParam("id");
			logger.info("Получение попытки теста по ID: {}", idParam);

			UUID id = UUID.fromString(idParam);
			TestAttemptModel attempt = testAttemptService.findById(id)
					.orElse(null);

			if (attempt == null) {
				ctx.status(HttpStatus.NOT_FOUND)
						.json(createErrorResponse("Попытка теста с ID " + id + " не найдена"));
				logger.warn("Попытка теста с ID {} не найдена", id);
				return;
			}

			TestAttempt dto = toDTO(attempt);
			ctx.json(dto);
			ctx.status(HttpStatus.OK);
			logger.info("Возвращена попытка теста с ID: {}", id);

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат ID", e);
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат ID"));
		} catch (Exception e) {
			logger.error("Ошибка при получении попытки по ID", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при получении попытки"));
		}
	}

	/**
	 * Получить попытки теста по ID студента.
	 * GET /test-attempts/student/{studentId}
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void getTestAttemptsByStudentId(Context ctx) {
		try {
			String studentIdParam = ctx.pathParam("studentId");
			logger.info("Получение попыток тестов для студента с ID: {}", studentIdParam);

			UUID studentId = UUID.fromString(studentIdParam);
			List<TestAttemptModel> attempts = testAttemptService.getTestAttemptsByStudentId(studentId);
			List<TestAttempt> dtos = attempts.stream()
					.map(this::toDTO)
					.collect(Collectors.toList());

			ctx.json(dtos);
			ctx.status(HttpStatus.OK);
			logger.info("Возвращено {} попыток для студента {}", dtos.size(), studentId);

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат ID студента", e);
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат ID студента"));
		} catch (Exception e) {
			logger.error("Ошибка при получении попыток по студенту", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при получении попыток студента"));
		}
	}

	/**
	 * Получить попытки теста по ID теста.
	 * GET /test-attempts/test/{testId}
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void getTestAttemptsByTestId(Context ctx) {
		try {
			String testIdParam = ctx.pathParam("testId");
			logger.info("Получение попыток тестов для теста с ID: {}", testIdParam);

			UUID testId = UUID.fromString(testIdParam);
			List<TestAttemptModel> attempts = testAttemptService.getTestAttemptsByTestId(testId);
			List<TestAttempt> dtos = attempts.stream()
					.map(this::toDTO)
					.collect(Collectors.toList());

			ctx.json(dtos);
			ctx.status(HttpStatus.OK);
			logger.info("Возвращено {} попыток для теста {}", dtos.size(), testId);

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат ID теста", e);
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат ID теста"));
		} catch (Exception e) {
			logger.error("Ошибка при получении попыток по тесту", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при получении попыток теста"));
		}
	}

	/**
	 * Получить попытки теста по дате.
	 * GET /test-attempts/date/{date}
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void getTestAttemptsByDate(Context ctx) {
		try {
			String dateParam = ctx.pathParam("date");
			logger.info("Получение попыток тестов за дату: {}", dateParam);

			LocalDate date = LocalDate.parse(dateParam);
			List<TestAttemptModel> attempts = testAttemptService.getTestAttemptsByDate(date);
			List<TestAttempt> dtos = attempts.stream()
					.map(this::toDTO)
					.collect(Collectors.toList());

			ctx.json(dtos);
			ctx.status(HttpStatus.OK);
			logger.info("Возвращено {} попыток за дату {}", dtos.size(), date);

		} catch (Exception e) {
			logger.error("Ошибка при получении попыток по дате", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при получении попыток по дате"));
		}
	}

	/**
	 * Создать новую попытку теста.
	 * POST /test-attempts
	 * 
	 * Примечание: Для создания попытки нужны дополнительные поля (studentId,
	 * testId, attemptVersion),
	 * которых нет в DTO. Можно получить их из query параметров или JWT токена.
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void createTestAttempt(Context ctx) {
		try {
			TestAttempt dto = ctx.bodyAsClass(TestAttempt.class);
			logger.info("Создание новой попытки теста: {}", dto);

			// Получаем обязательные поля из query параметров
			String studentIdParam = ctx.queryParam("studentId");
			String testIdParam = ctx.queryParam("testId");
			String attemptVersion = ctx.queryParam("attemptVersion");

			if (studentIdParam == null || testIdParam == null || attemptVersion == null) {
				ctx.status(HttpStatus.BAD_REQUEST)
						.json(createErrorResponse("Необходимые параметры: studentId, testId, attemptVersion"));
				return;
			}

			if (dto.getDate_of_attempt() == null) {
				ctx.status(HttpStatus.BAD_REQUEST)
						.json(createErrorResponse("Дата попытки обязательна"));
				return;
			}

			UUID studentId = UUID.fromString(studentIdParam);
			UUID testId = UUID.fromString(testIdParam);

			// Создаем доменную модель
			TestAttemptModel attempt = new TestAttemptModel(
					null, // ID будет сгенерирован
					studentId,
					testId,
					dto.getDate_of_attempt().toLocalDate(),
					dto.getPoint() // может быть null
			);

			TestAttemptModel savedAttempt = testAttemptService.updateTestAttempt(attempt);
			TestAttempt responseDto = toDTO(savedAttempt);

			ctx.json(responseDto);
			ctx.status(HttpStatus.CREATED);
			logger.info("Создана новая попытка теста с ID: {}", responseDto.getId());

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат UUID", e);
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат ID: " + e.getMessage()));
		} catch (Exception e) {
			logger.error("Ошибка при создании попытки", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при создании попытки: " + e.getMessage()));
		}
	}

	/**
	 * Обновить существующую попытку теста.
	 * PUT /test-attempts/{id}
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void updateTestAttempt(Context ctx) {
		try {
			String idParam = ctx.pathParam("id");
			TestAttempt dto = ctx.bodyAsClass(TestAttempt.class);
			logger.info("Обновление попытки теста с ID: {}, данные: {}", idParam, dto);

			UUID id = UUID.fromString(idParam);

			// Получаем существующую попытку
			TestAttemptModel existingAttempt = testAttemptService.findById(id)
					.orElse(null);

			if (existingAttempt == null) {
				ctx.status(HttpStatus.NOT_FOUND)
						.json(createErrorResponse("Попытка теста с ID " + id + " не найдена"));
				logger.warn("Попытка теста с ID {} не найдена для обновления", id);
				return;
			}

			// Создаем обновленную модель
			TestAttemptModel updatedAttempt = new TestAttemptModel(
					existingAttempt.getId(),
					existingAttempt.getStudentId(),
					existingAttempt.getTestId(),
					dto.getDate_of_attempt() != null
							? dto.getDate_of_attempt().toLocalDate()
							: existingAttempt.getDateOfAttempt(),
					dto.getPoint() != null ? dto.getPoint() : existingAttempt.getPoint());

			TestAttemptModel savedAttempt = testAttemptService.updateTestAttempt(updatedAttempt);
			TestAttempt responseDto = toDTO(savedAttempt);

			ctx.json(responseDto);
			ctx.status(HttpStatus.OK);
			logger.info("Попытка теста с ID {} обновлена", id);

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат ID", e);
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат ID"));
		} catch (Exception e) {
			logger.error("Ошибка при обновлении попытки", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при обновлении попытки"));
		}
	}

	/**
	 * Удалить попытку теста по ID.
	 * DELETE /test-attempts/{id}
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void deleteTestAttempt(Context ctx) {
		try {
			String idParam = ctx.pathParam("id");
			logger.info("Удаление попытки теста с ID: {}", idParam);

			UUID id = UUID.fromString(idParam);
			testAttemptService.deleteTestAttempt(id);

			ctx.status(HttpStatus.NO_CONTENT);
			logger.info("Попытка теста с ID {} успешно удалена", id);

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат ID", e);
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат ID"));
		} catch (Exception e) {
			logger.error("Ошибка при удалении попытки", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при удалении попытки"));
		}
	}

	/**
	 * Получить завершенные попытки тестов.
	 * GET /test-attempts/completed
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void getCompletedTestAttempts(Context ctx) {
		try {
			logger.info("Получение завершенных попыток тестов");
			List<TestAttemptModel> attempts = testAttemptService.getCompletedTestAttempts();
			List<TestAttempt> dtos = attempts.stream()
					.map(this::toDTO)
					.collect(Collectors.toList());

			ctx.json(dtos);
			ctx.status(HttpStatus.OK);
			logger.info("Возвращено {} завершенных попыток", dtos.size());

		} catch (Exception e) {
			logger.error("Ошибка при получении завершенных попыток", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при получении завершенных попыток"));
		}
	}

	/**
	 * Получить незавершенные попытки тестов.
	 * GET /test-attempts/incomplete
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void getIncompleteTestAttempts(Context ctx) {
		try {
			logger.info("Получение незавершенных попыток тестов");
			List<TestAttemptModel> attempts = testAttemptService.getIncompleteTestAttempts();
			List<TestAttempt> dtos = attempts.stream()
					.map(this::toDTO)
					.collect(Collectors.toList());

			ctx.json(dtos);
			ctx.status(HttpStatus.OK);
			logger.info("Возвращено {} незавершенных попыток", dtos.size());

		} catch (Exception e) {
			logger.error("Ошибка при получении незавершенных попыток", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при получении незавершенных попыток"));
		}
	}

	/**
	 * Получить лучшую попытку студента по тесту.
	 * GET /test-attempts/best/student/{studentId}/test/{testId}
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	// public void getBestAttemptByStudentAndTest(Context ctx) {
	// try {
	// String studentIdParam = ctx.pathParam("studentId");
	// String testIdParam = ctx.pathParam("testId");
	// logger.info("Получение лучшей попытки студента {} по тесту {}",
	// studentIdParam, testIdParam);

	// UUID studentId = UUID.fromString(studentIdParam);
	// UUID testId = UUID.fromString(testIdParam);

	// TestAttemptModel bestAttempt = testAttemptService
	// .findBestAttemptByStudentAndTest(studentId, testId)
	// .orElse(null);

	// if (bestAttempt == null) {
	// ctx.status(HttpStatus.NOT_FOUND)
	// .json(createErrorResponse("Не найдено завершенных попыток студента " +
	// studentId + " по тесту " + testId));
	// logger.warn("Лучшая попытка студента {} по тесту {} не найдена", studentId,
	// testId);
	// return;
	// }

	// TestAttempt dto = toDTO(bestAttempt);
	// ctx.json(dto);
	// ctx.status(HttpStatus.OK);
	// logger.info("Возвращена лучшая попытка студента {} по тесту {}", studentId,
	// testId);

	// } catch (IllegalArgumentException e) {
	// logger.error("Неверный формат ID", e);
	// ctx.status(HttpStatus.BAD_REQUEST)
	// .json(createErrorResponse("Неверный формат ID"));
	// } catch (Exception e) {
	// logger.error("Ошибка при получении лучшей попытки", e);
	// ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
	// .json(createErrorResponse("Ошибка сервера при получении лучшей попытки"));
	// }
	// }

	/**
	 * Получить количество попыток студента по тесту.
	 * GET /test-attempts/count/student/{studentId}/test/{testId}
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void countAttemptsByStudentAndTest(Context ctx) {
		try {
			String studentIdParam = ctx.pathParam("studentId");
			String testIdParam = ctx.pathParam("testId");
			logger.info("Подсчет попыток студента {} по тесту {}", studentIdParam, testIdParam);

			UUID studentId = UUID.fromString(studentIdParam);
			UUID testId = UUID.fromString(testIdParam);

			int count = testAttemptService.countAttemptsByStudentAndTest(studentId, testId);

			var response = new CountResponse(count);
			ctx.json(response);
			ctx.status(HttpStatus.OK);
			logger.info("Студент {} имеет {} попыток по тесту {}", studentId, count, testId);

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат ID", e);
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат ID"));
		} catch (Exception e) {
			logger.error("Ошибка при подсчете попыток", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при подсчете попыток"));
		}
	}

	/**
	 * Получить попытки конкретного студента по конкретному тесту.
	 * GET /test-attempts/student/{studentId}/test/{testId}
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void getAttemptsByStudentAndTest(Context ctx) {
		try {
			String studentIdParam = ctx.pathParam("studentId");
			String testIdParam = ctx.pathParam("testId");
			logger.info("Получение попыток студента {} по тесту {}", studentIdParam, testIdParam);

			UUID studentId = UUID.fromString(studentIdParam);
			UUID testId = UUID.fromString(testIdParam);

			List<TestAttemptModel> attempts = testAttemptService.getAttemptsByStudentAndTest(studentId, testId);
			List<TestAttempt> dtos = attempts.stream()
					.map(this::toDTO)
					.collect(Collectors.toList());

			ctx.json(dtos);
			ctx.status(HttpStatus.OK);
			logger.info("Возвращено {} попыток студента {} по тесту {}", dtos.size(), studentId, testId);

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат ID", e);
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат ID"));
		} catch (Exception e) {
			logger.error("Ошибка при получении попыток студента по тесту", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при получении попыток студента по тесту"));
		}
	}

	/**
	 * Проверить существование попытки по ID.
	 * GET /test-attempts/exists/{id}
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void existsById(Context ctx) {
		try {
			String idParam = ctx.pathParam("id");
			logger.debug("Проверка существования попытки с ID: {}", idParam);

			UUID id = UUID.fromString(idParam);
			boolean exists = testAttemptService.existsById(id);

			var response = new ExistsResponse(exists);
			ctx.json(response);
			ctx.status(HttpStatus.OK);
			logger.debug("Попытка с ID {} существует: {}", id, exists);

		} catch (IllegalArgumentException e) {
			logger.error("Неверный формат ID", e);
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse("Неверный формат ID"));
		} catch (Exception e) {
			logger.error("Ошибка при проверке существования попытки", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при проверке существования попытки"));
		}
	}

	/**
	 * Завершить попытку теста (установить баллы).
	 * POST /test-attempts/{id}/complete
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void completeTestAttempt(Context ctx) {
		try {
			String idParam = ctx.pathParam("id");
			logger.info("Завершение попытки теста с ID: {}", idParam);

			UUID id = UUID.fromString(idParam);

			// Получаем баллы из тела запроса
			CompleteRequest request = ctx.bodyAsClass(CompleteRequest.class);
			if (request.points < 0 || request.points > 100) {
				ctx.status(HttpStatus.BAD_REQUEST)
						.json(createErrorResponse("Баллы должны быть в диапазоне 0-100"));
				return;
			}

			testAttemptService.completeTestAttempt(id, request.points);

			// Возвращаем обновленную попытку
			TestAttemptModel updatedAttempt = testAttemptService.getTestAttemptById(id);
			TestAttempt dto = toDTO(updatedAttempt);

			ctx.json(dto);
			ctx.status(HttpStatus.OK);
			logger.info("Попытка теста с ID {} завершена с баллами: {}", id, request.points);

		} catch (IllegalArgumentException e) {
			logger.error("Неверные параметры", e);
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse(e.getMessage()));
		} catch (Exception e) {
			logger.error("Ошибка при завершении попытки", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при завершении попытки"));
		}
	}

	/**
	 * Прикрепить сертификат к попытке.
	 * POST /test-attempts/{id}/attach-certificate
	 *
	 * @param ctx HTTP контекст Javalin
	 */
	public void attachCertificate(Context ctx) {
		try {
			String idParam = ctx.pathParam("id");
			logger.info("Прикрепление сертификата к попытке с ID: {}", idParam);

			UUID id = UUID.fromString(idParam);

			// Получаем ID сертификата из тела запроса
			CertificateRequest request = ctx.bodyAsClass(CertificateRequest.class);
			if (request.certificateId == null) {
				ctx.status(HttpStatus.BAD_REQUEST)
						.json(createErrorResponse("ID сертификата обязателен"));
				return;
			}

			testAttemptService.attachCertificate(id, request.certificateId);

			// Возвращаем обновленную попытку
			TestAttemptModel updatedAttempt = testAttemptService.getTestAttemptById(id);
			TestAttempt dto = toDTO(updatedAttempt);

			ctx.json(dto);
			ctx.status(HttpStatus.OK);
			logger.info("К попытке с ID {} прикреплен сертификат: {}", id, request.certificateId);

		} catch (IllegalArgumentException e) {
			logger.error("Неверные параметры", e);
			ctx.status(HttpStatus.BAD_REQUEST)
					.json(createErrorResponse(e.getMessage()));
		} catch (Exception e) {
			logger.error("Ошибка при прикреплении сертификата", e);
			ctx.status(HttpStatus.INTERNAL_SERVER_ERROR)
					.json(createErrorResponse("Ошибка сервера при прикреплении сертификата"));
		}
	}

	// --------------------------------------------------------------------
	// Вспомогательные методы
	// --------------------------------------------------------------------

	/**
	 * Преобразует доменную модель в DTO.
	 */
	private TestAttempt toDTO(TestAttemptModel attempt) {
		TestAttempt dto = new TestAttempt();

		dto.setId(attempt.getId());
		dto.setDate_of_attempt(Date.valueOf(attempt.getDateOfAttempt()));
		dto.setPoint(attempt.getPoint());

		// Вычисляем результат на основе баллов
		if (attempt.getPoint() == null) {
			dto.setResult("incomplete");
		} else if (attempt.getPoint() >= 70) {
			dto.setResult("passed");
		} else {
			dto.setResult("failed");
		}

		return dto;
	}

	/**
	 * Создает объект ошибки для ответа.
	 */
	private ErrorResponse createErrorResponse(String message) {
		return new ErrorResponse(message);
	}

	// Вспомогательные классы для запросов

	private static class CompleteRequest {
		public int points;
	}

	private static class CertificateRequest {
		public UUID certificateId;
	}

	private static class ErrorResponse {
		private final String error;

		public ErrorResponse(String error) {
			this.error = error;
		}

		public String getError() {
			return error;
		}
	}

	private static class CountResponse {
		private final int count;

		public CountResponse(int count) {
			this.count = count;
		}

		public int getCount() {
			return count;
		}
	}

	private static class ExistsResponse {
		private final boolean exists;

		public ExistsResponse(boolean exists) {
			this.exists = exists;
		}

		public boolean isExists() {
			return exists;
		}
	}
}