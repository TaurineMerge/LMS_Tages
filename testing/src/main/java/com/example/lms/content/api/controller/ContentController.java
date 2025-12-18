package com.example.lms.content.api.controller;

import com.example.lms.content.api.dto.Content;
import com.example.lms.content.domain.service.ContentService;
import com.github.jknack.handlebars.Handlebars;
import com.github.jknack.handlebars.Template;

import com.example.lms.tracing.SimpleTracer;

import io.javalin.http.Context;

import java.io.StringWriter;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Контроллер для управления элементами контента.
 * <p>
 * Обрабатывает HTTP-запросы к ресурсам:
 * <ul>
 * <li>GET /contents — получение списка элементов контента</li>
 * <li>GET /contents/html — получение списка в HTML-шаблоне</li>
 * <li>POST /contents — создание нового элемента контента</li>
 * <li>GET /contents/{id} — получение элемента контента по ID</li>
 * <li>PUT /contents/{id} — обновление элемента контента</li>
 * <li>DELETE /contents/{id} — удаление элемента контента</li>
 * <li>GET /contents/search?content=... — поиск по содержимому</li>
 * <li>GET /contents/type?type=... — фильтрация по типу контента</li>
 * <li>GET /contents/question/{questionId} — поиск по вопросу</li>
 * <li>GET /contents/answer/{answerId} — поиск по ответу</li>
 * </ul>
 *
 * Контроллер использует:
 * <ul>
 * <li>{@link ContentService} — бизнес-логику работы с элементами контента</li>
 * <li>{@link Handlebars} — генерацию HTML-шаблонов</li>
 * <li>{@link SimpleTracer} — OpenTelemetry-трейсинг</li>
 * </ul>
 */
public class ContentController {

    /** Сервисный слой, содержащий бизнес-логику работы с элементами контента. */
    private final ContentService contentService;

    /** Шаблонизатор Handlebars для генерации HTML. */
    private final Handlebars handlebars;

    /**
     * Создаёт контроллер элементов контента.
     *
     * @param contentService сервис управления элементами контента
     * @param handlebars шаблонизатор Handlebars
     */
    public ContentController(ContentService contentService, Handlebars handlebars) {
        this.contentService = contentService;
        this.handlebars = handlebars;
    }

    /**
     * Метод для рендеринга шаблонов.
     *
     * @param ctx контекст HTTP-запроса
     * @param templatePath путь к шаблону
     * @param model модель данных для шаблона
     */
    private void renderTemplate(Context ctx, String templatePath, Map<String, Object> model) {
        try {
            Template template = handlebars.compile(templatePath);
            StringWriter writer = new StringWriter();
            template.apply(model, writer);
            ctx.contentType("text/html; charset=utf-8");
            ctx.result(writer.toString());
        } catch (Exception e) {
            e.printStackTrace();
            ctx.status(500).result("Ошибка рендеринга шаблона: " + e.getMessage());
        }
    }

    /**
     * Обработчик GET /contents/html.
     * <p>
     * Возвращает HTML-страницу со списком элементов контента.
     *
     * @param ctx контекст HTTP-запроса
     */
    public void getContentsHtml(Context ctx) {
        SimpleTracer.runWithSpan("getContentsHtml", () -> {
            List<Content> contents = contentService.getAllContents();
            Map<String, Object> model = new HashMap<>();
            model.put("contents", contents);
            model.put("title", "Список элементов контента");

            renderTemplate(ctx, "layouts/content-list", model);
        });
    }

    /**
     * Обработчик GET /contents.
     * <p>
     * Возвращает JSON со списком всех элементов контента.
     *
     * @param ctx контекст HTTP-запроса
     */
    public void getAllContents(Context ctx) {
        SimpleTracer.runWithSpan("getAllContents", () -> {
            List<Content> contents = contentService.getAllContents();
            ctx.json(contents);
        });
    }

    /**
     * Обработчик POST /contents.
     * <p>
     * Получает JSON с данными элемента контента, создаёт новый элемент и возвращает его.
     *
     * @param ctx контекст HTTP-запроса
     */
    public void createContent(Context ctx) {
        SimpleTracer.runWithSpan("createContent", () -> {
            Content dto = ctx.bodyAsClass(Content.class);
            
            // Валидация данных
            try {
                contentService.validateContentData(dto);
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", e.getMessage()));
                return;
            }
            
            Content created = contentService.createContent(dto);
            ctx.status(201).json(created);
        });
    }

    /**
     * Обработчик GET /contents/{id}.
     * <p>
     * Возвращает элемент контента по его идентификатору.
     *
     * @param ctx контекст HTTP-запроса
     */
    public void getContentById(Context ctx) {
        SimpleTracer.runWithSpan("getContentById", () -> {
            String id = ctx.pathParam("id");
            try {
                Content dto = contentService.getContentById(id);
                ctx.json(dto);
            } catch (java.util.NoSuchElementException e) {
                ctx.status(404).json(Map.of("error", e.getMessage()));
            }
        });
    }

