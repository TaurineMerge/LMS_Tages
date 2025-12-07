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

public class JwtHandler {

    private static final Logger log = LoggerFactory.getLogger(JwtHandler.class);

    private static final Map<String, CachedKey> publicKeys = new ConcurrentHashMap<>();
    private static final Object lock = new Object();
    private static final String KEYCLOAK_INTERNAL_URL = System.getenv().getOrDefault("KEYCLOAK_INTERNAL_URL",
            "http://keycloak:8080");
    private static final String KEYCLOAK_URL = System.getenv().getOrDefault("KEYCLOAK_URL", "http://localhost:8080");
    private static final String REALM = System.getenv().getOrDefault("KEYCLOAK_REALM", "lms-realm");
    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

    private static final Duration CACHE_TTL = Duration.ofHours(1);

    private record CachedKey(RSAPublicKey key, Instant loadedAt) {
    }

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

    private static String extractKid(String token) {
        DecodedJWT jwt = JWT.decode(token);
        return jwt.getKeyId();
    }

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

                var realmAccess = jwt.getClaim("realm_access").asMap();
                if (realmAccess != null && realmAccess.containsKey("roles")) {
                    ctx.attribute("roles", realmAccess.get("roles"));
                }

                ctx.attribute("jwt", jwt);

            } catch (JWTVerificationException e) {
                throw new UnauthorizedResponse("Неверный JWT токен: " + e.getMessage());
            }
        };
    }
}
