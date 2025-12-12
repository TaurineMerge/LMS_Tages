package com.example.lms.tracing;

import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.context.Scope;

/**
 * Утилитарный класс для упрощённой работы с OpenTelemetry-трейсингом.
 * <p>
 * Позволяет:
 * <ul>
 *     <li>создавать кастомные спаны вручную</li>
 *     <li>добавлять атрибуты к текущему спану</li>
 *     <li>добавлять события (events)</li>
 *     <li>получать текущий traceId для логгирования</li>
 * </ul>
 *
 * Класс предназначен для использования в сервисах, контроллерах, репозиториях,
 * когда требуется вручную пометить участки кода или записать дополнительные
 * метаданные в трассировку.
 *
 * <p><b>Внимание:</b> OpenTelemetry Java Agent автоматически создаёт спаны
 * для HTTP-запросов, JDBC и других интеграций. Данный класс используется только
 * для кастомных областей, которые нужно отследить вручную.
 */
public class SimpleTracer {

    /** Глобальный OpenTelemetry tracer, полученный из Java Agent. */
    private static final Tracer TRACER =
            GlobalOpenTelemetry.getTracer("testing-service");

    /**
     * Возвращает traceId текущего активного спана.
     * <p>
     * Используется для записи traceId в логи, чтобы связать записи логов
     * с конкретной трассировкой в Jaeger / Zipkin / Tempo.
     *
     * @return строковый traceId или {@code "no-trace-id"}, если спан отсутствует
     */
    public static String getCurrentTraceId() {
        Span current = Span.current();
        if (current != null && current.getSpanContext().isValid()) {
            return current.getSpanContext().getTraceId();
        }
        return "no-trace-id";
    }

    /**
     * Создаёт кастомный спан с заданным именем, выполняет операцию внутри него
     * и корректно закрывает спан после выполнения.
     *
     * <p>Пример:
     * <pre>
     * SimpleTracer.runWithSpan("db.transform.results", () -> {
     *     performComplexTransformation();
     * });
     * </pre>
     *
     * @param spanName имя создаваемого спана
     * @param operation выполняемая логика
     * @throws RuntimeException если операция выбрасывает ошибку — она будет
     * записана в спан, а затем переброшена выше
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
     * Добавляет строковый атрибут в текущий активный спан.
     *
     * @param key ключ атрибута
     * @param value значение
     */
    public static void addAttribute(String key, String value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавляет числовой атрибут (int).
     */
    public static void addAttribute(String key, int value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавляет числовой атрибут (long).
     */
    public static void addAttribute(String key, long value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавляет числовой атрибут (double).
     */
    public static void addAttribute(String key, double value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавляет логический атрибут (boolean).
     */
    public static void addAttribute(String key, boolean value) {
        Span.current().setAttribute(key, value);
    }

    /**
     * Добавляет атрибут на основе объекта: вызывается {@code toString()}.
     *
     * @param key ключ атрибута
     * @param value объект, который будет приведён к строке
     */
    public static void addAttribute(String key, Object value) {
        Span.current().setAttribute(key, value.toString());
    }

    /**
     * Добавляет событие (event) в текущий активный спан.
     * <p>
     * Используется для логирования ключевых точек выполнения:
     * сохранение в БД, вызов внешнего сервиса, обработка ошибок и т.д.
     *
     * @param name название события
     */
    public static void addEvent(String name) {
        Span.current().addEvent(name);
    }
}