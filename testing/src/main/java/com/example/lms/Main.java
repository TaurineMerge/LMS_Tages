package com.example.lms;

import io.github.cdimascio.dotenv.Dotenv;

import io.javalin.Javalin;
import com.example.lms.test.api.router.TestRouter;
import com.example.lms.test.domain.service.TestService;
import com.example.lms.test.infrastructure.repositories.TestRepository;

import com.example.lms.config.DatabaseConfig;
import com.example.lms.test.api.controller.TestController;

/**
 * Главная точка входа в приложение LMS Testing Service.
 * <p>
 * Выполняет:
 * <ul>
 *     <li>загрузку переменных окружения через Dotenv</li>
 *     <li>инициализацию конфигурации БД</li>
 *     <li>ручное внедрение зависимостей (Manual Dependency Injection)</li>
 *     <li>создание и настройку Javalin-приложения</li>
 *     <li>регистрацию HTTP-маршрутов</li>
 * </ul>
 *
 * Приложение использует Javalin как лёгкий HTTP-фреймворк и следует принципам
 * разделения уровней: Controller → Service → Repository → DB.
 */
public class Main {

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

        // Конфигурация подключения к базе
        DatabaseConfig dbConfig = new DatabaseConfig(DB_URL, DB_USER, DB_PASSWORD);

        // Репозиторий с логикой работы с БД
        TestRepository testRepository = new TestRepository(dbConfig);

        // Сервисный слой (бизнес-логика)
        TestService testService = new TestService(testRepository);

        // Контроллер, принимающий HTTP-запросы
        TestController testController = new TestController(testService);

        // ---------------------------------------------------------------
        // 3. Создание и запуск Javalin HTTP-сервера
        // ---------------------------------------------------------------
        Javalin app = Javalin.create(config -> {
            // Регистрация маршрутов через ApiBuilder
            config.router.apiBuilder(() -> {
                TestRouter.register(testController);
            });
        }).start(APP_PORT);

        // Приложение успешно запустилось
        System.out.printf("Testing Service started on port %d%n", APP_PORT);
    }
}