package com.example.lms.config;

import io.opentelemetry.api.GlobalOpenTelemetry;
import io.opentelemetry.api.trace.Tracer;

public class OpenTelemetryConfig {

    public static void init() {
        System.out.println("OpenTelemetry will be initialized by Java Agent");
    }

    public static Tracer getTracer() {
        return GlobalOpenTelemetry.getTracer("testing-service");
    }
}