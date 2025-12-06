package com.example.lms;

import io.javalin.Javalin;

import com.example.lms.test.api.router.TestRouter;
import com.example.lms.test_attempt.api.router.TestAttemptRouter;

public class Main {

    public static void main(String[] args) {
        var app = Javalin.create(config -> {
            config.router.apiBuilder(() -> {
                TestRouter.register();
                TestAttemptRouter.register();
            });
        }).start(8085);
    }
}