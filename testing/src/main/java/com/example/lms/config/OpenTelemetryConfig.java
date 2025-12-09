// OpenTelemetryConfig.java - ЗАКОММЕНТИРУЙТЕ или УДАЛИТЕ
package com.example.lms.config;

import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.trace.Tracer;

public class OpenTelemetryConfig {
    
    public static void init() {
        // НИЧЕГО НЕ ДЕЛАЕМ! Java Agent сам все настроит
        System.out.println("OpenTelemetry will be initialized by Java Agent");
    }
    
    public static Tracer getTracer() {
        // Используем глобальный tracer от Java Agent
        return GlobalOpenTelemetry.getTracer("testing-service");
    }
}