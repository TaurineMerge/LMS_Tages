package com.example.lms.tracing;

import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.context.Scope;

/**
 * Утилитный класс для работы с OpenTelemetry-трейсингом.
 * <p>
 * Предоставляет:
 * <ul>
 *     <li>ручное создание спанов</li>
 *     <li>логирование traceId</li>
 *     <li>добавление атрибутов различных типов</li>
 *     <li>добавление событий в текущий спан</li>
 * </ul>
 *
 * Используется в роутерах и сервисах, чтобы:
 * <ul>
 *     <li>улучшать наблюдаемость</li>
 *     <li>детализировать операции</li>
 *     <li>передавать контекст запроса в мониторинг</li>
 * </ul>
 *
 * Все методы статические, что делает класс удобным
 * для использования как лёгкой обёртки вокруг OTel API.
 */
public class SimpleTracer {

    /**
     * Tracer, предоставляемый OpenTelemetry Java Agent.
     * <p>
     * Java Agent автоматически настраивает SDK, поэтому
     * объект берётся через {@link GlobalOpenTelemetry#getTracer(String)}.
     */
    private static final Tracer TRACER = GlobalOpenTelemetry.getTracer("lms-service");

    /**
     * Возвращает traceId текущего спана.
     * <p>
     * Используется в логах, чтобы связать логи и трейс.
     *
     * @return traceId или строку "no-trace-id", если спан отсутствует
     */
    public static String getCurrentTraceId() {
        Span current = Span.current();
        if (current != null && current.getSpanContext().isValid()) {
            return current.getSpanContext().getTraceId();
        }
        return "no-trace-id";
    }

    /**
     * Создаёт новый спан, выполняет переданную операцию и корректно завершает спан.
     * <p>
     * Использование:
     * <pre>
     * SimpleTracer.runWithSpan("loadUser", () -> {
     *     ... логика ...
     * });
     * </pre>
     *
     * В случае ошибки:
     * <ul>
     *     <li>ошибка записывается в спан</li>
     *     <li>исключение пробрасывается вверх</li>
     * </ul>
     *
     * @param spanName имя создаваемого спана
     * @param operation выполняемая операция
     */
    public static void runWithSpan(String spanName, Runnable operation) {
        Span span = TRACER.spanBuilder(spanName).startSpan();
        try (Scope scope = span.makeCurrent()) {
            operation.run();
        } catch (Exception e) {
            span.recordException(e);
            throw e;
        } finally {
            span.end();
        }
    }

    /**
     * Добавляет строковый атрибут в текущий спан.
     *
     * @param key ключ атрибута
     * @param value значение
     */
    public static void addAttribute(String key, String value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавляет числовой атрибут типа int.
     */
    public static void addAttribute(String key, int value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавляет числовой атрибут типа long.
     */
    public static void addAttribute(String key, long value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавляет числовой атрибут типа double.
     */
    public static void addAttribute(String key, double value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавляет булевый атрибут.
     */
    public static void addAttribute(String key, boolean value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавляет объект как строковый атрибут (вызывает {@code toString()}).
     *
     * @param key ключ
     * @param value объект для записи
     */
    public static void addAttribute(String key, Object value) {
        Span.current().setAttribute(key, value.toString());
    }

    /**
     * Добавляет событие (event) в текущий спан.
     * <p>
     * Используется для фиксации важных бизнес-событий:
     * <ul>
     *     <li>"test_attempt.created"</li>
     *     <li>"test_attempt.retrieved.by.id"</li>
     *     <li>"user.authenticated"</li>
     * </ul>
     *
     * @param name имя события
     */
    public static void addEvent(String name) {
        Span.current().addEvent(name);
    }
}