package com.example.lms;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.config.HandlebarsConfig;
import com.example.lms.question.domain.service.QuestionService;
import com.example.lms.question.infrastructure.repositories.QuestionRepository;
import com.example.lms.question.api.controller.QuestionController;
import com.example.lms.question.api.router.QuestionRouter;
import com.example.lms.answer.api.controller.AnswerController;
import com.example.lms.answer.api.router.AnswerRouter;
import com.example.lms.answer.domain.service.AnswerService;
import com.example.lms.answer.infrastructure.repositories.AnswerRepository;
import com.example.lms.config.DatabaseConfig;
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
        QuestionRepository questionRepository = new QuestionRepository(dbConfig);

        // Сервисный слой (бизнес-логика)
        TestService testService = new TestService(testRepository);
        AnswerService answerService = new AnswerService(answerRepository);
        QuestionService questionService = new QuestionService(questionRepository);

        // Контроллер, принимающий HTTP-запросы
        TestController testController = new TestController(testService, handlebars);
        AnswerController answerController = new AnswerController(answerService);
        QuestionController questionController = new QuestionController(questionService);

        // ---------------------------------------------------------------
        // 3. Создание и запуск Javalin HTTP-сервера
        // ---------------------------------------------------------------
        Javalin app = Javalin.create(config -> {
            // Регистрация маршрутов через ApiBuilder
            config.router.apiBuilder(() -> {
                TestRouter.register(testController);
                AnswerRouter.register(answerController);
                QuestionRouter.register(questionController);
            });
        }).start("0.0.0.0", APP_PORT);

        // Приложение успешно запустилось
        logger.info("Testing Service started on port %d%n", APP_PORT);
    }
}