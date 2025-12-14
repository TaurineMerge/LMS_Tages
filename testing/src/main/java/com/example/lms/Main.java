package com.example.lms;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.config.HandlebarsConfig;
import com.example.lms.answer.api.controller.AnswerController;
import com.example.lms.answer.api.router.AnswerRouter;
import com.example.lms.answer.domain.service.AnswerService;
import com.example.lms.answer.infrastructure.repositories.AnswerRepository;
import com.example.lms.config.DatabaseConfig;
import com.example.lms.security.JwtHandler;
import com.example.lms.shared.router.RouterUtils;
import com.example.lms.test.api.controller.TestController;
import com.example.lms.test.api.router.TestRouter;
import com.example.lms.test.domain.service.TestService;
import com.example.lms.test.infrastructure.repositories.TestRepository;
import com.github.jknack.handlebars.Handlebars;

import io.github.cdimascio.dotenv.Dotenv;
import io.javalin.Javalin;

/**
 * Главная точка входа в приложение LMS Testing Service.
 * <p>
 * Выполняет:
 * <ul>
 * <li>загрузку переменных окружения через Dotenv</li>
 * <li>инициализацию конфигурации БД</li>
 * <li>ручное внедрение зависимостей (Manual Dependency Injection)</li>
 * <li>создание и настройку Javalin-приложения</li>
 * <li>регистрацию HTTP-маршрутов</li>
 * </ul>
 *
 * Приложение использует Javalin как лёгкий HTTP-фреймворк и следует принципам
 * разделения уровней: Controller → Service → Repository → DB.
 */
public class Main {
	private static final Logger logger = LoggerFactory.getLogger(TestRouter.class);

	/**
	 * Главный метод запуска сервиса.
	 *
	 * @param args аргументы командной строки (не используются)
	 */
	public static void main(String[] args) {
		// ---------------------------------------------------------------
		// 1. Загрузка переменных окружения из файла .env
		// ---------------------------------------------------------------
		Dotenv dotenv = Dotenv.load();

		final Integer APP_PORT = Integer.parseInt(dotenv.get("APP_PORT"));
		final String DB_URL = dotenv.get("DB_URL");
		final String DB_USER = dotenv.get("DB_USER");
		final String DB_PASSWORD = dotenv.get("DB_PASSWORD");

		// ---------------------------------------------------------------
		// 2. Настройка зависимостей (Manual Dependency Injection)
		// ---------------------------------------------------------------
		// Конфигурация Handlebars
		Handlebars handlebars = HandlebarsConfig.configureHandlebars();

		// Конфигурация подключения к базе
		DatabaseConfig dbConfig = new DatabaseConfig(DB_URL, DB_USER, DB_PASSWORD);

		// Репозитории с логикой работы с БД
		TestRepository testRepository = new TestRepository(dbConfig);
		AnswerRepository answerRepository = new AnswerRepository(dbConfig);

		// Сервисный слой (бизнес-логика)
		TestService testService = new TestService(testRepository);
		AnswerService answerService = new AnswerService(answerRepository);

		// Контроллер, принимающий HTTP-запросы
		TestController testController = new TestController(testService, handlebars);
		AnswerController answerController = new AnswerController(answerService);

		// ---------------------------------------------------------------
		// 3. Создание и запуск Javalin HTTP-сервера
		// ---------------------------------------------------------------
		Javalin app = Javalin.create(config -> {
			// Регистрация маршрутов через ApiBuilder
			config.router.apiBuilder(() -> {
				TestRouter.register(testController);
				AnswerRouter.register(answerController);
			});
		}).start("0.0.0.0", APP_PORT);

		// Swagger JSON с аутентификацией
		app.get("/swagger.json", ctx -> {
			try {
				String swaggerJson = new String(
						java.nio.file.Files.readAllBytes(java.nio.file.Paths.get("docs/swagger.json")));
				ctx.contentType("application/json").result(swaggerJson);
			} catch (Exception e) {
				ctx.status(500).result("Error loading swagger.json");
			}
		});

		// Swagger UI страница
		app.get("/swagger", ctx -> {
			String html = """
					<!DOCTYPE html>
					<html>
					<head>
						<title>Testing Service API - Swagger UI</title>
						<link rel="stylesheet" type="text/css" href="/swagger-ui/swagger-ui.css" />
						<link rel="icon" type="image/png" href="/swagger-ui/favicon-32x32.png" sizes="32x32" />
						<link rel="icon" type="image/png" href="/swagger-ui/favicon-16x16.png" sizes="16x16" />
						<style>
							html {
								box-sizing: border-box;
								overflow: -moz-scrollbars-vertical;
								overflow-y: scroll;
							}
							*, *:before, *:after {
								box-sizing: inherit;
							}
							body {
								margin:0;
								background: #fafafa;
							}
						</style>
					</head>
					<body>
						<div id="swagger-ui"></div>
						<script src="/swagger-ui/swagger-ui-bundle.js"></script>
						<script src="/swagger-ui/swagger-ui-standalone-preset.js"></script>
						<script>
							window.onload = function() {
								const ui = SwaggerUIBundle({
									url: '/swagger.json',
									dom_id: '#swagger-ui',
									deepLinking: true,
									presets: [
										SwaggerUIBundle.presets.apis,
										SwaggerUIStandalonePreset
									],
									plugins: [
										SwaggerUIBundle.plugins.DownloadUrl
									],
									layout: "StandaloneLayout"
								});
							};
						</script>
					</body>
					</html>
					""";
			ctx.contentType("text/html").result(html);
		});

		// Статические файлы Swagger UI
		app.get("/swagger-ui/*", ctx -> {
			String path = ctx.path().replace("/swagger-ui/", "");
			try {
				var resourcePath = "/META-INF/resources/webjars/swagger-ui/5.10.3/" + path;
				var inputStream = Main.class.getResourceAsStream(resourcePath);
				if (inputStream != null) {
					byte[] bytes = inputStream.readAllBytes();
					inputStream.close();

					// Определяем content type на основе расширения файла
					if (path.endsWith(".css")) {
						ctx.contentType("text/css");
					} else if (path.endsWith(".js")) {
						ctx.contentType("application/javascript");
					} else if (path.endsWith(".png")) {
						ctx.contentType("image/png");
					} else if (path.endsWith(".json")) {
						ctx.contentType("application/json");
					} else {
						ctx.contentType("text/plain");
					}

					ctx.result(bytes);
				} else {
					ctx.status(404).result("File not found: " + path);
				}
			} catch (Exception e) {
				ctx.status(500).result("Error loading file: " + path);
			}
		});

		// ЗАЩИТА SWAGGER МАРШРУТОВ - добавляем ПОСЛЕ регистрации всех маршрутов!
		app.before("/swagger*", ctx -> {
			try {
				JwtHandler.authenticate().handle(ctx);
				JwtHandler.requireRealm(RouterUtils.READ_ACCESS_REALMS).handle(ctx);
			} catch (Exception e) {
				// Обработка ошибок аутентификации
				ctx.status(401).result("Unauthorized: " + e.getMessage());
			}
		});

		// Приложение успешно запустилось
		logger.info("Testing Service started on port %d%n", APP_PORT);
		logger.info("Protected Swagger UI available at: http://localhost:%d/swagger%n", APP_PORT);
		logger.info("Protected Swagger JSON available at: http://localhost:%d/swagger.json%n", APP_PORT);
	}
}
