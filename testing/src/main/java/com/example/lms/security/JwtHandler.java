package com.example.lms.security;

import io.javalin.http.Handler;
import io.javalin.http.UnauthorizedResponse;
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
import java.util.concurrent.ConcurrentHashMap;

/**
 * Обработчик JWT-токенов для интеграции с Keycloak.
 * <p>
 * Выполняет:
 * <ul>
 *     <li>извлечение и валидацию JWT из заголовка {@code Authorization}</li>
 *     <li>получение и кэширование публичных RSA-ключей из JWKS Keycloak</li>
 *     <li>проверку подписи и издателя токена</li>
 *     <li>прокладку пользовательских атрибутов (userId, username, email) в контекст Javalin</li>
 * </ul>
 *
 * Предназначен для использования как middleware:
 * <pre>
 * app.before("/api/*", JwtHandler.authenticate());
 * </pre>
 */
public class JwtHandler {

    private static final Logger log = LoggerFactory.getLogger(JwtHandler.class);

    /**
     * Кэш публичных ключей по {@code kid} из JWKS.
     * Используется для уменьшения количества запросов к Keycloak.
     */
    private static final Map<String, CachedKey> publicKeys = new ConcurrentHashMap<>();

    /**
     * Объект-блокировка для синхронизации обновления кэша ключей.
     */
    private static final Object lock = new Object();

    /**
     * Внутренний URL Keycloak (используется из Docker-сети).
     */
    private static final String KEYCLOAK_INTERNAL_URL = System.getenv().getOrDefault("KEYCLOAK_INTERNAL_URL",
            "http://keycloak:8080");

    /**
     * Внешний URL Keycloak (для проверки {@code iss} в токене).
     */
    private static final String KEYCLOAK_URL = System.getenv().getOrDefault("KEYCLOAK_URL", "http://localhost:8080");

    /**
     * Realm Keycloak, в котором выдан токен.
     */
    private static final String REALM = System.getenv().getOrDefault("KEYCLOAK_REALM", "student");

    /**
     * Объект для парсинга JSON-ответов от JWKS эндпоинта.
     */
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

    /**
     * Время жизни кэша публичных ключей.
     * После истечения ключи будут заново загружены из Keycloak.
     */
    private static final Duration CACHE_TTL = Duration.ofHours(1);

    /**
     * Обёртка для кэшируемого публичного ключа и времени его загрузки.
     *
     * @param key      RSA публичный ключ
     * @param loadedAt время загрузки ключа
     */
    private record CachedKey(RSAPublicKey key, Instant loadedAt) {
    }

    /**
     * Гарантирует наличие актуального публичного ключа для указанного {@code kid}.
     * <p>
     * Если ключ отсутствует или протух по TTL, выполняется запрос JWKS к Keycloak
     * и кэш обновляется. В случае ошибки при обновлении:
     * <ul>
     *     <li>если старый ключ есть — продолжаем использовать его</li>
     *     <li>если ключа не было — пробрасывается исключение</li>
     * </ul>
     *
     * @param kid идентификатор ключа из заголовка токена
     * @throws Exception если не удалось загрузить ключи и нет старого значения
     */
    private static void ensurePublicKeyLoaded(String kid) throws Exception {
        CachedKey cached = publicKeys.get(kid);
        if (cached == null || cached.loadedAt.plus(CACHE_TTL).isBefore(Instant.now())) {
            synchronized (lock) {
                cached = publicKeys.get(kid);
                if (cached == null || cached.loadedAt.plus(CACHE_TTL).isBefore(Instant.now())) {
                    try {
                        loadPublicKeys();
                        log.info("JWKS успешно загружены и кэш обновлён");
                    } catch (Exception e) {
                        if (cached != null) {
                            log.warn("Не удалось обновить JWKS, использую старый ключ: {}", e.getMessage());
                        } else {
                            throw e;
                        }
                    }
                }
            }
        }
    }

    /**
     * Загружает публичные ключи из JWKS эндпоинта Keycloak и обновляет кэш {@link #publicKeys}.
     *
     * @throws Exception если запрос к Keycloak или парсинг ответа завершился ошибкой
     */
    private static void loadPublicKeys() throws Exception {
        HttpClient client = HttpClient.newBuilder()
                .connectTimeout(Duration.ofSeconds(10))
                .build();

        String certsUrl = KEYCLOAK_INTERNAL_URL + "/realms/" + REALM + "/protocol/openid-connect/certs";
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(certsUrl))
                .timeout(Duration.ofSeconds(10))
                .GET()
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        if (response.statusCode() != 200) {
            throw new RuntimeException("Keycloak вернул статус " + response.statusCode());
        }

        JsonNode root = OBJECT_MAPPER.readTree(response.body());
        JsonNode keysNode = root.path("keys");
        if (!keysNode.isArray()) {
            throw new RuntimeException("JWKS не содержит keys");
        }

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

            publicKeys.put(kid, new CachedKey(key, Instant.now()));
        }
    }

    /**
     * Извлекает значение {@code kid} из заголовка JWT
     * без проверки подписи.
     *
     * @param token строка JWT-токена
     * @return значение {@code kid} или null, если оно отсутствует
     */
    private static String extractKid(String token) {
        DecodedJWT jwt = JWT.decode(token);
        return jwt.getKeyId();
    }

    /**
     * Возвращает Javalin {@link Handler}, который проверяет JWT-токен
     * в заголовке {@code Authorization} и, в случае успеха:
     * <ul>
     *     <li>проверяет подпись и издателя токена</li>
     *     <li>кладёт в контекст:
     *         <ul>
     *             <li>{@code userId} – subject токена</li>
     *             <li>{@code username} – claim {@code preferred_username}</li>
     *             <li>{@code email} – claim {@code email}</li>
     *             <li>{@code jwt} – полностью декодированный токен</li>
     *         </ul>
     *     </li>
     * </ul>
     * В случае любой ошибки кидает {@link UnauthorizedResponse}.
     *
     * @return Javalin Handler для аутентификации запросов по JWT
     */
    public static Handler authenticate() {
        return ctx -> {
            String authHeader = ctx.header("Authorization");

            if (authHeader == null || !authHeader.startsWith("Bearer ")) {
                throw new UnauthorizedResponse("Отсутствует или неверный заголовок Authorization");
            }

            String token = authHeader.substring(7);
            String kid = extractKid(token);

            log.debug("KID токена: {}", kid);

            try {
                ensurePublicKeyLoaded(kid);
            } catch (Exception e) {
                throw new UnauthorizedResponse("Сервер не готов: " + e.getMessage());
            }

            CachedKey cached = publicKeys.get(kid);
            if (cached == null) {
                throw new UnauthorizedResponse("JWT key not found for kid: " + kid);
            }

            try {
                Algorithm algorithm = Algorithm.RSA256(cached.key(), null);
                DecodedJWT jwt = JWT.require(algorithm)
                        .withIssuer(KEYCLOAK_URL + "/realms/" + REALM)
                        .build()
                        .verify(token);

                ctx.attribute("userId", jwt.getSubject());
                ctx.attribute("username", jwt.getClaim("preferred_username").asString());
                ctx.attribute("email", jwt.getClaim("email").asString());
                ctx.attribute("jwt", jwt);

            } catch (JWTVerificationException e) {
                throw new UnauthorizedResponse("Неверный JWT токен: " + e.getMessage());
            }
        };
    }
}
