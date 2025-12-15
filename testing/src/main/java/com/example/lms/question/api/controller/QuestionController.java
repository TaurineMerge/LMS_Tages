package com.example.lms.question.api.controller;

import com.example.lms.question.api.dto.Question;
import com.example.lms.question.domain.service.QuestionService;
import com.example.lms.tracing.SimpleTracer;
import io.javalin.http.Context;

import java.util.List;
import java.util.Map;
import java.util.UUID;

/**
 * HTTP-контроллер для работы с вопросами тестов.
 * <p>
 * Отвечает за:
 * <ul>
 *     <li>приём и разбор HTTP-запросов</li>
 *     <li>вызов {@link QuestionService}</li>
 *     <li>формирование HTTP-ответов (JSON)</li>
 * </ul>
 *
 * Все методы предполагается вызывать из Javalin-роутера.
 */
public class QuestionController {

    private final QuestionService questionService;

    /**
     * Создаёт контроллер вопросов.
     *
     * @param questionService сервисный слой для работы с вопросами
     */
    public QuestionController(QuestionService questionService) {
        this.questionService = questionService;
    }

    /**
     * GET /questions
     * <p>
     * Если передан query-параметр {@code testId}, возвращает вопросы только этого теста.
     * Иначе — все вопросы в системе.
     *
     * @param ctx контекст Javalin
     */
    public void getQuestions(Context ctx) {
        SimpleTracer.runWithSpan("getQuestions", () -> {
            String testIdParam = ctx.queryParam("testId");
            List<Question> result;

            if (testIdParam != null && !testIdParam.isBlank()) {
                UUID testId = UUID.fromString(testIdParam);
                SimpleTracer.addAttribute("question.testId", testId.toString());
                result = questionService.getQuestionsByTestId(testId);
            } else {
                result = questionService.getAllQuestions();
            }

            ctx.json(result);
            SimpleTracer.addEvent("questions.list.returned");
        });
    }

    /**
     * GET /questions/{id}
     * <p>
     * Возвращает вопрос по его идентификатору.
     *
     * @param ctx контекст Javalин
     */
    public void getQuestionById(Context ctx) {
        SimpleTracer.runWithSpan("getQuestionById", () -> {
            UUID id = UUID.fromString(ctx.pathParam("id"));
            SimpleTracer.addAttribute("question.id", id.toString());

            Question question = questionService.getQuestionById(id);
            ctx.json(question);

            SimpleTracer.addEvent("question.returned.by.id");
        });
    }

    /**
     * POST /questions
     * <p>
     * Создаёт новый вопрос.
     * Ожидает JSON тела, соответствующего {@link Question}.
     *
     * @param ctx контекст Javalin
     */
    public void createQuestion(Context ctx) {
        SimpleTracer.runWithSpan("createQuestion", () -> {
            Question dto = ctx.bodyAsClass(Question.class);

            if (dto.getId() != null) {
                // На всякий случай игнорируем присланный извне id
                dto.setId(null);
            }

            Question created = questionService.createQuestion(dto);
            ctx.status(201).json(created);

            SimpleTracer.addAttribute("question.testId", created.getTestId().toString());
            SimpleTracer.addEvent("question.created");
        });
    }

    /**
     * PUT /questions/{id}
     * <p>
     * Обновляет существующий вопрос.
     * ID берётся из path-параметра, тело — из JSON.
     *
     * @param ctx контекст Javalin
     */
    public void updateQuestion(Context ctx) {
        SimpleTracer.runWithSpan("updateQuestion", () -> {
            UUID id = UUID.fromString(ctx.pathParam("id"));
            Question dto = ctx.bodyAsClass(Question.class);
            dto.setId(id);

            SimpleTracer.addAttribute("question.id", id.toString());

            Question updated = questionService.updateQuestion(dto);
            ctx.json(updated);

            SimpleTracer.addEvent("question.updated");
        });
    }

    /**
     * DELETE /questions/{id}
     * <p>
     * Удаляет вопрос по ID. Возвращает JSON с признаком удаления.
     *
     * @param ctx контекст Javalin
     */
    public void deleteQuestion(Context ctx) {
        SimpleTracer.runWithSpan("deleteQuestion", () -> {
            UUID id = UUID.fromString(ctx.pathParam("id"));
            SimpleTracer.addAttribute("question.id", id.toString());

            boolean deleted = questionService.deleteQuestion(id);
            ctx.json(Map.of(
                    "id", id.toString(),
                    "deleted", deleted
            ));

            SimpleTracer.addEvent("question.deleted");
        });
    }
}