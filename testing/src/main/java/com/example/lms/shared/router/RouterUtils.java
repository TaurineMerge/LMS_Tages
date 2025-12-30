package com.example.lms.shared.router;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.security.JwtHandler;
import com.example.lms.shared.validation.JsonSchemaValidator;
import com.example.lms.tracing.SimpleTracer;

import io.github.cdimascio.dotenv.Dotenv;
import io.javalin.http.Handler;

import java.util.List;
import java.util.Set;

/**
 * Утилиты для роутеров приложения.
 * <p>
 * Предоставляет переиспользуемые компоненты для настройки маршрутов:
 * <ul>
 * <li>Общие middleware для логирования и аутентификации</li>
 * <li>Вспомогательные методы для авторизации по realm</li>
 * <li>Методы для композиции валидации и авторизации</li>
 * <li>Константы realm из конфигурации</li>
 * </ul>
 * 
 * <h2>Пример использования:</h2>
 * 
 * <pre>
 * public class UserRouter {
 * 	public static void register(UserController controller) {
 * 		path("/users", () -> {
 * 			applyStandardBeforeMiddleware();
 * 
 * 			// Простая авторизация
 * 			get(withRealm(TEACHER_REALM, controller::getUsers));
 * 
 * 			// Валидация + авторизация
 * 			post(withValidationAndRealm(
 * 					"/schemas/user-schema.json",
 * 					TEACHER_REALM,
 * 					controller::createUser));
 * 
 * 			applyStandardAfterMiddleware();
 * 		});
 * 	}
 * }
 * </pre>
 * 
 * @see JwtHandler
 * @see JsonSchemaValidator
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

	/** Набор realm с доступом только на чтение (студенты и преподаватели) */
	public static final Set<String> READ_ACCESS_REALMS = Set.of(STUDENT_REALM, TEACHER_REALM);

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
	 * delete(withRealm(TEACHER_REALM, controller::deleteUser));
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

	// ============================================================
	// ВАЛИДАЦИЯ + АВТОРИЗАЦИЯ
	// ============================================================

	/**
	 * Создает композитный Handler с валидацией JSON и проверкой realm.
	 * <p>
	 * Последовательность выполнения:
	 * <ol>
	 * <li>Валидация JSON по схеме</li>
	 * <li>Проверка realm</li>
	 * <li>Выполнение бизнес-логики</li>
	 * </ol>
	 * 
	 * <p>
	 * <b>Пример использования:</b>
	 * 
	 * <pre>
	 * post(withValidationAndRealm(
	 * 		"/schemas/answer-schema.json",
	 * 		TEACHER_REALM,
	 * 		controller::createAnswer));
	 * </pre>
	 * 
	 * @param schemaPath путь к JSON Schema в classpath
	 * @param realm      требуемый realm
	 * @param handler    обработчик бизнес-логики
	 * @return композитный Handler
	 * @throws io.javalin.http.BadRequestResponse если валидация не пройдена
	 * @throws io.javalin.http.ForbiddenResponse  если realm не совпадает
	 */
	public static Handler withValidationAndRealm(String schemaPath, String realm, Handler handler) {
		return ctx -> {
			JsonSchemaValidator.validate(schemaPath).handle(ctx);
			JwtHandler.requireRealm(realm).handle(ctx);
			handler.handle(ctx);
		};
	}

	/**
	 * Создает композитный Handler с валидацией JSON и проверкой набора realm.
	 * 
	 * <p>
	 * <b>Пример использования:</b>
	 * 
	 * <pre>
	 * put(withValidationAndRealm(
	 * 		"/schemas/answer-schema.json",
	 * 		READ_ACCESS_REALMS,
	 * 		controller::updateAnswer));
	 * </pre>
	 * 
	 * @param schemaPath путь к JSON Schema в classpath
	 * @param realms     набор требуемых realm
	 * @param handler    обработчик бизнес-логики
	 * @return композитный Handler
	 */
	public static Handler withValidationAndRealm(String schemaPath, Set<String> realms, Handler handler) {
		return ctx -> {
			JsonSchemaValidator.validate(schemaPath).handle(ctx);
			JwtHandler.requireRealm(realms).handle(ctx);
			handler.handle(ctx);
		};
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
	 * 	RouterUtils.applyStandardBeforeMiddleware();
	 * 	// ... маршруты
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
	 * Применяет стандартный набор before() middleware с кастомным логгером.
	 * 
	 * @param customLogger логгер для использования вместо стандартного
	 */
	public static void applyStandardBeforeMiddleware(Logger customLogger, List<String> publicPaths) {
		io.javalin.apibuilder.ApiBuilder.before(ctx -> {
			String path = ctx.path();
			
			for (String pattern : publicPaths) {
				if (path.matches(pattern)) {
					return;
				}
			}
			
			// Для остальных маршрутов применяем аутентификацию
			authenticate().handle(ctx);
		});
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
	 * 	// ... маршруты
	 * 	RouterUtils.applyStandardAfterMiddleware();
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
	// ВАЛИДАЦИЯ ПАРАМЕТРОВ
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
