package com.example.lms.security;

import io.javalin.http.Handler;
import io.javalin.http.UnauthorizedResponse;
import com.auth0.jwt.JWT;
import com.auth0.jwt.algorithms.Algorithm;
import com.auth0.jwt.interfaces.DecodedJWT;
import com.auth0.jwt.exceptions.JWTVerificationException;

import java.math.BigInteger;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.security.KeyFactory;
import java.security.interfaces.RSAPublicKey;
import java.security.spec.RSAPublicKeySpec;
import java.util.Base64;

public class JwtHandler {

    private static volatile RSAPublicKey publicKey = null;
    private static final Object lock = new Object();
    private static final String KEYCLOAK_INTERNAL_URL = System.getenv().getOrDefault("KEYCLOAK_INTERNAL_URL",
            "http://keycloak:8080");
    private static final String KEYCLOAK_URL = System.getenv().getOrDefault("KEYCLOAK_URL", "http://localhost:8080");
    private static final String REALM = System.getenv().getOrDefault("KEYCLOAK_REALM", "lms-realm");

    private static void ensurePublicKeyLoaded() throws Exception {
        if (publicKey == null) {
            synchronized (lock) {
                if (publicKey == null) {
                    System.out.println("Загрузка публичного ключа Keycloak...");
                    loadPublicKeyFromCerts();
                    System.out.println("Публичный ключ Keycloak успешно загружен из realm: " + REALM);
                }
            }
        }
    }

    private static void loadPublicKeyFromCerts() throws Exception {
        try {
            HttpClient client = HttpClient.newBuilder()
                    .connectTimeout(java.time.Duration.ofSeconds(10))
                    .build();

            String certsUrl = KEYCLOAK_INTERNAL_URL + "/realms/" + REALM + "/protocol/openid-connect/certs";
            System.out.println("Получение JWKS с: " + certsUrl);

            HttpRequest request = HttpRequest.newBuilder()
                    .uri(URI.create(certsUrl))
                    .timeout(java.time.Duration.ofSeconds(10))
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

            int firstAqab = body.indexOf("\"e\":\"AQAB\"");
            System.out.println("Первый AQAB на позиции: " + firstAqab);

            if (firstAqab == -1) {
                throw new RuntimeException("Не удалось найти RSA ключи");
            }

            int secondKeyStart = body.indexOf("\"n\":\"", firstAqab + 20);
            System.out.println("Второй ключ 'n' начинается на: " + secondKeyStart);

            if (secondKeyStart == -1) {
                System.out.println("Найден только один ключ, используем его");
                secondKeyStart = body.indexOf("\"n\":\"");
            }

            String remainingBody = body.substring(secondKeyStart);
            System.out.println("Извлечение из: " + remainingBody.substring(0, Math.min(100, remainingBody.length())));

            String nValue = extractJsonValue(remainingBody, "\"n\":\"", "\"");
            String eValue = extractJsonValue(remainingBody, "\"e\":\"", "\"");

            System.out.println("Извлечено n: " + (nValue != null ? ("длина=" + nValue.length()) : "null"));
            System.out.println("Извлечено e: " + (eValue != null ? eValue : "null"));

            if (nValue == null || eValue == null) {
                throw new RuntimeException("Не удалось извлечь n или e из ключа");
            }

            System.out.println("Декодирование base64...");
            byte[] nBytes = Base64.getUrlDecoder().decode(nValue);
            byte[] eBytes = Base64.getUrlDecoder().decode(eValue);

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

    private static String extractJsonValue(String json, String startMarker, String endMarker) {
        try {
            int start = json.indexOf(startMarker);
            if (start == -1) {
                System.err.println("Маркер начала не найден: " + startMarker);
                return null;
            }
            start += startMarker.length();

            int end = json.indexOf(endMarker, start);
            if (end == -1) {
                System.err.println("Маркер конца не найден: " + endMarker);
                return null;
            }

            String result = json.substring(start, end);
            System.out.println("Извлечено значение между маркерами (длина=" + result.length() + ")");
            return result;
        } catch (Exception e) {
            System.err.println("Ошибка в extractJsonValue: " + e.getMessage());
            e.printStackTrace();
            return null;
        }
    }

    public static Handler authenticate() {
        return ctx -> {
            try {
                ensurePublicKeyLoaded();
            } catch (Exception e) {
                System.err.println("Не удалось загрузить публичный ключ: " + e.getMessage());
                e.printStackTrace();
                throw new UnauthorizedResponse("Сервер не готов: " + e.getMessage());
            }

            String authHeader = ctx.header("Authorization");

            if (authHeader == null || !authHeader.startsWith("Bearer ")) {
                throw new UnauthorizedResponse("Отсутствует или неверный заголовок Authorization");
            }

            String token = authHeader.substring(7);

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

                System.out
                        .println("Аутентифицированный пользователь: " + jwt.getClaim("preferred_username").asString());

            } catch (JWTVerificationException e) {
                System.err.println("Ошибка проверки JWT: " + e.getMessage());
                throw new UnauthorizedResponse("Неверный JWT токен: " + e.getMessage());
            }
        };
    }
}
