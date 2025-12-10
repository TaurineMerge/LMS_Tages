package com.example.lms.test.api.router;

import com.example.lms.test.api.controller.TestController;
import com.example.lms.tracing.SimpleTracer;
import com.example.lms.Main;
import com.example.lms.security.JwtHandler;
import static io.javalin.apibuilder.ApiBuilder.*;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class TestRouter {
    private static final Logger logger = LoggerFactory.getLogger(Main.class);

    public static void register(TestController testController) {
        path("/tests", () -> {
            before(JwtHandler.authenticate());
            before(ctx -> {
                logger.info("▶️  Request started: {} {} (traceId: {})",
                        ctx.method(),
                        ctx.path(),
                        SimpleTracer.getCurrentTraceId());
            });

            get(testController::getTests);
            post(testController::createTest);

            path("/{id}", () -> {
                get(testController::getTestById);
                put(testController::updateTest);
                delete(testController::deleteTest);
            });

            after(ctx -> {
                logger.info("✅ Request completed: {} {} -> {} (traceId: {})",
                        ctx.method(),
                        ctx.path(),
                        ctx.status(),
                        SimpleTracer.getCurrentTraceId());
            });
        });
    }
}