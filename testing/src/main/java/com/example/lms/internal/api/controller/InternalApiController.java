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
 * Контроллер для внутреннего API, предоставляющий эндпоинты для взаимодействия с другими сервисами системы LMS.
 * <p>
 * Контроллер обрабатывает HTTP-запросы от Python-сервиса Personal Account и других внутренних компонентов системы.
 * Все методы контроллера обеспечивают валидацию входных параметров, обработку ошибок и логирование операций.
 * </p>
 * <p>
 * Основные функции:
 * <ul>
 *   <li>Предоставление информации о попытках прохождения тестов</li>
 *   <li>Агрегация статистики пользователей</li>
 *   <li>Получение данных о тестах и черновиках, связанных с курсами</li>
 * </ul>
 *
 * @author Команда разработки LMS
 * @version 1.0
 * @see InternalApiService
 * @see Context
 */
public class InternalApiController {

    private static final Logger logger = LoggerFactory.getLogger(InternalApiController.class);
    private final InternalApiService internalApiService;

    /**
     * Создает новый экземпляр контроллера внутреннего API.
     *
     * @param internalApiService сервис для обработки бизнес-логики внутреннего API
     * @throws NullPointerException если {@code internalApiService} равен {@code null}
     */
    public InternalApiController(InternalApiService internalApiService) {
        this.internalApiService = internalApiService;
    }

    /**
     * Получает детальную информацию о конкретной попытке прохождения теста.
     * <p>
     * Эндпоинт: {@code GET /internal/attempts/{attemptId}}
     * <p>
     * Ответ содержит полную информацию о попытке, включая результат прохождения и связанные данные.
     *
     * @param ctx HTTP-контекст запроса, содержащий параметры пути и используемый для формирования ответа
     * @throws IllegalArgumentException если параметр {@code attemptId} имеет неверный формат UUID
     * @see AttemptDetail
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
     * Получает список всех попыток прохождения тестов для указанного пользователя.
     * <p>
     * Эндпоинт: {@code GET /internal/users/{userId}/attempts}
     * <p>
     * Ответ содержит список попыток с основной информацией по каждой из них.
     *
     * @param ctx HTTP-контекст запроса, содержащий параметры пути и используемый для формирования ответа
     * @throws IllegalArgumentException если параметр {@code userId} имеет неверный формат UUID
     * @see AttemptsListItem
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
     * Получает агрегированную статистику пользователя по всем попыткам прохождения тестов.
     * <p>
     * Эндпоинт: {@code GET /internal/users/{userId}/stats}
     * <p>
     * Ответ содержит общую статистику пользователя и детальную статистику по каждому тесту.
     *
     * @param ctx HTTP-контекст запроса, содержащий параметры пути и используемый для формирования ответа
     * @throws IllegalArgumentException если параметр {@code userId} имеет неверный формат UUID
     * @see UserStats
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
     * Получает информацию о тесте, связанном с указанным курсом.
     * <p>
     * Эндпоинт: {@code GET /internal/categories/{categoryId}/courses/{courseId}/test}
     * <p>
     * Ответ содержит данные теста или статус "not_found", если тест не найден.
     *
     * @param ctx HTTP-контекст запроса, содержащий параметры пути и используемый для формирования ответа
     * @throws IllegalArgumentException если параметр {@code courseId} имеет неверный формат UUID
     * @see CourseTestResponse
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
     * Получает информацию о черновике теста, связанном с указанным курсом.
     * <p>
     * Эндпоинт: {@code GET /internal/categories/{categoryId}/courses/{courseId}/draft}
     * <p>
     * Ответ содержит данные черновика или статус "not_found", если черновик не найден.
     *
     * @param ctx HTTP-контекст запроса, содержащий параметры пути и используемый для формирования ответа
     * @throws IllegalArgumentException если параметр {@code courseId} имеет неверный формат UUID
     * @see CourseDraftResponse
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

    /**
     * Создает объект для ответа с сообщением об ошибке.
     * <p>
     * Используется для стандартизированного форматирования ошибок в ответах API.
     *
     * @param message текстовое сообщение об ошибке
     * @return объект {@link ErrorResponse}, содержащий сообщение об ошибке
     */
    private ErrorResponse createErrorResponse(String message) {
        return new ErrorResponse(message);
    }

    /**
     * Класс для стандартизированного представления ошибок в ответах API.
     * <p>
     * Все ошибки в API возвращаются в едином формате с полем {@code error},
     * содержащим описание проблемы.
     */
    private static class ErrorResponse {
        
        private final String error;

        /**
         * Создает новый объект ответа с ошибкой.
         *
         * @param error текстовое описание ошибки
         * @throws NullPointerException если {@code error} равен {@code null}
         */
        public ErrorResponse(String error) {
            this.error = error;
        }

        /**
         * Возвращает текстовое описание ошибки.
         *
         * @return сообщение об ошибке
         */
        public String getError() {
            return error;
        }
    }
}