package com.example.controller;

import io.javalin.http.Context;

public class TestAttemptController {
    
    public static void getTestAttempts(Context ctx) {
        ctx.json("GET /test-attempts");
    }

    public static void createTestAttempt(Context ctx) {
        ctx.json("POST /test-attempts");
    }

    public static void getTestAttemptById(Context ctx) {
        ctx.json("GET /test-attempts/" + ctx.pathParam("id"));
    }

    public static void deleteTestAttempt(Context ctx) {
        ctx.json("DELETE /test-attempts/" + ctx.pathParam("id"));
    }
}