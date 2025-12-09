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
 * Обработчик JWT токенов для аутентификации и авторизации.
 * Поддерживает загрузку публичных ключей из Keycloak и верификацию токенов.
 * Кэширует публичные ключи для улучшения производительности.
 */
public class JwtHandler {

    private static final Logger log = LoggerFactory.getLogger(JwtHandler.class);

    /** Кэш публичных ключей по идентификатору ключа (kid) */
    private static final Map<String, CachedKey> publicKeys = new ConcurrentHashMap<>();
    
    /** Объект для синхронизации при загрузке ключей */
    private static final Object lock = new Object();
    
    /** Внутренний URL Keycloak для доступа из контейнерной сети */
    private static final String KEYCLOAK_INTERNAL_URL = System.getenv().getOrDefault("KEYCLOAK_INTERNAL_URL",
            "http://keycloak:8080");
    
    /** Публичный URL Keycloak для проверки issuer токена */
    private static final String KEYCLOAK_URL = System.getenv().getOrDefault("KEYCLOAK_URL", "http://localhost:8080");
    
    /** Название realm в Keycloak */
    private static final String REALM = System.getenv().getOrDefault("KEYCLOAK_REALM", "student");
    
    /** Объект для парсинга JSON */
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

    /** Время жизни кэша публичных ключей */
    private static final Duration CACHE_TTL = Duration.ofHours(1);

    /**
     * Запись для хранения публичного ключа и времени его загрузки.
     * 
     * @param key публичный RSA ключ
     * @param loadedAt время загрузки ключа
     */
    private record CachedKey(RSAPublicKey key, Instant loadedAt) {
    }

    /**
     * Гарантирует, что публичный ключ для указанного kid загружен и актуален.
     * Если ключ отсутствует или устарел, происходит загрузка из Keycloak.
     *
     * @param kid идентификатор ключа (Key ID) из JWT токена
     * @throws Exception если произошла ошибка при загрузке ключей
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
     * Загружает публичные ключи из JWKS эндпоинта Keycloak.
     * Парсит ответ и сохраняет RSA ключи в кэш.
     *
     * @throws Exception если произошла ошибка при загрузке или парсинге ключей
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
        if (!keysNode.isArray())
            throw new RuntimeException("JWKS не содержит keys");

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

            publicKeys.put(kid, new CachedKey(key, Instant.now()));
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
     * Проверяет наличие и валидность токена, извлекает claims и сохраняет их в контексте.
     *
     * @return Handler для использования в маршрутах Javalin
     * 
     * @throws UnauthorizedResponse если:
     *         - отсутствует заголовок Authorization
     *         - некорректный формат заголовка Authorization
     *         - не удалось загрузить публичные ключи
     *         - не найден ключ для указанного kid
     *         - токен не прошел верификацию (истек, неверная подпись и т.д.)
     * 
     * Устанавливает следующие атрибуты в контексте при успешной аутентификации:
     *   - userId: subject токена
     *   - username: preferred_username из claims
     *   - email: email из claims
     *   - jwt: полный объект DecodedJWT
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