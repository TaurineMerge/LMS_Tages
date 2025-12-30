package com.example.lms.security;

import io.github.cdimascio.dotenv.Dotenv;
import io.javalin.http.Handler;
import io.javalin.http.UnauthorizedResponse;
import io.javalin.http.ForbiddenResponse;
import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.interfaces.DecodedJWT;
import com.auth0.jwt.exceptions.JWTVerificationException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.math.BigInteger;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.security.KeyFactory;
import java.security.interfaces.RSAPublicKey;
import java.security.spec.RSAPublicKeySpec;
import java.time.Duration;
import java.time.Instant;
import java.util.Base64;
import java.util.Map;
import java.util.Set;
import java.util.concurrent.ConcurrentHashMap;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Обработчик JWT токенов для аутентификации и авторизации в системе LMS.
 * <p>
 * Класс обеспечивает:
 * <ul>
 *   <li>Верификацию JWT токенов с использованием публичных ключей RSA из Keycloak</li>
 *   <li>Поддержку работы с несколькими realm (областями безопасности)</li>
 *   <li>Кэширование публичных ключей для улучшения производительности</li>
 *   <li>Извлечение и проверку claims (утверждений) из токенов</li>
 *   <li>Авторизацию на основе принадлежности к определенному realm</li>
 * </ul>
 * <p>
 * Класс использует библиотеку Auth0 JWT для работы с токенами и Jackson для парсинга JSON.
 * Публичные ключи загружаются из JWKS (JSON Web Key Set) эндпоинтов Keycloak.
 *
 * @author Команда разработки LMS
 * @version 1.0
 * @see <a href="https://auth0.com/docs/secure/tokens/json-web-tokens">JSON Web Tokens</a>
 * @see <a href="https://www.keycloak.org/docs/latest/securing_apps/">Keycloak Documentation</a>
 */
public class JwtHandler {

    private static final Logger log = LoggerFactory.getLogger(JwtHandler.class);
    private static final Dotenv dotenv = Dotenv.load();

    private static final Map<String, CachedKey> publicKeys = new ConcurrentHashMap<>();
    private static final Object lock = new Object();
    private static final String KEYCLOAK_INTERNAL_URL = dotenv.get("KEYCLOAK_INTERNAL_URL");
    private static final String KEYCLOAK_URL = dotenv.get("KEYCLOAK_EXTERNAL_URL");
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();
    private static final Duration CACHE_TTL = Duration.ofHours(1);
    private static final Pattern REALM_PATTERN = Pattern.compile("/realms[-/]([^/]+)");

    /**
     * Внутренний класс для хранения публичного ключа и времени его загрузки в кэше.
     *
     * @param key публичный RSA ключ для верификации JWT токенов
     * @param loadedAt время загрузки ключа в кэш
     */
    private record CachedKey(RSAPublicKey key, Instant loadedAt) {
    }

    /**
     * Извлекает название realm из issuer claim JWT токена.
     *
     * @param token JWT токен в формате строки
     * @return название realm
     * @throws UnauthorizedResponse если токен не содержит issuer или не удалось извлечь realm
     * @throws IllegalArgumentException если токен имеет неверный формат
     */
    private static String extractRealm(String token) {
        DecodedJWT jwt = JWT.decode(token);
        String issuer = jwt.getIssuer();

        if (issuer == null) {
            throw new UnauthorizedResponse("Токен не содержит issuer");
        }

        Matcher matcher = REALM_PATTERN.matcher(issuer);
        if (!matcher.find()) {
            throw new UnauthorizedResponse("Не удалось извлечь realm из issuer: " + issuer);
        }

        return matcher.group(1);
    }

    /**
     * Гарантирует, что публичный ключ для указанного realm и kid загружен и актуален.
     * <p>
     * Если ключ отсутствует в кэше или устарел (время жизни превышает {@link #CACHE_TTL}),
     * выполняется загрузка ключей из Keycloak для указанного realm.
     *
     * @param realm название realm
     * @param kid идентификатор ключа (Key ID) из JWT токена
     * @throws Exception если произошла ошибка при загрузке ключей и в кэше нет запасного ключа
     */
    private static void ensurePublicKeyLoaded(String realm, String kid) throws Exception {
        String cacheKey = realm + ":" + kid;
        CachedKey cached = publicKeys.get(cacheKey);

        if (cached == null || cached.loadedAt.plus(CACHE_TTL).isBefore(Instant.now())) {
            synchronized (lock) {
                cached = publicKeys.get(cacheKey);
                if (cached == null || cached.loadedAt.plus(CACHE_TTL).isBefore(Instant.now())) {
                    try {
                        loadPublicKeys(realm);
                        log.info("JWKS для realm '{}' успешно загружены и кэш обновлён", realm);
                    } catch (Exception e) {
                        if (cached != null) {
                            log.warn("Не удалось обновить JWKS для realm '{}', использую старый ключ: {}",
                                    realm, e.getMessage());
                        } else {
                            throw e;
                        }
                    }
                }
            }
        }
    }

