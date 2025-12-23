package com.example.lms;

import java.io.InputStream;
import java.nio.charset.StandardCharsets;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.answer.api.controller.AnswerController;
import com.example.lms.answer.api.router.AnswerRouter;
import com.example.lms.answer.domain.service.AnswerService;
import com.example.lms.answer.infrastructure.repositories.AnswerRepository;
import com.example.lms.config.DatabaseConfig;
import com.example.lms.config.HandlebarsConfig;
import com.example.lms.content.api.controller.ContentController;
import com.example.lms.content.api.router.ContentRouter;
import com.example.lms.content.domain.service.ContentService;
import com.example.lms.content.infrastructure.repositories.ContentRepository;
import com.example.lms.draft.api.controller.DraftController;
import com.example.lms.draft.api.router.DraftRouter;
import com.example.lms.draft.domain.service.DraftService;
import com.example.lms.draft.infrastructure.repositories.DraftRepository;
import com.example.lms.question.api.controller.QuestionController;
import com.example.lms.question.api.router.QuestionRouter;
import com.example.lms.question.domain.service.QuestionService;
import com.example.lms.question.infrastructure.repositories.QuestionRepository;
import com.example.lms.test.api.controller.TestController;
import com.example.lms.test.api.router.TestRouter;
import com.example.lms.test.domain.service.TestService;
import com.example.lms.test.infrastructure.repositories.TestRepository;
import com.example.lms.test.web.controller.TestFormController;
import com.example.lms.test.web.router.TestWebRouter;
import com.example.lms.test_attempt.api.controller.TestAttemptController;
import com.example.lms.test_attempt.api.router.TestAttemptRouter;
import com.example.lms.test_attempt.domain.service.TestAttemptService;
import com.example.lms.test_attempt.infrastructure.repositories.TestAttemptRepository;
import com.example.lms.ui.UiRouter;
import com.example.lms.ui.UiTestController;
import com.github.jknack.handlebars.Handlebars;

import io.github.cdimascio.dotenv.Dotenv;
import io.javalin.Javalin;
import static io.javalin.apibuilder.ApiBuilder.get;
import io.javalin.http.staticfiles.Location;
import io.javalin.json.JavalinJackson;

/**
 * Главная точка входа в приложение LMS Testing Service.
 */
public class Main {
	private static final Logger logger = LoggerFactory.getLogger(Main.class);

