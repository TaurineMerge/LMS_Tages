package com.example.lms.shared.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.security.JwtHandler;
import com.example.lms.tracing.SimpleTracer;

import io.github.cdimascio.dotenv.Dotenv;
import io.javalin.http.Handler;

import java.util.Set;

/**
 * Утилиты для роутеров приложения.
 * <p>
 * Предоставляет переиспользуемые компоненты для настройки маршрутов:
 * <ul>
 * <li>Общие middleware для логирования и аутентификации</li>
 * <li>Вспомогательные методы для авторизации по realm</li>
 * <li>Константы realm из конфигурации</li>
 * </ul>
 * 
 * <h2>Пример использования:</h2>
 * 
 * <pre>
 * public class UserRouter {
 *     public static void register(UserController controller) {
 *         path("/users", () -> {
 *             // Применить стандартные middleware
 *             before(RouterUtils.authenticate());
 *             before(RouterUtils.logRequestStart());
 * 
 *             // Эндпоинты с авторизацией по realm
 *             get(RouterUtils.withRealm(RouterUtils.ADMIN_REALM, controller::getUsers));
 *             post(RouterUtils.withRealm(RouterUtils.TEACHER_REALM, controller::createUser));
 * 
 *             // Логирование завершения
 *             after(RouterUtils.logRequestEnd());
 *         });
 *     }
 * }
 * </pre>
 * 
 * @see JwtHandler
 * @see SimpleTracer
 */
public class RouterUtils {

    private static final Logger logger = LoggerFactory.getLogger(RouterUtils.class);
    private static final Dotenv dotenv = Dotenv.load();

    // ============================================================
    // КОНСТАНТЫ REALM
    // ============================================================

    /** Realm для студентов из переменных окружения */
    public static final String STUDENT_REALM = dotenv.get("KEYCLOAK_STUDENT_REALM");

    /** Realm для преподавателей из переменных окружения */
    public static final String TEACHER_REALM = dotenv.get("KEYCLOAK_TEACHER_REALM");

    /** Realm для администраторов из переменных окружения (если есть) */
    public static final String ADMIN_REALM = dotenv.get("KEYCLOAK_ADMIN_REALM", "admin-realm");

    /** Набор realm с доступом только на чтение (студенты и преподаватели) */
    public static final Set<String> READ_ACCESS_REALMS = Set.of(STUDENT_REALM, TEACHER_REALM);

    /** Набор realm с доступом на запись (преподаватели и администраторы) */
    public static final Set<String> WRITE_ACCESS_REALMS = Set.of(TEACHER_REALM, ADMIN_REALM);

    /** Набор всех realm системы */
    public static final Set<String> ALL_REALMS = Set.of(STUDENT_REALM, TEACHER_REALM, ADMIN_REALM);

    // ============================================================
    // MIDDLEWARE: Аутентификация
    // ============================================================

    /**
     * Возвращает Handler для JWT аутентификации.
     * <p>
     * Проверяет наличие и валидность JWT токена, извлекает информацию
     * о пользователе (userId, username, email, realm) и сохраняет в контексте.
     * 
     * @return Handler для аутентификации
     * @throws io.javalin.http.UnauthorizedResponse если токен отсутствует или
     *                                              невалиден
     * 
     * @see JwtHandler#authenticate()
     */
    public static Handler authenticate() {
        return JwtHandler.authenticate();
    }

    // ============================================================
    // MIDDLEWARE: Логирование
    // ============================================================

    /**
     * Возвращает Handler для логирования начала обработки запроса.
     * <p>
     * Логирует HTTP метод, путь, traceId и realm пользователя.
     * Используется как before() middleware.
     * 
     * @return Handler для логирования начала запроса
     */
    public static Handler logRequestStart() {
        return ctx -> {
            logger.info("▶️  Request started: {} {} (traceId: {}, realm: {})",
                    ctx.method(),
                    ctx.path(),
                    SimpleTracer.getCurrentTraceId(),
                    ctx.attribute("realm"));
        };
    }

    /**
     * Возвращает Handler для логирования завершения обработки запроса.
     * <p>
     * Логирует HTTP метод, путь, статус ответа, traceId и realm пользователя.
     * Используется как after() middleware.
     * 
     * @return Handler для логирования завершения запроса
     */
    public static Handler logRequestEnd() {
        return ctx -> {
            logger.info("✅ Request completed: {} {} -> {} (traceId: {}, realm: {})",
                    ctx.method(),
                    ctx.path(),
                    ctx.status(),
                    SimpleTracer.getCurrentTraceId(),
                    ctx.attribute("realm"));
        };
    }

    /**
     * Возвращает Handler для логирования начала запроса с кастомным логгером.
     * <p>
     * Полезно когда нужно использовать специфичный логгер роутера.
     * 
     * @param customLogger логгер для вывода сообщений
     * @return Handler для логирования начала запроса
     */
    public static Handler logRequestStart(Logger customLogger) {
        return ctx -> {
            customLogger.info("▶️  Request started: {} {} (traceId: {}, realm: {})",
                    ctx.method(),
                    ctx.path(),
                    SimpleTracer.getCurrentTraceId(),
                    ctx.attribute("realm"));
        };
    }

    /**
     * Возвращает Handler для логирования завершения запроса с кастомным логгером.
     * 
     * @param customLogger логгер для вывода сообщений
     * @return Handler для логирования завершения запроса
     */
    public static Handler logRequestEnd(Logger customLogger) {
        return ctx -> {
            customLogger.info("✅ Request completed: {} {} -> {} (traceId: {}, realm: {})",
                    ctx.method(),
                    ctx.path(),
                    ctx.status(),
                    SimpleTracer.getCurrentTraceId(),
                    ctx.attribute("realm"));
        };
    }

