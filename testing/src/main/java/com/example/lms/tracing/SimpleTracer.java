package com.example.lms.tracing;

import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.context.Scope;

public class SimpleTracer {

    private static final Tracer TRACER = GlobalOpenTelemetry.getTracer("lms-service");

    /**
     * Получить текущий traceId для логгирования
     */
    public static String getCurrentTraceId() {
        Span current = Span.current();
        if (current != null && current.getSpanContext().isValid()) {
            return current.getSpanContext().getTraceId();
        }
        return "no-trace-id";
    }

    /**
     * Метод для ручного создания спанов
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
     * Добавить атрибут к текущему спану (String)
     */
    public static void addAttribute(String key, String value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавить атрибут к текущему спану (int)
     */
    public static void addAttribute(String key, int value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавить атрибут к текущему спану (long)
     */
    public static void addAttribute(String key, long value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавить атрибут к текущему спану (double)
     */
    public static void addAttribute(String key, double value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавить атрибут к текущему спану (boolean)
     */
    public static void addAttribute(String key, boolean value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавить атрибут к текущему спану (Object - вызывает toString)
     */
    public static void addAttribute(String key, Object value) {
        Span.current().setAttribute(key, value.toString());
    }

    /**
     * Добавить событие в текущий спан
     */
    public static void addEvent(String name) {
        Span.current().addEvent(name);
    }
}