    /**
     * Обработчик PUT /contents/{id}.
     * <p>
     * Обновляет элемент контента: данные принимаются в JSON, а ID берётся из path parameter.
     *
     * @param ctx контекст HTTP-запроса
     */
    public void updateContent(Context ctx) {
        SimpleTracer.runWithSpan("updateContent", () -> {
            String id = ctx.pathParam("id");
            Content dto = ctx.bodyAsClass(Content.class);
            dto.setId(id);

            // Валидация данных
            try {
                contentService.validateContentData(dto);
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", e.getMessage()));
                return;
            }

            try {
                Content updated = contentService.updateContent(dto);
                ctx.json(updated);
            } catch (java.util.NoSuchElementException e) {
                ctx.status(404).json(Map.of("error", e.getMessage()));
            } catch (IllegalArgumentException e) {
                ctx.status(400).json(Map.of("error", e.getMessage()));
            }
        });
    }

    /**
     * Обработчик DELETE /contents/{id}.
     * <p>
     * Удаляет элемент контента по идентификатору.
     *
     * @param ctx контекст HTTP-запроса
     */
    public void deleteContent(Context ctx) {
        SimpleTracer.runWithSpan("deleteContent", () -> {
            String id = ctx.pathParam("id");
            boolean deleted = contentService.deleteContent(id);
            if (deleted) {
                ctx.status(204);
            } else {
                ctx.status(404).json(Map.of("error", "Элемент контента с ID " + id + " не найден"));
            }
        });
    }

    /**
     * Обработчик GET /contents/search.
     * <p>
     * Поиск элементов контента по содержимому.
     *
     * @param ctx контекст HTTP-запроса
     */
    public void searchByContent(Context ctx) {
        SimpleTracer.runWithSpan("searchByContent", () -> {
            String content = ctx.queryParam("content");
            if (content == null || content.trim().isEmpty()) {
                ctx.status(400).json(Map.of("error", "Параметр 'content' обязателен"));
                return;
            }
            
            List<Content> contents = contentService.findByContentContaining(content);
            ctx.json(contents);
        });
    }

    /**
     * Обработчик GET /contents/type.
     * <p>
     * Фильтрация элементов контента по типу.
     *
     * @param ctx контекст HTTP-запроса
     */
    public void filterByType(Context ctx) {
        SimpleTracer.runWithSpan("filterByType", () -> {
            String typeParam = ctx.queryParam("type");
            Boolean type = null;
            
            if (typeParam != null) {
                if ("true".equalsIgnoreCase(typeParam)) {
                    type = true;
                } else if ("false".equalsIgnoreCase(typeParam)) {
                    type = false;
                } else {
                    ctx.status(400).json(Map.of("error", "Параметр 'type' должен быть 'true' или 'false'"));
                    return;
                }
            }
            
            List<Content> contents = contentService.findByTypeOfContent(type);
            ctx.json(contents);
        });
    }

    /**
     * Обработчик GET /contents/question/{questionId}.
     * <p>
     * Поиск элементов контента по ID вопроса.
     *
     * @param ctx контекст HTTP-запроса
     */
    public void getByQuestionId(Context ctx) {
        SimpleTracer.runWithSpan("getByQuestionId", () -> {
            String questionId = ctx.pathParam("questionId");
            List<Content> contents = contentService.findByQuestionId(questionId);
            ctx.json(contents);
        });
    }

    /**
     * Обработчик GET /contents/answer/{answerId}.
     * <p>
     * Поиск элементов контента по ID ответа.
     *
     * @param ctx контекст HTTP-запроса
     */
    public void getByAnswerId(Context ctx) {
        SimpleTracer.runWithSpan("getByAnswerId", () -> {
            String answerId = ctx.pathParam("answerId");
            List<Content> contents = contentService.findByAnswerId(answerId);
            ctx.json(contents);
        });
    }

    /**
     * Обработчик GET /contents/validate/{id}.
     * <p>
     * Проверяет существование элемента контента по ID.
     *
     * @param ctx контекст HTTP-запроса
     */
    public void existsById(Context ctx) {
        SimpleTracer.runWithSpan("existsById", () -> {
            String id = ctx.pathParam("id");
            boolean exists = contentService.existsById(id);
            ctx.json(Map.of("exists", exists));
        });
    }

    /**
     * Обработчик GET /contents/max-order.
     * <p>
     * Возвращает элементы контента с максимальным порядковым номером.
     *
     * @param ctx контекст HTTP-запроса
     */
    public void getWithMaxOrder(Context ctx) {
        SimpleTracer.runWithSpan("getWithMaxOrder", () -> {
            List<Content> contents = contentService.findWithMaxOrder();
            ctx.json(contents);
        });
    }
}