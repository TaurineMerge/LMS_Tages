package com.example.lms.test_attempt.api.controller;

import io.javalin.http.Context;
import java.util.Map;

public class TestAttemptController {

    public static void getTestAttempts(Context ctx) {
        // TODO: Заглушка, пока нет подключения к БД
        ctx.json("[]");
    }

    public static void createTestAttempt(Context ctx) {
        // TODO: Заглушка, пока нет подключения к БД
        // Предположим, что есть DTO класс TestAttempt
        // TestAttempt dto = ctx.bodyAsClass(TestAttempt.class);
        // ctx.json(dto);
        ctx.json("POST /test-attempts (stub)");
    }

    public static void getTestAttemptById(Context ctx) {
        // TODO: Заглушка, пока нет подключения к БД
        String id = ctx.pathParam("id");
        ctx.json(Map.of("id", id, "message", "Not implemented yet"));
    }

    public static void deleteTestAttempt(Context ctx) {
        // TODO: Заглушка, пока нет подключения к БД
        String id = ctx.pathParam("id");
        ctx.json(Map.of("deleted", true, "id", id));
    }
}