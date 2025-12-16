package com.example.lms.config;

import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.trace.Tracer;

/**
 * Конфигурационный класс для работы с OpenTelemetry.
 * <p>
 * Поскольку сервис запускается с OpenTelemetry Java Agent,
 * ручная инициализация SDK внутри приложения не требуется:
 * агент автоматически настраивает экспортеры, ресурс (service.name),
 * обработчики и сборщики метрик/трейсов.
 *
 * Этот класс выполняет две задачи:
 * <ul>
 * <li>фиксирует факт инициализации (для логов и читаемости)</li>
 * <li>предоставляет доступ к глобальному {@link Tracer}</li>
 * </ul>
 *
 * Tracer используется внутри приложения для создания пользовательских
 * (manual) span-ов — например, в SimpleTracer, сервисах и роутерах.
 */
public class OpenTelemetryConfig {

	/**
	 * Метод вызывается при старте сервиса.
	 * <p>
	 * Инициализация трейсинга полностью выполняется Java Agent,
	 * поэтому метод не содержит дополнительной логики.
	 * Оставлен для читаемости архитектуры и удобства дебага.
	 */
	public static void init() {
		System.out.println("OpenTelemetry will be initialized by Java Agent");
	}

	/**
	 * Возвращает глобальный {@link Tracer}, предоставляемый OpenTelemetry Java
	 * Agent.
	 * <p>
	 * Используется для ручного создания span-ов:
	 * 
	 * <pre>
	 * Tracer tracer = OpenTelemetryConfig.getTracer();
	 * Span span = tracer.spanBuilder("my-operation").startSpan();
	 * </pre>
	 *
	 * @return Tracer, связанный с сервисом "testing-service"
	 */
	public static Tracer getTracer() {
		return GlobalOpenTelemetry.getTracer("testing-service");
	}
}