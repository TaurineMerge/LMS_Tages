package com.example.lms.test.api.controller;

import io.javalin.http.Context;

public class TestController {

    public static void getTests(Context ctx) {
        ctx.json("GET /tests");
    }

    public static void createTest(Context ctx) {
        ctx.json("POST /tests");
    }

    public static void getTestById(Context ctx) {
        ctx.json("GET /tests/" + ctx.pathParam("id"));
    }

    public static void updateTest(Context ctx) {
        ctx.json("PUT /tests/" + ctx.pathParam("id"));
    }

    public static void deleteTest(Context ctx) {
        ctx.json("DELETE /tests/" + ctx.pathParam("id"));
    }
}