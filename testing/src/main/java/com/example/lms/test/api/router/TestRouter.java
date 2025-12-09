package com.example.lms.test.api.router;

import static io.javalin.apibuilder.ApiBuilder.*;

import com.example.lms.security.JwtHandler;
import com.example.lms.test.api.controller.TestController;
import com.example.lms.tracing.SimpleTracer;
import io.javalin.http.Context;

public class TestRouter {
    public static void register() {
        path("/tests", () -> {
            before(JwtHandler.authenticate());

            get(ctx -> {
                SimpleTracer.runWithSpan("getTests", () -> {
                    captureUserAttributes(ctx);
                    // Добавляем атрибуты из запроса
                    String filter = ctx.queryParam("filter");
                    String sortBy = ctx.queryParam("sortBy");
                    String page = ctx.queryParam("page");
                    
                    if (filter != null) {
                        SimpleTracer.addAttribute("query.filter", filter);
                    }
                    if (sortBy != null) {
                        SimpleTracer.addAttribute("query.sortBy", sortBy);
                    }
                    if (page != null) {
                        SimpleTracer.addAttribute("query.page", page);
                    }
                    
                    // Вызываем ваш контроллер
                    TestController.getTests(ctx);
                    
                    // Добавляем событие
                    SimpleTracer.addEvent("tests.retrieved");
                });
            });
            
            post(ctx -> {
                SimpleTracer.runWithSpan("createTest", () -> {
                    captureUserAttributes(ctx);
                    // Логируем информацию о запросе
                    SimpleTracer.addAttribute("content.type", ctx.contentType());
                    SimpleTracer.addAttribute("content.length", String.valueOf(ctx.contentLength()));
                    
                    // Вызываем ваш контроллер
                    TestController.createTest(ctx);
                    
                    // Добавляем событие
                    SimpleTracer.addEvent("test.created");
                });
            });

            path("/{id}", () -> {
                get(ctx -> {
                    SimpleTracer.runWithSpan("getTestById", () -> {
                        captureUserAttributes(ctx);
                        String testId = ctx.pathParam("id");
                        SimpleTracer.addAttribute("test.id", testId);
                        
                        TestController.getTestById(ctx);
                        
                        SimpleTracer.addEvent("test.retrieved.by.id");
                    });
                });
                
                put(ctx -> {
                    SimpleTracer.runWithSpan("updateTest", () -> {
                        captureUserAttributes(ctx);
                        String testId = ctx.pathParam("id");
                        SimpleTracer.addAttribute("test.id", testId);
                        SimpleTracer.addAttribute("content.type", ctx.contentType());
                        SimpleTracer.addAttribute("content.length", String.valueOf(ctx.contentLength()));
                        
                        TestController.updateTest(ctx);
                        
                        SimpleTracer.addEvent("test.updated");
                    });
                });
                
                delete(ctx -> {
                    SimpleTracer.runWithSpan("deleteTest", () -> {
                        captureUserAttributes(ctx);
                        String testId = ctx.pathParam("id");
                        SimpleTracer.addAttribute("test.id", testId);
                        
                        TestController.deleteTest(ctx);
                        
                        SimpleTracer.addEvent("test.deleted");
                    });
                });
            });
        });
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