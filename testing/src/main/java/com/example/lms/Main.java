package com.example.lms;

import io.github.cdimascio.dotenv.Dotenv;

import io.javalin.Javalin;
import com.example.lms.test.api.router.TestRouter;
import com.example.lms.test.domain.service.TestService;
import com.example.lms.test.infrastructure.repositories.TestRepository;

import com.example.lms.config.DatabaseConfig;
import com.example.lms.test.api.controller.TestController;

public class Main {
    public static void main(String[] args) {
        Dotenv dotenv = null;
        dotenv = Dotenv.load();

        final Integer APP_PORT = Integer.parseInt(dotenv.get("APP_PORT"));
        final String DB_URL = dotenv.get("DB_URL");
        final String DB_USER = dotenv.get("DB_USER");
        final String DB_PASSWORD = dotenv.get("DB_PASSWORD");

        // Manual Dependency Injection
        DatabaseConfig dbConfig = new DatabaseConfig(DB_URL, DB_USER, DB_PASSWORD);
        TestRepository testRepository = new TestRepository(dbConfig);
        TestService testService = new TestService(testRepository);
        TestController testController = new TestController(testService);

        // Javalin
        Javalin app = Javalin.create(config -> {
            config.router.apiBuilder(() -> {
                TestRouter.register(testController);
            });
        }).start(APP_PORT);
    }
}
