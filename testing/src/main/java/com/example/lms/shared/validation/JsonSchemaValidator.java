package com.example.lms.shared.validation;

import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.networknt.schema.JsonSchema;
import com.networknt.schema.JsonSchemaFactory;
import com.networknt.schema.SpecVersion;
import com.networknt.schema.ValidationMessage;
import io.javalin.http.BadRequestResponse;
import io.javalin.http.Context;
import io.javalin.http.Handler;

import java.io.InputStream;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * Валидатор JSON данных на основе JSON Schema.
 * <p>
 * Использует библиотеку networknt/json-schema-validator для валидации входящих
 * JSON данных согласно определенным схемам.
 * 
 * @see JsonSchema
 * @see JsonSchemaFactory
 */
public class JsonSchemaValidator {

    private static final ObjectMapper OBJECT_MAPPER = new ObjectMapper();
    private static final JsonSchemaFactory SCHEMA_FACTORY = JsonSchemaFactory.getInstance(SpecVersion.VersionFlag.V7);

    /**
     * Создает Handler для валидации JSON тела запроса по схеме из файла.
     * <p>
     * Схема загружается из classpath (resources). Если валидация не проходит,
     * выбрасывается {@link BadRequestResponse} с описанием ошибок.
     * 
     * @param schemaPath путь к JSON Schema файлу в classpath (например,
     *                   "/schemas/user-schema.json")
     * @return Handler для валидации запроса
     * @throws BadRequestResponse       если JSON не соответствует схеме
     * @throws IllegalArgumentException если схема не найдена или невалидна
     * 
     * @see #validate(String)
     */
    public static Handler validate(String schemaPath) {
        // Загружаем схему из classpath
        InputStream schemaStream = JsonSchemaValidator.class.getResourceAsStream(schemaPath);
        if (schemaStream == null) {
            throw new IllegalArgumentException("JSON Schema не найдена: " + schemaPath);
        }

        try {
            JsonNode schemaNode = OBJECT_MAPPER.readTree(schemaStream);
            JsonSchema schema = SCHEMA_FACTORY.getSchema(schemaNode);

            return validateWithSchema(schema);
        } catch (Exception e) {
            throw new IllegalArgumentException("Ошибка загрузки JSON Schema: " + schemaPath, e);
        }
    }

    /**
     * Создает Handler для валидации JSON тела запроса по схеме из строки.
     * <p>
     * Полезно для простых схем или динамической генерации схем.
     * 
     * @param schemaJson JSON Schema в виде строки
     * @return Handler для валидации запроса
     * @throws BadRequestResponse       если JSON не соответствует схеме
     * @throws IllegalArgumentException если схема невалидна
     */
    public static Handler validateFromString(String schemaJson) {
        try {
            JsonNode schemaNode = OBJECT_MAPPER.readTree(schemaJson);
            JsonSchema schema = SCHEMA_FACTORY.getSchema(schemaNode);

            return validateWithSchema(schema);
        } catch (Exception e) {
            throw new IllegalArgumentException("Ошибка парсинга JSON Schema", e);
        }
    }

    /**
     * Создает Handler для валидации с уже загруженной схемой.
     * <p>
     * Внутренний метод, выполняющий непосредственную валидацию.
     * 
     * @param schema скомпилированная JSON Schema
     * @return Handler для валидации запроса
     */
    private static Handler validateWithSchema(JsonSchema schema) {
        return ctx -> {
            try {
                // Парсим тело запроса как JSON
                JsonNode jsonNode = OBJECT_MAPPER.readTree(ctx.body());

                // Выполняем валидацию
                Set<ValidationMessage> errors = schema.validate(jsonNode);

                if (!errors.isEmpty()) {
                    // Формируем читаемое сообщение об ошибках
                    String errorMessage = errors.stream()
                            .map(ValidationMessage::getMessage)
                            .collect(Collectors.joining(", "));

                    throw new BadRequestResponse("Ошибка валидации: " + errorMessage);
                }

                // Валидация пройдена, сохраняем распарсенный JSON в контекст
                ctx.attribute("validatedJson", jsonNode);

            } catch (BadRequestResponse e) {
                throw e;
            } catch (Exception e) {
                throw new BadRequestResponse("Невалидный JSON: " + e.getMessage());
            }
        };
    }

    /**
     * Получает валидированный JSON из контекста запроса.
     * <p>
     * Используется в контроллере после прохождения валидации для получения
     * уже распарсенного и проверенного JSON.
     * 
     * @param ctx контекст Javalin
     * @return валидированный JsonNode или null, если валидация не выполнялась
     * 
     * @see #validate(String)
     */
    public static JsonNode getValidatedJson(Context ctx) {
        return ctx.attribute("validatedJson");
    }

    /**
     * Валидирует JSON напрямую, без использования в качестве Handler.
     * <p>
     * Полезно для валидации данных внутри сервисов или тестов.
     * 
     * @param jsonNode   JSON для валидации
     * @param schemaPath путь к схеме в classpath
     * @return true если валидация пройдена
     * @throws ValidationException если валидация не пройдена
     */
    public static boolean validateJson(JsonNode jsonNode, String schemaPath) {
        InputStream schemaStream = JsonSchemaValidator.class.getResourceAsStream(schemaPath);
        if (schemaStream == null) {
            throw new IllegalArgumentException("JSON Schema не найдена: " + schemaPath);
        }

        try {
            JsonNode schemaNode = OBJECT_MAPPER.readTree(schemaStream);
            JsonSchema schema = SCHEMA_FACTORY.getSchema(schemaNode);

            Set<ValidationMessage> errors = schema.validate(jsonNode);

            if (!errors.isEmpty()) {
                String errorMessage = errors.stream()
                        .map(ValidationMessage::getMessage)
                        .collect(Collectors.joining(", "));

                throw new ValidationException("Валидация не пройдена: " + errorMessage);
            }

            return true;
        } catch (Exception e) {
            throw new ValidationException("Ошибка валидации", e);
        }
    }

    /**
     * Исключение для ошибок валидации.
     */
    public static class ValidationException extends RuntimeException {
        public ValidationException(String message) {
            super(message);
        }

        public ValidationException(String message, Throwable cause) {
            super(message, cause);
        }
    }
}