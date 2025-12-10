package com.example.lms.config;

import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.trace.Tracer;

/**
 * Конфигурационный класс для OpenTelemetry.
 * <p>
 * Поскольку проект использует OpenTelemetry Java Agent,
 * инициализация трейсинга происходит автоматически
 * при запуске приложения с аргументами:
 *
 * <pre>
 * -javaagent:opentelemetry-javaagent.jar
 * -Dotel.service.name=testing-service
 * </pre>
 *
 * Поэтому ручная настройка SDK внутри приложения НЕ требуется.
 * Этот класс остается тонкой оберткой над глобальными объектами OpenTelemetry.
 */
public class OpenTelemetryConfig {

    /**
     * Инициализация OpenTelemetry.
     * <p>
     * Метод умышленно пустой, так как все настройки
     * выполняет OpenTelemetry Java Agent.
     */
    public static void init() {
        // Java Agent полностью берет на себя настройку OpenTelemetry
        System.out.println("OpenTelemetry will be initialized by Java Agent");
    }

    /**
     * Возвращает глобальный Tracer, зарегистрированный Java Agent.
     * <p>
     * Используется для создания span-ов в сервисах и контроллерах.
     *
     * @return глобальный Tracer, связанный с сервисом testing-service
     */
    public static Tracer getTracer() {
        return GlobalOpenTelemetry.getTracer("testing-service");
    }
}