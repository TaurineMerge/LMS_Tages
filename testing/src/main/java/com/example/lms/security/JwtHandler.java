package com.example.lms.security;

import io.javalin.http.Handler;
import io.javalin.http.UnauthorizedResponse;
import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.interfaces.DecodedJWT;
import com.auth0.jwt.exceptions.JWTVerificationException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;

import java.math.BigInteger;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.security.KeyFactory;
import java.security.interfaces.RSAPublicKey;
import java.security.spec.RSAPublicKeySpec;
import java.util.Base64;
import java.time.Duration;

public class JwtHandler {

    private static volatile RSAPublicKey publicKey = null;
    private static final Object lock = new Object();
    private static final String KEYCLOAK_INTERNAL_URL = System.getenv().getOrDefault("KEYCLOAK_INTERNAL_URL",
            "http://keycloak:8080");
    private static final String KEYCLOAK_URL = System.getenv().getOrDefault("KEYCLOAK_URL", "http://localhost:8080");
    private static final String REALM = System.getenv().getOrDefault("KEYCLOAK_REALM", "lms-realm");

    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();

    private static void ensurePublicKeyLoaded(String kid) throws Exception {
        if (publicKey == null) {
            synchronized (lock) {
                if (publicKey == null) {
                    loadPublicKeyFromCerts(kid);
                }
            }
        }
    }

    private static void loadPublicKeyFromCerts(String kid) throws Exception {
        try {
            HttpClient client = HttpClient.newBuilder()
                    .connectTimeout(Duration.ofSeconds(10))
                    .build();

            String certsUrl = KEYCLOAK_INTERNAL_URL + "/realms/" + REALM + "/protocol/openid-connect/certs";
            System.out.println("Получение JWKS с: " + certsUrl);

            HttpRequest request = HttpRequest.newBuilder()
                    .uri(URI.create(certsUrl))
                    .timeout(Duration.ofSeconds(10))
                    .GET()
                    .build();

            System.out.println("Отправка HTTP-запроса...");
            HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
            System.out.println("Получен HTTP-ответ: " + response.statusCode());

            if (response.statusCode() != 200) {
                throw new RuntimeException("Keycloak вернул статус " + response.statusCode());
            }

            String body = response.body();
            System.out.println("Длина тела ответа: " + (body != null ? body.length() : "null"));

            if (body == null || body.isEmpty()) {
                throw new RuntimeException("Тело ответа пустое");
            }

            JsonNode root = OBJECT_MAPPER.readTree(body);
            JsonNode keysNode = root.path("keys");
            if (keysNode == null || !keysNode.isArray() || keysNode.size() == 0) {
                throw new RuntimeException("JWKS не содержит keys");
            }

            JsonNode chosenKey = null;
            for (JsonNode keyNode : keysNode) {
                if (keyNode.path("kid").asText("").equals(kid)) {
                    chosenKey = keyNode;
                    break;
                }
            }

            if (chosenKey == null) {
                throw new RuntimeException("Не найден ключ в JWKS с kid=" + kid);
            }

            String nValue = chosenKey.path("n").asText(null);
            String eValue = chosenKey.path("e").asText(null);

            System.out.println("Извлечено n: " + (nValue != null ? ("длина=" + nValue.length()) : "null"));
            System.out.println("Извлечено e: " + (eValue != null ? eValue : "null"));

            if (nValue == null || eValue == null) {
                throw new RuntimeException("Не удалось извлечь n или e из ключа");
            }

            System.out.println("Декодирование base64...");
            byte[] nBytes;
            byte[] eBytes;
            try {
                nBytes = Base64.getUrlDecoder().decode(nValue);
                eBytes = Base64.getUrlDecoder().decode(eValue);
            } catch (IllegalArgumentException ex) {
                throw new RuntimeException("Ошибка при base64url декодировании n/e: " + ex.getMessage(), ex);
            }

            System.out.println("Создание BigInteger...");
            BigInteger modulus = new BigInteger(1, nBytes);
            BigInteger exponent = new BigInteger(1, eBytes);
            RSAPublicKeySpec spec = new RSAPublicKeySpec(modulus, exponent);

            KeyFactory factory = KeyFactory.getInstance("RSA");
            publicKey = (RSAPublicKey) factory.generatePublic(spec);

            if (publicKey == null) {
                throw new RuntimeException("Публичный ключ равен null после генерации");
            }

            System.out.println("Публичный ключ успешно создан: " + publicKey.getModulus().bitLength() + " бит");

        } catch (Exception e) {
            System.err.println("Исключение при загрузке ключа: " + e.getClass().getName());
            System.err.println("Сообщение исключения: " + e.getMessage());
            e.printStackTrace();
            throw e;
        }
    }

    private static String extractKid(String token) {
        try {
            DecodedJWT jwt = JWT.decode(token);
            return jwt.getKeyId();
        } catch (Exception e) {
            throw new RuntimeException("Не удалось извлечь kid из токена: " + e.getMessage(), e);
        }
    }

    public static Handler authenticate() {
        return ctx -> {
            String authHeader = ctx.header("Authorization");

            if (authHeader == null || !authHeader.startsWith("Bearer ")) {
                throw new UnauthorizedResponse("Отсутствует или неверный заголовок Authorization");
            }

            String token = authHeader.substring(7);

            String kid = extractKid(token);
            System.out.println("KID токена: " + kid);

            try {
                ensurePublicKeyLoaded(kid);
            } catch (Exception e) {
                System.err.println("Не удалось загрузить публичный ключ: " + e.getMessage());
                throw new UnauthorizedResponse("Сервер не готов: " + e.getMessage());
            }

            try {
                Algorithm algorithm = Algorithm.RSA256(publicKey, null);

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

                System.out.println("Аутентифицирован пользователь: "
                        + jwt.getClaim("preferred_username").asString());

            } catch (JWTVerificationException e) {
                System.err.println("Ошибка проверки JWT: " + e.getMessage());
                throw new UnauthorizedResponse("Неверный JWT токен: " + e.getMessage());
            }
        };
    }
}
