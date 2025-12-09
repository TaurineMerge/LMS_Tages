package com.example.lms;

import io.javalin.Javalin;
import com.example.lms.test.api.router.TestRouter;
import com.example.lms.test_attempt.api.router.TestAttemptRouter;
import com.example.lms.tracing.SimpleTracer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

public class Main {
    private static final Logger logger = LoggerFactory.getLogger(Main.class);
    
    public static void main(String[] args) {
        logger.info("🚀 Starting Testing Service on port 8085");
        logger.info("📊 OpenTelemetry Java Agent ENABLED");
        logger.info("🔗 OTLP Endpoint: http://otel-collector:4318");
        logger.info("📈 Jaeger UI will be available at: http://localhost:16686");
        
        var app = Javalin.create(config -> {
            config.router.apiBuilder(() -> {
                TestRouter.register();
                TestAttemptRouter.register();
            });
        }).start(8085);
        
        // Добавляем middleware ДО старта приложения
        app.before(ctx -> {
            logger.info("▶️  Request started: {} {} (traceId: {})", 
                ctx.method(), 
                ctx.path(),
                SimpleTracer.getCurrentTraceId());
        });
        
        app.after(ctx -> {
            logger.info("✅ Request completed: {} {} -> {} (traceId: {})", 
                ctx.method(), 
                ctx.path(),
                ctx.status(),
                SimpleTracer.getCurrentTraceId());
        });
        
        app.get("/health", ctx -> {
            ctx.json("{\"status\": \"OK\", \"service\": \"testing\", \"traceId\": \"" + 
                    SimpleTracer.getCurrentTraceId() + "\"}");
        });
        
        logger.info("✅ Testing service successfully started!");
    }
}