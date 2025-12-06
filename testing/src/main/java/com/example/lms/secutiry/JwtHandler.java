package com.example.lms.secutiry;

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

    private static RSAPublicKey publicKey;
    private static final String KEYCLOAK_URL = System.getenv().getOrDefault("KEYCLOAK_URL", "http://keycloak:8080");
    private static final String REALM = System.getenv().getOrDefault("KEYCLOAK_REALM", "master");

    static {
        try {
            loadPublicKeyFromCerts();
        } catch (Exception e) {
            System.err.println("Failed to load Keycloak public key: " + e.getMessage());
            e.printStackTrace();
        }
    }

    private static void loadPublicKeyFromCerts() throws Exception {
        HttpClient client = HttpClient.newHttpClient();

        // Получаем JWKS (JSON Web Key Set) от Keycloak
        String certsUrl = KEYCLOAK_URL + "/realms/" + REALM + "/protocol/openid-connect/certs";
        HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(certsUrl))
                .GET()
                .build();

        HttpResponse<String> response = client.send(request, HttpResponse.BodyHandlers.ofString());
        String body = response.body();

        // Простой парсинг JSON для получения n и e (modulus и exponent)
        // Keycloak возвращает массив keys, берём первый ключ с use="sig"
        String nValue = extractJsonValue(body, "\"n\":\"", "\"");
        String eValue = extractJsonValue(body, "\"e\":\"", "\"");

        if (nValue == null || eValue == null) {
            throw new RuntimeException("Could not extract public key from Keycloak certs endpoint");
        }

        // Декодируем Base64Url в RSA компоненты
        byte[] nBytes = Base64.getUrlDecoder().decode(nValue);
        byte[] eBytes = Base64.getUrlDecoder().decode(eValue);

        BigInteger modulus = new BigInteger(1, nBytes);
        BigInteger exponent = new BigInteger(1, eBytes);

        RSAPublicKeySpec spec = new RSAPublicKeySpec(modulus, exponent);
        KeyFactory factory = KeyFactory.getInstance("RSA");
        publicKey = (RSAPublicKey) factory.generatePublic(spec);

        System.out.println("Successfully loaded Keycloak public key from realm: " + REALM);
    }

    private static String extractJsonValue(String json, String startMarker, String endMarker) {
        int start = json.indexOf(startMarker);
        if (start == -1)
            return null;
        start += startMarker.length();

        int end = json.indexOf(endMarker, start);
        if (end == -1)
            return null;

        return json.substring(start, end);
    }

    public static Handler authenticate() {
        return ctx -> {
            String authHeader = ctx.header("Authorization");

            if (authHeader == null || !authHeader.startsWith("Bearer ")) {
                throw new UnauthorizedResponse("Missing or invalid Authorization header");
            }

            String token = authHeader.substring(7);

            try {
                Algorithm algorithm = Algorithm.RSA256(publicKey, null);
                DecodedJWT jwt = JWT.require(algorithm)
                        .withIssuer(KEYCLOAK_URL + "/realms/" + REALM)
                        .build()
                        .verify(token);

                // Сохраняем информацию о пользователе в контексте
                ctx.attribute("userId", jwt.getSubject());
                ctx.attribute("username", jwt.getClaim("preferred_username").asString());
                ctx.attribute("email", jwt.getClaim("email").asString());
                ctx.attribute("roles", jwt.getClaim("realm_access").asMap().get("roles"));
                ctx.attribute("jwt", jwt);

                System.out.println("Authenticated user: " + jwt.getClaim("preferred_username").asString());

            } catch (JWTVerificationException e) {
                throw new UnauthorizedResponse("Invalid JWT token: " + e.getMessage());
            }
        };
    }
}