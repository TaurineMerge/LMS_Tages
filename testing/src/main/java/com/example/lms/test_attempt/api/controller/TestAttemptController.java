package com.example.lms.test_attempt.api.controller;

import io.javalin.http.Context;

/**
 * Контроллер для работы с попытками прохождения тестов.
 * <p>
 * На текущем этапе методы являются упрощёнными заглушками
 * и возвращают текстовое описание вызванного действия.
 * <p>
 * Предполагаемые реальные эндпоинты:
 * <ul>
 *     <li>GET /test-attempts — получение списка попыток</li>
 *     <li>POST /test-attempts — создание новой попытки</li>
 *     <li>GET /test-attempts/{id} — получение попытки по ID</li>
 *     <li>DELETE /test-attempts/{id} — удаление попытки</li>
 * </ul>
 *
 * Позже контроллер будет использовать сервисный слой и доменные модели,
 * как это реализовано в модуле тестов.
 */
public class TestAttemptController {

    /**
     * Обработчик GET /test-attempts.
     * <p>
     * На текущем этапе возвращает заглушку:
     * <pre>"GET /test-attempts"</pre>
     *
     * @param ctx HTTP-контекст Javalin
     */
    public static void getTestAttempts(Context ctx) {
        ctx.json("GET /test-attempts");
    }

    /**
     * Обработчик POST /test-attempts.
     * <p>
     * Пока что возвращает заглушку:
     * <pre>"POST /test-attempts"</pre>
     *
     * @param ctx HTTP-контекст Javalин
     */
    public static void createTestAttempt(Context ctx) {
        ctx.json("POST /test-attempts");
    }

    /**
     * Обработчик GET /test-attempts/{id}.
     * <p>
     * Возвращает простую строку с ID попытки.
     * <p>
     * Пример:
     * <pre>GET /test-attempts/123</pre>
     *
     * @param ctx HTTP-контекст
     */
    public static void getTestAttemptById(Context ctx) {
        ctx.json("GET /test-attempts/" + ctx.pathParam("id"));
    }

    /**
     * Обработчик DELETE /test-attempts/{id}.
     * <p>
     * На текущем этапе отвечает только текстом заглушкой.
     *
     * @param ctx HTTP-контекст
     */
    public static void deleteTestAttempt(Context ctx) {
        ctx.json("DELETE /test-attempts/" + ctx.pathParam("id"));
    }
}