    /**
     * Загружает публичные ключи из JWKS (JSON Web Key Set) эндпоинта Keycloak для указанного realm.
     * <p>
     * Метод выполняет HTTP-запрос к Keycloak, парсит JSON-ответ и извлекает RSA публичные ключи,
     * которые затем сохраняются в кэше.
     *
     * @param realm название realm, для которого загружаются ключи
     * @throws Exception если произошла ошибка HTTP-запроса, парсинга JSON или создания ключей
     */
    private static void loadPublicKeys(String realm) throws Exception {
        HttpClient client = HttpClient.newBuilder()
                .connectTimeout(Duration.ofSeconds(10))
                .build();

        String certsUrl = KEYCLOAK_INTERNAL_URL + "/realms/" + realm + "/protocol/openid-connect/certs";
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(certsUrl))
                .timeout(Duration.ofSeconds(10))
                .GET()
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        if (response.statusCode() != 200) {
            throw new RuntimeException("Keycloak вернул статус " + response.statusCode() +
                    " для realm '" + realm + "'");
        }

        JsonNode root = OBJECT_MAPPER.readTree(response.body());
        JsonNode keysNode = root.path("keys");
        if (!keysNode.isArray())
            throw new RuntimeException("JWKS не содержит keys для realm '" + realm + "'");

