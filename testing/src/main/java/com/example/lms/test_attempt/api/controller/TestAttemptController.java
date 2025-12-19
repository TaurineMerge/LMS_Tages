package com.example.lms.test_attempt.api.controller;

import com.example.lms.test_attempt.api.dto.TestAttempt;
import com.example.lms.test_attempt.domain.service.TestAttemptService;
import com.example.lms.tracing.SimpleTracer;

import io.javalin.http.Context;

import java.util.List;
import java.util.Map;
import java.util.UUID;

/**
 * Контроллер для управления попытками прохождения тестов.
 * <p>
 * Обрабатывает HTTP-запросы к ресурсам:
 * <ul>
 * <li>GET /test-attempts — получение списка попыток</li>
 * <li>POST /test-attempts — создание новой попытки</li>
 * <li>GET /test-attempts/{id} — получение попытки по ID</li>
 * <li>PUT /test-attempts/{id} — обновление попытки</li>
 * <li>DELETE /test-attempts/{id} — удаление попытки</li>
 * <li>POST /test-attempts/{id}/complete — завершение попытки</li>
 * <li>PUT /test-attempts/{id}/snapshot — обновление снапшота</li>
 * <li>GET /test-attempts/student/{studentId} — получение попыток студента</li>
 * <li>GET /test-attempts/test/{testId} — получение попыток по тесту</li>
 * </ul>
 *
 * Контроллер использует:
 * <ul>
 * <li>{@link TestAttemptService} — бизнес-логику работы с попытками тестов</li>
 * <li>{@link SimpleTracer} — OpenTelemetry-трейсинг</li>
 * </ul>
 */
public class TestAttemptController {

    /** Сервисный слой, содержащий бизнес-логику работы с попытками тестов. */
    private final TestAttemptService testAttemptService;

    /**
     * Создаёт контроллер попыток тестов.
     *
     * @param testAttemptService сервис управления попытками тестов
     */
    public TestAttemptController(TestAttemptService testAttemptService) {
        this.testAttemptService = testAttemptService;
    }

    /**
     * Обработчик GET /test-attempts.
     * <p>
     * Возвращает список всех попыток тестов.
     */
    public void getTestAttempts(Context ctx) {
        SimpleTracer.runWithSpan("getTestAttempts", () -> {
            List<TestAttempt> attempts = testAttemptService.getAllTestAttempts();
            ctx.json(attempts);
        });
    }

    /**
     * Обработчик POST /test-attempts.
     * <p>
     * Получает JSON с данными попытки, создаёт новую попытку и возвращает её.
     */
    public void createTestAttempt(Context ctx) {
        SimpleTracer.runWithSpan("createTestAttempt", () -> {
            TestAttempt dto = ctx.bodyAsClass(TestAttempt.class);
            TestAttempt created = testAttemptService.createTestAttempt(dto);
            ctx.json(created);
        });
    }

    /**
     * Обработчик GET /test-attempts/{id}.
     * <p>
     * Возвращает попытку по её идентификатору.
     */
    public void getTestAttemptById(Context ctx) {
        SimpleTracer.runWithSpan("getTestAttemptById", () -> {
            String id = ctx.pathParam("id");
            TestAttempt dto = testAttemptService.getTestAttemptById(id);
            ctx.json(dto);
        });
    }

    /**
     * Обработчик PUT /test-attempts/{id}.
     * <p>
     * Обновляет попытку теста.
     */
    public void updateTestAttempt(Context ctx) {
        SimpleTracer.runWithSpan("updateTestAttempt", () -> {
            String id = ctx.pathParam("id");
            TestAttempt dto = ctx.bodyAsClass(TestAttempt.class);
            dto.setId(id);

            TestAttempt updated = testAttemptService.updateTestAttempt(dto);
            ctx.json(updated);
        });
    }

    /**
     * Обработчик DELETE /test-attempts/{id}.
     * <p>
     * Удаляет попытку теста по идентификатору.
     */
    public void deleteTestAttempt(Context ctx) {
        SimpleTracer.runWithSpan("deleteTestAttempt", () -> {
            String id = ctx.pathParam("id");
            boolean deleted = testAttemptService.deleteTestAttempt(id);
            ctx.json(Map.of("deleted", deleted));
        });
    }

    /**
     * Обработчик POST /test-attempts/{id}/complete.
     * <p>
     * Завершает попытку теста с указанием итогового балла.
     */
    public void completeTestAttempt(Context ctx) {
        SimpleTracer.runWithSpan("completeTestAttempt", () -> {
            String id = ctx.pathParam("id");
            Map<String, Object> body = ctx.bodyAsClass(Map.class);
            Integer finalPoint = (Integer) body.get("finalPoint");
            
            if (finalPoint == null) {
                ctx.status(400).json(Map.of("error", "finalPoint is required"));
                return;
            }
            
            TestAttempt completed = testAttemptService.completeTestAttempt(id, finalPoint);
            ctx.json(completed);
        });
    }

    /**
     * Обработчик PUT /test-attempts/{id}/snapshot.
     * <p>
     * Обновляет снапшот попытки теста.
     */
    public void updateSnapshot(Context ctx) {
        SimpleTracer.runWithSpan("updateSnapshot", () -> {
            String id = ctx.pathParam("id");
            Map<String, Object> body = ctx.bodyAsClass(Map.class);
            String snapshot = (String) body.get("snapshot");
            
            if (snapshot == null || snapshot.trim().isEmpty()) {
                ctx.status(400).json(Map.of("error", "snapshot is required"));
                return;
            }
            
            TestAttempt updated = testAttemptService.updateSnapshot(id, snapshot);
            ctx.json(updated);
        });
    }

    /**
     * Обработчик GET /test-attempts/student/{studentId}.
     * <p>
     * Возвращает все попытки указанного студента.
     */
    public void getTestAttemptsByStudentId(Context ctx) {
        SimpleTracer.runWithSpan("getTestAttemptsByStudentId", () -> {
            String studentId = ctx.pathParam("studentId");
            List<TestAttempt> attempts = testAttemptService.getTestAttemptsByStudentId(studentId);
            ctx.json(attempts);
        });
    }

    /**
     * Обработчик GET /test-attempts/test/{testId}.
     * <p>
     * Возвращает все попытки указанного теста.
     */
    public void getTestAttemptsByTestId(Context ctx) {
        SimpleTracer.runWithSpan("getTestAttemptsByTestId", () -> {
            String testId = ctx.pathParam("testId");
            List<TestAttempt> attempts = testAttemptService.getTestAttemptsByTestId(testId);
            ctx.json(attempts);
        });
    }

    /**
     * Внутренний класс для запроса на завершение попытки.
     */
    public static class CompleteRequest {
        private Integer finalPoint;

        public CompleteRequest() {}

        public CompleteRequest(Integer finalPoint) {
            this.finalPoint = finalPoint;
        }

        public Integer getFinalPoint() {
            return finalPoint;
        }

        public void setFinalPoint(Integer finalPoint) {
            this.finalPoint = finalPoint;
        }
    }

    /**
     * Внутренний класс для запроса на обновление снапшота.
     */
    public static class SnapshotRequest {
        private String snapshot;

        public SnapshotRequest() {}

        public SnapshotRequest(String snapshot) {
            this.snapshot = snapshot;
        }

        public String getSnapshot() {
            return snapshot;
        }

        public void setSnapshot(String snapshot) {
            this.snapshot = snapshot;
        }
    }
}