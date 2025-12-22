package com.example.lms.config;

import io.github.cdimascio.dotenv.Dotenv;
import io.minio.MinioClient;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Конфигурация для подключения к MinIO Object Storage.
 * <p>
 * MinIO используется для хранения снимков (snapshots) попыток прохождения
 * тестов.
 * Снимки содержат полную информацию о вопросах, ответах и выборе студента.
 * 
 * <h2>Переменные окружения:</h2>
 * <ul>
 * <li>{@code MINIO_ENDPOINT} - адрес MinIO сервера (например,
 * localhost:9000)</li>
 * <li>{@code MINIO_ACCESS_KEY} - ключ доступа (по умолчанию: minioadmin)</li>
 * <li>{@code MINIO_SECRET_KEY} - секретный ключ (по умолчанию: minioadmin)</li>
 * <li>{@code MINIO_BUCKET} - название bucket для снимков (по умолчанию:
 * test-attempts)</li>
 * <li>{@code MINIO_USE_SSL} - использовать SSL (по умолчанию: false)</li>
 * </ul>
 * 
 * <h2>Структура хранилища:</h2>
 * 
 * <pre>
 * bucket: test-attempts/
 *   ├── {student_id}/
 *   │   ├── {test_id}/
 *   │   │   ├── {attempt_id}.json
 * </pre>
 * 
 * @see io.minio.MinioClient
 * @see com.example.lms.test_attempt.infrastructure.storage.MinioStorageService
 */
public class MinioConfig {

    private static final Logger logger = LoggerFactory.getLogger(MinioConfig.class);
    private static final Dotenv dotenv = Dotenv.configure().ignoreIfMissing().load();

    /** Адрес MinIO сервера (host:port) */
    private final String endpoint;

    /** Access Key для аутентификации */
    private final String accessKey;

    /** Secret Key для аутентификации */
    private final String secretKey;

    /** Название bucket для хранения снимков */
    private final String bucket;

    /** Использовать SSL для подключения */
    private final boolean useSSL;

    /**
     * Создаёт конфигурацию MinIO из переменных окружения.
     * <p>
     * Загружает настройки из файла {@code .env} или системных переменных.
     * Использует значения по умолчанию, если переменные не установлены.
     */
    public MinioConfig() {
        this.endpoint = dotenv.get("MINIO_ENDPOINT");
        this.accessKey = dotenv.get("MINIO_ACCESS_KEY");
        this.secretKey = dotenv.get("MINIO_SECRET_KEY");
        this.bucket = dotenv.get("MINIO_BUCKET");
        this.useSSL = Boolean.parseBoolean(dotenv.get("MINIO_USE_SSL"));

        logger.info("MinIO конфигурация загружена: endpoint={}, bucket={}, useSSL={}",
                endpoint, bucket, useSSL);
    }

    /**
     * Создаёт клиент MinIO с текущей конфигурацией.
     * <p>
     * Клиент используется для выполнения операций с объектами в MinIO:
     * загрузка, скачивание, удаление файлов.
     * 
     * @return настроенный клиент MinioClient
     * @throws IllegalStateException если не удалось создать клиент
     */
    public MinioClient createClient() {
        try {
            MinioClient client = MinioClient.builder()
                    .endpoint(endpoint, useSSL ? 443 : 9000, useSSL)
                    .credentials(accessKey, secretKey)
                    .build();

            logger.info("MinIO клиент успешно создан");
            return client;

        } catch (Exception e) {
            logger.error("Ошибка при создании MinIO клиента", e);
            throw new IllegalStateException("Не удалось создать MinIO клиент", e);
        }
    }

    /**
     * Получает значение переменной окружения или возвращает значение по умолчанию.
     * 
     * @param key          название переменной окружения
     * @param defaultValue значение по умолчанию
     * @return значение переменной или defaultValue
     */
    private String getEnvOrDefault(String key, String defaultValue) {
        String value = dotenv.get(key);
        return value != null ? value : defaultValue;
    }

    // ======================================================================
    // GETTERS
    // ======================================================================

    /**
     * Получить адрес MinIO сервера.
     * 
     * @return endpoint в формате host:port
     */
    public String getEndpoint() {
        return endpoint;
    }

    /**
     * Получить Access Key.
     * 
     * @return ключ доступа для аутентификации
     */
    public String getAccessKey() {
        return accessKey;
    }

    /**
     * Получить Secret Key.
     * 
     * @return секретный ключ для аутентификации
     */
    public String getSecretKey() {
        return secretKey;
    }

    /**
     * Получить название bucket.
     * 
     * @return название bucket для снимков попыток
     */
    public String getBucket() {
        return bucket;
    }

    /**
     * Проверить, используется ли SSL.
     * 
     * @return true если используется HTTPS, false для HTTP
     */
    public boolean isUseSSL() {
        return useSSL;
    }

    /**
     * Получить базовый URL для доступа к объектам.
     * <p>
     * Формат: {@code http(s)://endpoint/bucket/}
     * 
     * @return базовый URL для объектов
     */
    public String getBaseUrl() {
        String protocol = useSSL ? "https" : "http";
        return String.format("%s://%s/%s", protocol, endpoint, bucket);
    }

    @Override
    public String toString() {
        return "MinioConfig{" +
                "endpoint='" + endpoint + '\'' +
                ", bucket='" + bucket + '\'' +
                ", useSSL=" + useSSL +
                ", accessKey='***'" + // Не показываем ключи в логах
                ", secretKey='***'" +
                '}';
    }
}