    // ============================================================
    // АВТОРИЗАЦИЯ: Проверка realm
    // ============================================================

    /**
     * Создает композитный Handler с проверкой одного разрешенного realm.
     * <p>
     * Сначала проверяет, что realm пользователя совпадает с требуемым,
     * затем выполняет основную бизнес-логику обработчика.
     * 
     * <p>
     * <b>Пример использования:</b>
     * 
     * <pre>
     * delete(withRealm(ADMIN_REALM, controller::deleteUser));
     * </pre>
     * 
     * @param realm   разрешенный realm для доступа к эндпоинту
     * @param handler обработчик контроллера для выполнения бизнес-логики
     * @return композитный Handler с проверкой realm
     * @throws io.javalin.http.ForbiddenResponse если realm пользователя не
     *                                           совпадает
     */
    public static Handler withRealm(String realm, Handler handler) {
        return ctx -> {
            JwtHandler.requireRealm(realm).handle(ctx);
            handler.handle(ctx);
        };
    }

    /**
     * Создает композитный Handler с проверкой набора разрешенных realm.
     * <p>
     * Сначала проверяет, что realm пользователя входит в список разрешенных,
     * затем выполняет основную бизнес-логику обработчика.
     * 
     * <p>
     * <b>Пример использования:</b>
     * 
     * <pre>
     * get(withRealm(READ_ACCESS_REALMS, controller::getUsers));
     * get(withRealm(Set.of(STUDENT_REALM, TEACHER_REALM), controller::getData));
     * </pre>
     * 
     * @param realms  набор разрешенных realm для доступа к эндпоинту
     * @param handler обработчик контроллера для выполнения бизнес-логики
     * @return композитный Handler с проверкой realm
     * @throws io.javalin.http.ForbiddenResponse если realm пользователя не входит в
     *                                           список
     */
    public static Handler withRealm(Set<String> realms, Handler handler) {
        return ctx -> {
            JwtHandler.requireRealm(realms).handle(ctx);
            handler.handle(ctx);
        };
    }

    /**
     * Создает Handler только с проверкой realm, без дополнительной логики.
     * <p>
     * Полезно когда нужно просто запретить доступ определенным realm,
     * а основную логику обрабатывает другой обработчик.
     * 
     * <p>
     * <b>Пример использования:</b>
     * 
     * <pre>
     * before(requireRealm(ADMIN_REALM)); // Применить ко всей группе маршрутов
     * get(controller::getData); // Логика без дополнительных проверок
     * </pre>
     * 
     * @param realm требуемый realm
     * @return Handler для проверки realm
     */
    public static Handler requireRealm(String realm) {
        return JwtHandler.requireRealm(realm);
    }

    /**
     * Создает Handler только с проверкой набора realm, без дополнительной логики.
     * 
     * @param realms набор требуемых realm
     * @return Handler для проверки realm
     */
    public static Handler requireRealm(Set<String> realms) {
        return JwtHandler.requireRealm(realms);
    }

    // ============================================================
    // КОМБИНИРОВАННЫЕ MIDDLEWARE
    // ============================================================

    /**
     * Применяет стандартный набор before() middleware для роутера.
     * <p>
     * Включает:
     * <ul>
     * <li>JWT аутентификацию</li>
     * <li>Логирование начала запроса</li>
     * </ul>
     * 
     * <p>
     * <b>Пример использования:</b>
     * 
     * <pre>
     * path("/users", () -> {
     *     RouterUtils.applyStandardBeforeMiddleware();
     *     // ... маршруты
     * });
     * </pre>
     * 
     * <p>
     * <i>Примечание:</i> Этот метод должен вызываться внутри path() блока,
     * так как использует ApiBuilder методы.
     */
    public static void applyStandardBeforeMiddleware() {
        io.javalin.apibuilder.ApiBuilder.before(authenticate());
        io.javalin.apibuilder.ApiBuilder.before(logRequestStart());
    }

    /**
     * Применяет стандартный набор before() middleware с кастомным логгером.
     * 
     * @param customLogger логгер для использования вместо стандартного
     */
    public static void applyStandardBeforeMiddleware(Logger customLogger) {
        io.javalin.apibuilder.ApiBuilder.before(authenticate());
        io.javalin.apibuilder.ApiBuilder.before(logRequestStart(customLogger));
    }

    /**
     * Применяет стандартный набор after() middleware для роутера.
     * <p>
     * Включает логирование завершения запроса.
     * 
     * <p>
     * <b>Пример использования:</b>
     * 
     * <pre>
     * path("/users", () -> {
     *     // ... маршруты
     *     RouterUtils.applyStandardAfterMiddleware();
     * });
     * </pre>
     */
    public static void applyStandardAfterMiddleware() {
        io.javalin.apibuilder.ApiBuilder.after(logRequestEnd());
    }

    /**
     * Применяет стандартный набор after() middleware с кастомным логгером.
     * 
     * @param customLogger логгер для использования вместо стандартного
     */
    public static void applyStandardAfterMiddleware(Logger customLogger) {
        io.javalin.apibuilder.ApiBuilder.after(logRequestEnd(customLogger));
    }

    // ============================================================
    // ВАЛИДАЦИЯ
    // ============================================================

    /**
     * Проверяет, что контроллер не равен null.
     * <p>
     * Используется в начале метода register() роутеров для валидации параметров.
     * 
     * @param controller     объект контроллера для проверки
     * @param controllerName имя контроллера для сообщения об ошибке
     * @throws IllegalArgumentException если controller равен null
     */
    public static void validateController(Object controller, String controllerName) {
        if (controller == null) {
            throw new IllegalArgumentException(controllerName + " не может быть null");
        }
    }
}