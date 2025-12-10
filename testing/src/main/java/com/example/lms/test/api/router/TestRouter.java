package com.example.lms.test.api.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.security.JwtHandler;
import com.example.lms.test.api.controller.TestController;
import com.example.lms.tracing.SimpleTracer;

import static io.javalin.apibuilder.ApiBuilder.after;
import static io.javalin.apibuilder.ApiBuilder.before;
import static io.javalin.apibuilder.ApiBuilder.delete;
import static io.javalin.apibuilder.ApiBuilder.get;
import static io.javalin.apibuilder.ApiBuilder.path;
import static io.javalin.apibuilder.ApiBuilder.post;
import static io.javalin.apibuilder.ApiBuilder.put;

public class TestRouter {
    private static final Logger logger = LoggerFactory.getLogger(TestRouter.class);

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