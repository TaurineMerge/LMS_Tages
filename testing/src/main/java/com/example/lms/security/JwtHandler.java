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
 * Обработчик JWT токенов для аутентификации и авторизации.
 * Поддерживает загрузку публичных ключей из Keycloak и верификацию токенов.
 * Поддерживает работу с несколькими realm и авторизацию на основе realm.
 * Кэширует публичные ключи для улучшения производительности.
 */
public class JwtHandler {

    private static final Logger log = LoggerFactory.getLogger(JwtHandler.class);
    private static final Dotenv dotenv = Dotenv.load();

    /** Кэш публичных ключей по составному ключу "realm:kid" */
    private static final Map<String, CachedKey> publicKeys = new ConcurrentHashMap<>();

    /** Объект для синхронизации при загрузке ключей */
    private static final Object lock = new Object();

    /** Внутренний URL Keycloak для доступа из контейнерной сети */
    private static final String KEYCLOAK_INTERNAL_URL = dotenv.get("KEYCLOAK_INTERNAL_URL");

    /** Публичный URL Keycloak для проверки issuer токена */
    private static final String KEYCLOAK_URL = dotenv.get("KEYCLOAK_EXTERNAL_URL");

    /** Объект для парсинга JSON */
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

    /** Время жизни кэша публичных ключей */
    private static final Duration CACHE_TTL = Duration.ofHours(1);

    /** Паттерн для извлечения realm из issuer URL */
    private static final Pattern REALM_PATTERN = Pattern.compile("/realms/([^/]+)");

    /**
     * Запись для хранения публичного ключа и времени его загрузки.
     * 
     * @param key      публичный RSA ключ
     * @param loadedAt время загрузки ключа
     */
    private record CachedKey(RSAPublicKey key, Instant loadedAt) {
    }

    /**
     * Извлекает realm из issuer claim токена.
     *
     * @param token JWT токен в формате строки
     * @return название realm
     * @throws UnauthorizedResponse если не удалось извлечь realm
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
     * Гарантирует, что публичный ключ для указанного realm и kid загружен и
     * актуален.
     * Если ключ отсутствует или устарел, происходит загрузка из Keycloak.
     *
     * @param realm название realm
     * @param kid   идентификатор ключа (Key ID) из JWT токена
     * @throws Exception если произошла ошибка при загрузке ключей
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
     * Загружает публичные ключи из JWKS эндпоинта Keycloak для указанного realm.
     * Парсит ответ и сохраняет RSA ключи в кэш.
     *
     * @param realm название realm
     * @throws Exception если произошла ошибка при загрузке или парсинге ключей
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
            if (!"RSA".equals(keyNode.path("kty").asText()))
                continue;

            String kid = keyNode.path("kid").asText(null);
            String nValue = keyNode.path("n").asText(null);
            String eValue = keyNode.path("e").asText(null);

            if (kid == null || nValue == null || eValue == null)
                continue;

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
     * Извлекает идентификатор ключа (kid) из JWT токена без верификации.
     *
     * @param token JWT токен в формате строки
     * @return идентификатор ключа (kid)
     */
    private static String extractKid(String token) {
        DecodedJWT jwt = JWT.decode(token);
        return jwt.getKeyId();
    }

    /**
     * Создает обработчик (Handler) для аутентификации JWT токенов в Javalin.
     * Проверяет наличие и валидность токена, извлекает claims и сохраняет их в
     * контексте.
     * Поддерживает работу с несколькими realm - realm извлекается из issuer токена.
     *
     * @return Handler для использования в маршрутах Javalin
     * 
     * @throws UnauthorizedResponse если:
     *                              - отсутствует заголовок Authorization
     *                              - некорректный формат заголовка Authorization
     *                              - не удалось извлечь realm из токена
     *                              - не удалось загрузить публичные ключи
     *                              - не найден ключ для указанного realm и kid
     *                              - токен не прошел верификацию (истек, неверная
     *                              подпись и т.д.)
     * 
     *                              Устанавливает следующие атрибуты в контексте при
     *                              успешной аутентификации:
     *                              - userId: subject токена
     *                              - username: preferred_username из claims
     *                              - email: email из claims
     *                              - realm: название realm из issuer
     *                              - jwt: полный объект DecodedJWT
     */
    public static Handler authenticate() {
        return ctx -> {
            String authHeader = ctx.header("Authorization");

            if (authHeader == null || !authHeader.startsWith("Bearer ")) {
                throw new UnauthorizedResponse("Отсутствует или неверный заголовок Authorization");
            }

            String token = authHeader.substring(7);

            // Извлекаем realm и kid из токена
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
                        .withIssuer(KEYCLOAK_URL + "/realms/" + realm)
                        .build()
                        .verify(token);

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
     * Создает обработчик для авторизации на основе realm.
     * Проверяет, что realm пользователя входит в список разрешенных realm.
     * Должен использоваться ПОСЛЕ authenticate().
     *
     * @param allowedRealms набор разрешенных realm
     * @return Handler для использования в маршрутах Javalin
     * 
     * @throws ForbiddenResponse если realm пользователя не входит в список
     *                           разрешенных
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
     * Создает обработчик для авторизации на основе одного realm.
     * Удобный метод для случаев, когда нужен только один конкретный realm.
     *
     * @param allowedRealm разрешенный realm
     * @return Handler для использования в маршрутах Javalin
     * 
     *         Пример использования:
     * 
     *         <pre>
     * app.get("/admin/settings", 
     *     JwtHandler.authenticate(),
     *     JwtHandler.requireRealm("admin-realm"),
     *     ctx -> { ... }
     * );
     *         </pre>
     */
    public static Handler requireRealm(String allowedRealm) {
        return requireRealm(Set.of(allowedRealm));
    }
}