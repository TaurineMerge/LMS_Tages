package com.example.lms.test.api.controller;

import com.example.lms.test.api.dto.Test;
import com.example.lms.test.domain.service.TestService;
import com.github.jknack.handlebars.Handlebars;
import com.github.jknack.handlebars.Template;

import com.example.lms.tracing.SimpleTracer;

import io.javalin.http.Context;

import java.io.IOException;
import java.util.Map;

/**
 * Контроллер для управления тестами.
 * <p>
 * Обрабатывает HTTP-запросы к ресурсам:
 * <ul>
 *     <li>GET /tests — получение списка тестов в HTML-шаблоне</li>
 *     <li>POST /tests — создание нового теста</li>
 *     <li>GET /tests/{id} — получение теста по ID</li>
 *     <li>PUT /tests/{id} — обновление теста</li>
 *     <li>DELETE /tests/{id} — удаление теста</li>
 * </ul>
 *
 * Контроллер использует:
 * <ul>
 *     <li>{@link TestService} — бизнес-логику работы с тестами</li>
 *     <li>{@link Handlebars} — генерацию HTML-шаблонов</li>
 *     <li>{@link SimpleTracer} — OpenTelemetry-трейсинг</li>
 * </ul>
 */
public class TestController {

    /** Сервисный слой, содержащий бизнес-логику работы с тестами. */
    private final TestService testService;

    /** Шаблонизатор Handlebars для генерации HTML. */
    private final Handlebars handlebars = new Handlebars();

    /**
     * Создаёт контроллер тестов.
     *
     * @param testService сервис управления тестами
     */
    public TestController(TestService testService) {
        this.testService = testService;
    }

    /**
     * Обработчик GET /tests.
     * <p>
     * Возвращает HTML-страницу со списком тестов, используя шаблон
     * {@code /templates/test-list.hbs}.
     * <p>
     * Если при компиляции или применении шаблона возникает ошибка —
     * клиенту возвращается статус 500.
     *
     * @param ctx контекст HTTP-запроса
     */
    public void getTests(Context ctx) {
        SimpleTracer.runWithSpan("getTests", () -> {
            try {
                var tests = testService.getAllTests();

                Template template = handlebars.compile("/templates/test-list");
                String html = template.apply(Map.of("tests", tests));

                ctx.html(html);
            } catch (IOException e) {
                e.printStackTrace();
                ctx.status(500).html("Internal Server Error: " + e.getMessage());
            }
        });
    }

    /**
     * Обработчик POST /tests.
     * <p>
     * Получает JSON с данными теста, создаёт новый тест и возвращает его.
     *
     * @param ctx HTTP-контекст
     */
    public void createTest(Context ctx) {
        SimpleTracer.runWithSpan("createTest", () -> {
            Test dto = ctx.bodyAsClass(Test.class);
            Test created = testService.createTest(dto);
            ctx.json(created);
        });
    }

    /**
     * Обработчик GET /tests/{id}.
     * <p>
     * Возвращает тест по его идентификатору.
     *
     * @param ctx HTTP-контекст
     */
    public void getTestById(Context ctx) {
        SimpleTracer.runWithSpan("getTestById", () -> {
            String id = ctx.pathParam("id");
            Test dto = testService.getTestById(id);
            ctx.json(dto);
        });
    }

    /**
     * Обработчик PUT /tests/{id}.
     * <p>
     * Обновляет тест: данные принимаются в JSON, а ID берётся из path parameter.
     *
     * @param ctx HTTP-контекст
     */
    public void updateTest(Context ctx) {
        SimpleTracer.runWithSpan("updateTest", () -> {
            String id = ctx.pathParam("id");
            Test dto = ctx.bodyAsClass(Test.class);
            dto.setId(Long.parseLong(id));

            Test updated = testService.updateTest(dto);
            ctx.json(updated);
        });
    }

    /**
     * Обработчик DELETE /tests/{id}.
     * <p>
     * Удаляет тест по идентификатору и возвращает JSON вида:
     * <pre>
     * {"deleted": true}
     * </pre>
     *
     * @param ctx HTTP-контекст
     */
    public void deleteTest(Context ctx) {
        SimpleTracer.runWithSpan("deleteTest", () -> {
            String id = ctx.pathParam("id");
            boolean deleted = testService.deleteTest(id);
            ctx.json(Map.of("deleted", deleted));
        });
    }
}