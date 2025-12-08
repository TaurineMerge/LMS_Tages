package com.example.lms.test.api.controller;

import com.example.lms.test.api.dto.Test;
import com.example.lms.test.domain.service.TestService;
import com.github.jknack.handlebars.Handlebars;
import com.github.jknack.handlebars.Template;

import io.javalin.http.Context;

import java.io.IOException;
import java.util.Map;

public class TestController {

    private final TestService testService;
    private final Handlebars handlebars = new Handlebars();

    public TestController(TestService testService) {
        this.testService = testService;
    }

    // GET /tests
    public void getTests(Context ctx) {
        try {
            var tests = testService.getAllTests();
            Template template = handlebars.compile("/templates/test-list");
            String html = template.apply(Map.of("tests", tests));
            ctx.html(html);
        } catch (IOException e) {
            e.printStackTrace();
            ctx.status(500).html("Internal Server Error: " + e.getMessage());
        }
    }

    // POST /tests
    public void createTest(Context ctx) {
        Test dto = ctx.bodyAsClass(Test.class);
        Test created = testService.createTest(dto);
        ctx.json(created);
    }

    // GET /tests/{id}
    public void getTestById(Context ctx) {
        String id = ctx.pathParam("id");
        Test dto = testService.getTestById(id);
        ctx.json(dto);
    }

    // PUT /tests/{id}
    public void updateTest(Context ctx) {
        String id = ctx.pathParam("id");
        Test dto = ctx.bodyAsClass(Test.class);
        dto.setId(Long.parseLong(id));
        Test updated = testService.updateTest(dto);
        ctx.json(updated);
    }

    // DELETE /tests/{id}
    public void deleteTest(Context ctx) {
        String id = ctx.pathParam("id");
        boolean deleted = testService.deleteTest(id);
        ctx.json(Map.of("deleted", deleted));
    }
}