	/**
	 * Главный метод запуска сервиса.
	 *
	 * @param args аргументы командной строки
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
		QuestionRepository questionRepository = new QuestionRepository(dbConfig);
		TestAttemptRepository testAttemptRepository = new TestAttemptRepository(dbConfig);
		DraftRepository draftRepository = new DraftRepository(dbConfig);
		ContentRepository contentRepository = new ContentRepository(dbConfig);

		// Сервисный слой (бизнес-логика)
		TestService testService = new TestService(testRepository);
		AnswerService answerService = new AnswerService(answerRepository);
		QuestionService questionService = new QuestionService(questionRepository);
		TestAttemptService testAttemptService = new TestAttemptService(testAttemptRepository);
		DraftService draftService = new DraftService(draftRepository);
		ContentService contentService = new ContentService(contentRepository);

		// Контроллер, принимающий HTTP-запросы
		TestController testController = new TestController(testService, handlebars);
		AnswerController answerController = new AnswerController(answerService);
		QuestionController questionController = new QuestionController(questionService);
		TestAttemptController testAttemptController = new TestAttemptController(testAttemptService);
		DraftController draftController = new DraftController(draftService);
		ContentController contentController = new ContentController(contentService, handlebars);
		
		// ИСПРАВЛЕННАЯ СТРОКА: добавлен draftService в конструктор TestFormController
		TestFormController testWebController = new TestFormController(
			testService, 
			questionService, 
			answerService, 
			draftService,  // Добавлен draftService
			handlebars
		);

		// UI controller (ВАЖНО: добавили testAttemptService)
		var uiTestController = new UiTestController(testService, questionService, answerService, testAttemptService);

		// ---------------------------------------------------------------
		// 3. Создание и запуск Javalin HTTP-сервера
		// ---------------------------------------------------------------
		Javalin app = Javalin.create(config -> {
			// Используем JavalinJackson по умолчанию (без кастомного ObjectMapper)
			config.jsonMapper(new JavalinJackson());
			
			// Добавляем логирование запросов для отладки
			config.requestLogger.http((ctx, executionTimeMs) -> {
				logger.info("{} {} - {}ms", ctx.method(), ctx.path(), executionTimeMs);
			});

			// Настройка статических файлов
			config.staticFiles.add(staticFiles -> {
				staticFiles.hostedPath = "/";
				staticFiles.directory = "/public";
				staticFiles.location = Location.CLASSPATH;
				staticFiles.precompress = false;
			});

			// Для webjars добавляем только если они существуют
			try {
				if (Main.class.getResource("/META-INF/resources/webjars") != null) {
					config.staticFiles.add(staticFiles -> {
						staticFiles.hostedPath = "/webjars";
						staticFiles.directory = "/META-INF/resources/webjars";
						staticFiles.location = Location.CLASSPATH;
					});
				}
			} catch (Exception e) {
				logger.warn("Webjars not found, skipping webjars static files configuration");
			}

			// Регистрация маршрутов
			config.router.apiBuilder(() -> {
				TestRouter.register(testController);
				AnswerRouter.register(answerController);
				// Веб-маршруты конструктора тестов
				TestWebRouter.register(testWebController);
				QuestionRouter.register(questionController);
				TestAttemptRouter.register(testAttemptController);
				DraftRouter.register(draftController);
				ContentRouter.register(contentController);

				// UI маршруты
				UiRouter.register(uiTestController);

				// Swagger UI
				get("/swagger", ctx -> {
					try (InputStream inputStream = Main.class.getClassLoader().getResourceAsStream("public/swagger.html")) {
						if (inputStream == null) {
							logger.error("Swagger HTML file not found in classpath");
							ctx.status(404).result("Swagger UI not found");
							return;
						}
						String html = new String(inputStream.readAllBytes());
						ctx.html(html);
					} catch (Exception e) {
						logger.error("Error loading swagger.html", e);
						ctx.status(500).result("Internal server error");
					}
				});

				// Swagger JSON
				get("/swagger.json", ctx -> {
					try (InputStream inputStream = Main.class.getClassLoader().getResourceAsStream("docs/swagger.json")) {
						if (inputStream == null) {
							logger.error("Swagger JSON file not found in classpath");
							ctx.status(404).result("Swagger JSON not found");
							return;
						}
						String json = new String(inputStream.readAllBytes());
						ctx.contentType("application/json").result(json);
					} catch (Exception e) {
						logger.error("Error loading swagger.json", e);
						ctx.status(500).result("Internal server error");
					}
				});

				get("/docs/swagger.json", ctx -> {
					logger.info("Loading swagger.json from /docs...");

					// Пробуем загрузить из /docs
					try (InputStream inputStream = Main.class.getClassLoader()
							.getResourceAsStream("docs/swagger.json")) {
						if (inputStream != null) {
							String json = new String(inputStream.readAllBytes(), StandardCharsets.UTF_8);
							logger.info("Successfully loaded swagger.json from /docs");
							ctx.contentType("application/json").result(json);
							return;
						}
					} catch (Exception e) {
						logger.warn("Error loading from /docs: {}", e.getMessage());
					}

					// Пробуем загрузить из /public
					try (InputStream inputStream = Main.class.getClassLoader()
							.getResourceAsStream("public/swagger.json")) {
						if (inputStream != null) {
							String json = new String(inputStream.readAllBytes(), StandardCharsets.UTF_8);
							logger.info("Successfully loaded swagger.json from /public");
							ctx.contentType("application/json").result(json);
							return;
						}
					} catch (Exception e) {
						logger.warn("Error loading from /public: {}", e.getMessage());
					}

					// Если файл не найден, возвращаем минимальную спецификацию
					logger.warn("swagger.json not found, returning default spec");
				});

				get("/oauth2-redirect.html", ctx -> {
					try (InputStream inputStream = Main.class.getClassLoader()
							.getResourceAsStream("public/oauth2-redirect.html")) {
						if (inputStream == null) {
							ctx.status(404).result("OAuth2 redirect page not found");
							return;
						}
						String html = new String(inputStream.readAllBytes());
						ctx.html(html);
					}
				});

				// Health check endpoint
				get("/health", ctx -> ctx.result("OK"));
			});
		}).start("0.0.0.0", APP_PORT);

		// Информация о запуске
		logger.info("Testing Service started on port {}", APP_PORT);
		logger.info("Swagger UI available at: http://localhost:{}/swagger", APP_PORT);
		logger.info("Swagger JSON available at: http://localhost:{}/swagger.json", APP_PORT);
		logger.info("Health check: http://localhost:{}/health", APP_PORT);

		// Добавляем обработчик завершения приложения
		Runtime.getRuntime().addShutdownHook(new Thread(() -> {
			logger.info("Shutting down Testing Service...");
			app.stop();
		}));

		// Приложение успешно запустилось
		logger.info("Testing Service started on port %d%n", APP_PORT);
	}
}