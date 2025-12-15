package com.example.lms.test.api.controller;

import com.example.lms.test.api.dto.Test;
import io.javalin.http.Context;
import java.util.Map;

public class TestController {

    // GET /tests
    public static void getTests(Context ctx) {
        // TODO: Заглушка, пока нет подключения к БД
        ctx.json("[]");
    }

    // POST /tests
    public static void createTest(Context ctx) {
        // TODO: Заглушка, пока нет подключения к БД
        Test dto = ctx.bodyAsClass(Test.class);
        ctx.json(dto);
    }

    // GET /tests/{id}
    public static void getTestById(Context ctx) {
        // TODO: Заглушка, пока нет подключения к БД
        String id = ctx.pathParam("id");
        ctx.json(Map.of("id", id, "message", "Not implemented yet"));
    }

    // PUT /tests/{id}
    public static void updateTest(Context ctx) {
        // TODO: Заглушка, пока нет подключения к БД
        String id = ctx.pathParam("id");
        Test dto = ctx.bodyAsClass(Test.class);
        ctx.json(dto);
    }

    // DELETE /tests/{id}
    public static void deleteTest(Context ctx) {
        // TODO: Заглушка, пока нет подключения к БД
        String id = ctx.pathParam("id");
        ctx.json(Map.of("deleted", true, "id", id));
    }
}
