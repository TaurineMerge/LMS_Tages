package com.example.lms.test_attempt.api.controller;

import io.javalin.http.Context;
import java.util.Map;

/**
 * Контроллер для работы с попытками прохождения тестов (TestAttempt).
 * <p>
 * Данный класс содержит HTTP-эндпоинты для CRUD-операций над попытками.
 * На текущем этапе методы являются заглушками и не имеют подключения к БД.
 * Позже будут реализованы через сервисный слой и репозиторий.
 *
 * Эндпоинты:
 * <ul>
 *     <li>GET /test-attempts — получение списка попыток</li>
 *     <li>POST /test-attempts — создание новой попытки</li>
 *     <li>GET /test-attempts/{id} — получение попытки по ID</li>
 *     <li>DELETE /test-attempts/{id} — удаление попытки</li>
 * </ul>
 */
public class TestAttemptController {

    /**
     * Обработчик GET /test-attempts.
     * <p>
     * Возвращает список попыток тестирования.
     * Пока реализована заглушка: метод возвращает пустой массив JSON.
     *
     * @param ctx контекст запроса Javalин
     */
    public static void getTestAttempts(Context ctx) {
        ctx.json("[]");
    }

    /**
     * Обработчик POST /test-attempts.
     * <p>
     * Создаёт новую попытку тестирования.
     * На текущем этапе возвращает текстовую заглушку без обработки входных данных.
     *
     * Пример будущей логики:
     * <pre>
     * TestAttempt dto = ctx.bodyAsClass(TestAttempt.class);
     * ctx.json(dto);
     * </pre>
     *
     * @param ctx контекст запроса Javalин
     */
    public static void createTestAttempt(Context ctx) {
        ctx.json("POST /test-attempts (stub)");
    }

    /**
     * Обработчик GET /test-attempts/{id}.
     * <p>
     * Возвращает данные о конкретной попытке по её идентификатору.
     * Пока возвращает заглушку с ID.
     *
     * @param ctx контекст запроса Javalин
     */
    public static void getTestAttemptById(Context ctx) {
        String id = ctx.pathParam("id");
        ctx.json(Map.of("id", id, "message", "Not implemented yet"));
    }

    /**
     * Обработчик DELETE /test-attempts/{id}.
     * <p>
     * Удаляет попытку тестирования по ID.
     * Пока метод не взаимодействует с БД и возвращает заглушку.
     *
     * @param ctx контекст запроса Javalин
     */
    public static void deleteTestAttempt(Context ctx) {
        String id = ctx.pathParam("id");
        ctx.json(Map.of("deleted", true, "id", id));
    }
}