        for (JsonNode keyNode : keysNode) {
            if (!"RSA".equals(keyNode.path("kty").asText())) {
                continue;
            }

            String kid = keyNode.path("kid").asText(null);
            String nValue = keyNode.path("n").asText(null);
            String eValue = keyNode.path("e").asText(null);

            if (kid == null || nValue == null || eValue == null) {
                continue;
            }

            byte[] nBytes = Base64.getUrlDecoder().decode(nValue);
            byte[] eBytes = Base64.getUrlDecoder().decode(eValue);
            BigInteger modulus = new BigInteger(1, nBytes);
            BigInteger exponent = new BigInteger(1, eBytes);

            KeyFactory factory = KeyFactory.getInstance("RSA");
            RSAPublicKey key = (RSAPublicKey) factory.generatePublic(new RSAPublicKeySpec(modulus, exponent));

            String cacheKey = realm + ":" + kid;
            publicKeys.put(cacheKey, new CachedKey(key, Instant.now()));
        }
    }

    /**
     * Извлекает идентификатор ключа (kid) из JWT токена без верификации подписи.
     *
     * @param token JWT токен в формате строки
     * @return идентификатор ключа (kid) или {@code null}, если токен не содержит kid
     * @throws IllegalArgumentException если токен имеет неверный формат
     */
    private static String extractKid(String token) {
        DecodedJWT jwt = JWT.decode(token);
        return jwt.getKeyId();
    }

    /**
     * Создает обработчик (Handler) для аутентификации JWT токенов в Javalin.
     * <p>
     * Обработчик выполняет следующие действия:
     * <ol>
     *   <li>Извлекает токен из заголовка Authorization (формат "Bearer {token}")</li>
     *   <li>Извлекает realm и kid из токена</li>
     *   <li>Загружает и кэширует публичные ключи при необходимости</li>
     *   <li>Верифицирует подпись токена с использованием RSA256</li>
     *   <li>Проверяет issuer (издателя) токена на соответствие ожидаемому</li>
     *   <li>Извлекает claims и сохраняет их в контексте запроса</li>
     * </ol>
     * <p>
     * При успешной аутентификации устанавливаются следующие атрибуты в контексте:
     * <ul>
     *   <li>{@code userId} - subject токена (идентификатор пользователя)</li>
     *   <li>{@code username} - preferred_username из claims</li>
     *   <li>{@code email} - email из claims</li>
     *   <li>{@code realm} - название realm из issuer</li>
     *   <li>{@code jwt} - полный объект DecodedJWT</li>
     * </ul>
     *
     * @return Handler для использования в маршрутах Javalin
     * @throws UnauthorizedResponse если:
     *                              <ul>
     *                                <li>Отсутствует заголовок Authorization</li>
     *                                <li>Некорректный формат заголовка Authorization</li>
     *                                <li>Не удалось извлечь realm из токена</li>
     *                                <li>Не удалось загрузить публичные ключи</li>
     *                                <li>Не найден ключ для указанного realm и kid</li>
     *                                <li>Токен не прошел верификацию (истек, неверная подпись и т.д.)</li>
     *                                <li>Issuer токена не соответствует ожидаемому</li>
     *                              </ul>
     */
    public static Handler authenticate() {
        return ctx -> {
            String authHeader = ctx.header("Authorization");

            if (authHeader == null || !authHeader.startsWith("Bearer ")) {
                throw new UnauthorizedResponse("Отсутствует или неверный заголовок Authorization");
            }

            String token = authHeader.substring(7);

            String realm = extractRealm(token);
            String kid = extractKid(token);

            log.debug("Realm: {}, KID: {}", realm, kid);

            try {
                ensurePublicKeyLoaded(realm, kid);
            } catch (Exception e) {
                throw new UnauthorizedResponse("Сервер не готов: " + e.getMessage());
            }

            String cacheKey = realm + ":" + kid;
            CachedKey cached = publicKeys.get(cacheKey);
            if (cached == null) {
                throw new UnauthorizedResponse("JWT key not found for realm '" + realm + "' and kid: " + kid);
            }

            try {
                Algorithm algorithm = Algorithm.RSA256(cached.key(), null);

                DecodedJWT jwt = JWT.require(algorithm)
                        .build()
                        .verify(token);

                String base = (KEYCLOAK_URL == null ? "" : KEYCLOAK_URL).replaceAll("/+$", "");
                String baseNoAuth = base.endsWith("/auth") ? base.substring(0, base.length() - 5) : base;

                String expectedIssuerNoAuth = baseNoAuth + "/realms/" + realm;
                String expectedIssuerWithAuth = baseNoAuth + "/auth/realms/" + realm;

                String actualIssuer = jwt.getIssuer();
                if (actualIssuer == null ||
                        (!actualIssuer.equals(expectedIssuerNoAuth) && !actualIssuer.equals(expectedIssuerWithAuth))) {
                    throw new UnauthorizedResponse(
                            "Неверный issuer в JWT. Ожидалось: '" + expectedIssuerNoAuth + "' или '" +
                                    expectedIssuerWithAuth + "', получили: '" + actualIssuer + "'");
                }

                ctx.attribute("userId", jwt.getSubject());
                ctx.attribute("username", jwt.getClaim("preferred_username").asString());
                ctx.attribute("email", jwt.getClaim("email").asString());
                ctx.attribute("realm", realm);
                ctx.attribute("jwt", jwt);

            } catch (JWTVerificationException e) {
                throw new UnauthorizedResponse("Неверный JWT токен: " + e.getMessage());
            }
        };
    }

    /**
     * Создает обработчик для авторизации на основе принадлежности к определенному realm.
     * <p>
     * Обработчик проверяет, что realm пользователя (установленный {@link #authenticate()})
     * входит в набор разрешенных realm. Должен использоваться ПОСЛЕ {@link #authenticate()}.
     *
     * @param allowedRealms набор разрешенных realm
     * @return Handler для использования в маршрутах Javalin
     * @throws UnauthorizedResponse если realm не найден в контексте (authenticate() не был вызван)
     * @throws ForbiddenResponse если realm пользователя не входит в список разрешенных
     */
    public static Handler requireRealm(Set<String> allowedRealms) {
        return ctx -> {
            String userRealm = ctx.attribute("realm");

            if (userRealm == null) {
                throw new UnauthorizedResponse(
                        "Realm не найден в контексте. Используйте authenticate() перед requireRealm()");
            }

            if (!allowedRealms.contains(userRealm)) {
                log.warn("Доступ запрещен: realm '{}' не входит в список разрешенных {}",
                        userRealm, allowedRealms);
                throw new ForbiddenResponse("Доступ запрещен для realm: " + userRealm);
            }

            log.debug("Доступ разрешен для realm: {}", userRealm);
        };
    }

    /**
     * Создает обработчик для авторизации на основе одного конкретного realm.
     * <p>
     * Удобный метод для случаев, когда требуется доступ только для одного realm.
     *
     * @param allowedRealm разрешенный realm
     * @return Handler для использования в маршрутах Javalin
     * 
     * @example Пример использования:
     * <pre>{@code
     * app.get("/admin/settings", 
     *     JwtHandler.authenticate(),
     *     JwtHandler.requireRealm("admin-realm"),
     *     ctx -> { ... }
     * );
     * }</pre>
     */
    public static Handler requireRealm(String allowedRealm) {
        return requireRealm(Set.of(allowedRealm));
    }
}