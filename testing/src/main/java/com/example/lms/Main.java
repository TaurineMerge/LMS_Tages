package com.example.lms;

import io.javalin.Javalin;
import com.example.lms.test.api.router.TestRouter;
import com.example.lms.test.domain.service.TestService;
import com.example.lms.test.infrastructure.repositories.TestRepository;
import com.example.lms.config.DatabaseConfig;
import com.example.lms.test.api.controller.TestController;

public class Main {

    public static void main(String[] args) {
        // Manual Dependency Injection
        DatabaseConfig dbConfig = new DatabaseConfig("jdbc:postgresql://app-db/appdb", "appuser", "password");
        TestRepository testRepository = new TestRepository(dbConfig);
        TestService testService = new TestService(testRepository);
        TestController testController = new TestController(testService);

        // Javalin
        Javalin app = Javalin.create(config -> {
            config.router.apiBuilder(() -> {
                TestRouter.register(testController);
            });
        }).start(8085);
    }
}
