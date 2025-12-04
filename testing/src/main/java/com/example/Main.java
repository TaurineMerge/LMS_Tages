package com.example;

import io.javalin.Javalin;
import io.javalin.http.Context;

import static io.javalin.apibuilder.ApiBuilder.*;

public class Main {

    public static void main(String[] args) {
        var app = Javalin.create(config -> {
            config.router.apiBuilder(() -> {
                registerTestRoutes();
                registerTestAttemptRoutes();
            });
        }).start(8085);
    }

    // Группировка эндпоинтов для Tests
    private static void registerTestRoutes() {
        path("/tests", () -> {
            get(Main::getTests);
            post(Main::createTest);
            
            path("/{id}", () -> {
                get(Main::getTestById);
                put(Main::updateTest);
                delete(Main::deleteTest);
            });
        });
    }

    // Группировка эндпоинтов для Test Attempts
    private static void registerTestAttemptRoutes() {
        path("/test-attempts", () -> {
            get(Main::getTestAttempts);
            post(Main::createTestAttempt);
            
            path("/{id}", () -> {
                get(Main::getTestAttemptById);
                delete(Main::deleteTestAttempt);
            });
        });
    }

    // Tests handlers (заглушки)
    private static void getTests(Context ctx) {
        ctx.json("GET /tests");
    }

    private static void createTest(Context ctx) {
        ctx.json("POST /tests");
    }

    private static void getTestById(Context ctx) {
        ctx.json("GET /tests/" + ctx.pathParam("id"));
    }

    private static void updateTest(Context ctx) {
        ctx.json("PUT /tests/" + ctx.pathParam("id"));
    }

    private static void deleteTest(Context ctx) {
        ctx.json("DELETE /tests/" + ctx.pathParam("id"));
    }

    // Test Attempts handlers (заглушки)
    private static void getTestAttempts(Context ctx) {
        ctx.json("GET /test-attempts");
    }

    private static void createTestAttempt(Context ctx) {
        ctx.json("POST /test-attempts");
    }

    private static void getTestAttemptById(Context ctx) {
        ctx.json("GET /test-attempts/" + ctx.pathParam("id"));
    }

    private static void deleteTestAttempt(Context ctx) {
        ctx.json("DELETE /test-attempts/" + ctx.pathParam("id"));
    }
}