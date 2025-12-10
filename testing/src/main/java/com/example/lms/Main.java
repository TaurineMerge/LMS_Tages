package com.example.lms;

import io.javalin.Javalin;
import com.example.lms.test.api.router.TestRouter;
import com.example.lms.test_attempt.api.router.TestAttemptRouter;
import com.example.lms.tracing.SimpleTracer;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Точка входа в сервис тестирования (Testing Service).
 * <p>
 * Отвечает за:
 * <ul>
 *     <li>запуск Javalin-приложения на порту 8085</li>
 *     <li>регистрацию роутеров тестов и попыток тестов</li>
 *     <li>логирование входящих и исходящих запросов</li>
 *     <li>интеграцию с OpenTelemetry/Tracing через {@link SimpleTracer}</li>
 *     <li>предоставление health-check эндпоинта {@code GET /health}</li>
 * </ul>
 *
 * Сервис ожидается, что запускается с OpenTelemetry Java Agent, поэтому
 * в логах дополнительно выводится информация по OTLP и Jaeger UI.
 */
public class Main {

    private static final Logger logger = LoggerFactory.getLogger(Main.class);

    /**
     * Точка входа Java-приложения.
     * <p>
     * Действия метода:
     * <ol>
     *     <li>Выводит в лог информацию о запуске сервиса и настройках трейсинга</li>
     *     <li>Создаёт и настраивает Javalin-приложение</li>
     *     <li>Регистрирует роутеры:
     *         <ul>
     *             <li>{@link TestRouter} — работа с тестами</li>
     *             <li>{@link TestAttemptRouter} — работа с попытками тестов</li>
     *         </ul>
     *     </li>
     *     <li>Добавляет middleware для логирования всех запросов (before/after)</li>
     *     <li>Регистрирует эндпоинт {@code GET /health} для проверки живости сервиса</li>
     * </ol>
     *
     * @param args аргументы командной строки (не используются)
     */
    public static void main(String[] args) {
        logger.info("Starting Testing Service on port 8085");
        logger.info("OpenTelemetry Java Agent ENABLED");
        logger.info("OTLP Endpoint: http://otel-collector:4318");
        logger.info("Jaeger UI will be available at: http://localhost:16686");

        // Создаём и настраиваем Javalин-приложение
        var app = Javalin.create(config -> {
            config.router.apiBuilder(() -> {
                TestRouter.register();
                TestAttemptRouter.register();
            });
        }).start(8085);

        // Middleware ДО обработки запроса
        app.before(ctx -> {
            logger.info("Request started: {} {} (traceId: {})",
                    ctx.method(),
                    ctx.path(),
                    SimpleTracer.getCurrentTraceId());
        });

        // Middleware ПОСЛЕ обработки запроса
        app.after(ctx -> {
            logger.info("Request completed: {} {} -> {} (traceId: {})",
                    ctx.method(),
                    ctx.path(),
                    ctx.status(),
                    SimpleTracer.getCurrentTraceId());
        });

        // Health-check эндпоинт
        app.get("/health", ctx -> {
            ctx.json("{\"status\": \"OK\", \"service\": \"testing\", \"traceId\": \"" +
                    SimpleTracer.getCurrentTraceId() + "\"}");
        });

        logger.info("Testing service successfully started!");
    }
}