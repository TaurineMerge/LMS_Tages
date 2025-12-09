package com.example.lms.test_attempt.api.router;

import static io.javalin.apibuilder.ApiBuilder.*;

import com.example.lms.security.JwtHandler;
import com.example.lms.test_attempt.api.controller.TestAttemptController;
import com.example.lms.tracing.SimpleTracer;
import io.javalin.http.Context;

public class TestAttemptRouter {
    public static void register() {
        path("/test-attempts", () -> {
            before(JwtHandler.authenticate());

            get(ctx -> {
                SimpleTracer.runWithSpan("getTestAttempts", () -> {
                    captureUserAttributes(ctx);
                    // Добавляем атрибуты из запроса
                    String userId = ctx.queryParam("userId");
                    String testId = ctx.queryParam("testId");
                    String status = ctx.queryParam("status");
                    
                    if (userId != null) {
                        SimpleTracer.addAttribute("query.userId", userId);
                    }
                    if (testId != null) {
                        SimpleTracer.addAttribute("query.testId", testId);
                    }
                    if (status != null) {
                        SimpleTracer.addAttribute("query.status", status);
                    }
                    
                    // Вызываем контроллер
                    TestAttemptController.getTestAttempts(ctx);
                    
                    // Добавляем событие
                    SimpleTracer.addEvent("test_attempts.retrieved");
                });
            });
            
            post(ctx -> {
                SimpleTracer.runWithSpan("createTestAttempt", () -> {
                    captureUserAttributes(ctx);
                    // Логируем информацию о запросе
                    SimpleTracer.addAttribute("content.type", ctx.contentType());
                    SimpleTracer.addAttribute("content.length", String.valueOf(ctx.contentLength()));
                    
                    // Вызываем контроллер
                    TestAttemptController.createTestAttempt(ctx);
                    
                    // Добавляем событие
                    SimpleTracer.addEvent("test_attempt.created");
                });
            });

            path("/{id}", () -> {
                get(ctx -> {
                    SimpleTracer.runWithSpan("getTestAttemptById", () -> {
                        captureUserAttributes(ctx);
                        String attemptId = ctx.pathParam("id");
                        SimpleTracer.addAttribute("test_attempt.id", attemptId);
                        
                        TestAttemptController.getTestAttemptById(ctx);
                        
                        SimpleTracer.addEvent("test_attempt.retrieved.by.id");
                    });
                });
                
                delete(ctx -> {
                    SimpleTracer.runWithSpan("deleteTestAttempt", () -> {
                        captureUserAttributes(ctx);
                        String attemptId = ctx.pathParam("id");
                        SimpleTracer.addAttribute("test_attempt.id", attemptId);
                        
                        TestAttemptController.deleteTestAttempt(ctx);
                        
                        SimpleTracer.addEvent("test_attempt.deleted");
                    });
                });
            });
        }); // <-- Закрывающая скобка для path("/test-attempts", () -> {
    }

    private static void captureUserAttributes(Context ctx) {
        Object userId = ctx.attribute("userId");
        Object username = ctx.attribute("username");
        Object email = ctx.attribute("email");
        Object roles = ctx.attribute("roles");

        if (userId != null) {
            SimpleTracer.addAttribute("user.id", userId.toString());
        }
        if (username != null) {
            SimpleTracer.addAttribute("user.username", username.toString());
        }
        if (email != null) {
            SimpleTracer.addAttribute("user.email", email.toString());
        }
        if (roles != null) {
            SimpleTracer.addAttribute("user.roles", roles.toString());
        }
    }
}