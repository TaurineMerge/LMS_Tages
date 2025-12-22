package com.example.lms.shared.infrastructure;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.config.MinioConfig;
import com.example.lms.test_attempt.api.dto.AttemptSnapshotDto;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;

import io.minio.BucketExistsArgs;
import io.minio.GetObjectArgs;
import io.minio.MakeBucketArgs;
import io.minio.MinioClient;
import io.minio.PutObjectArgs;
import io.minio.RemoveObjectArgs;

/**
 * Сервис для работы с MinIO Object Storage.
 * <p>
 * Отвечает за сохранение, получение и удаление снимков попыток тестов в MinIO.
 * Снимки хранятся в формате JSON и содержат полную информацию о попытке.
 * 
 * <h2>Структура хранилища:</h2>
 * 
 * <pre>
 * bucket: snapshots/
 *   ├── {student_id}/
 *   │   ├── {test_id}/
 *   │   │   ├── {attempt_id}.json
 * </pre>
 * 
 * <h2>Операции:</h2>
 * <ul>
 * <li>Сохранение снимка при завершении попытки</li>
 * <li>Получение снимка для просмотра результатов</li>
 * <li>Удаление снимка при удалении попытки</li>
 * </ul>
 * 
 * @see AttemptSnapshotDto
 * @see MinioConfig
 */
public class MinioStorage {

    private static final Logger logger = LoggerFactory.getLogger(MinioStorage.class);

    private final MinioClient minioClient;
    private final MinioConfig config;
    private final ObjectMapper objectMapper;

    /**
     * Создаёт сервис для работы с MinIO.
     * <p>
     * Инициализирует клиент MinIO и проверяет наличие bucket.
     * Если bucket не существует, создаёт его автоматически.
     * 
     * @param config конфигурация MinIO
     */
    public MinioStorage(MinioConfig config) {
        this.config = config;
        this.minioClient = config.createClient();
        this.objectMapper = new ObjectMapper()
                .enable(SerializationFeature.INDENT_OUTPUT); // Красивый JSON

        ensureBucketExists();
        logger.info("MinioStorageService инициализирован успешно");
    }

    /**
     * Сохраняет снимок попытки в MinIO.
     * <p>
     * Выполняет следующие действия:
     * <ol>
     * <li>Сериализует DTO в JSON</li>
     * <li>Формирует путь к объекту в MinIO</li>
     * <li>Загружает JSON в MinIO</li>
     * <li>Возвращает URL для доступа к снимку</li>
     * </ol>
     * 
     * <h3>Пример пути:</h3>
     * {@code student_uuid/test_uuid/attempt_uuid.json}
     * 
     * <h3>Пример URL:</h3>
     * {@code http://localhost:9000/snapshots/student_uuid/test_uuid/attempt_uuid.json}
     * 
     * @param studentId UUID студента
     * @param testId    UUID теста
     * @param snapshot  снимок попытки для сохранения
     * @return URL для доступа к сохранённому снимку
     * @throws RuntimeException если не удалось сохранить снимок
     */
    public String saveSnapshot(UUID studentId, UUID testId, AttemptSnapshotDto snapshot) {
        try {
            // Сериализуем DTO в JSON
            String jsonContent = objectMapper.writeValueAsString(snapshot);
            byte[] jsonBytes = jsonContent.getBytes("UTF-8");

            // Формируем путь к объекту
            String objectPath = buildObjectPath(studentId, testId, snapshot.getAttemptId());

            logger.info("Сохранение снимка попытки {} в MinIO: {}",
                    snapshot.getAttemptId(), objectPath);

            // Загружаем в MinIO
            minioClient.putObject(
                    PutObjectArgs.builder()
                            .bucket(config.getBucket())
                            .object(objectPath)
                            .stream(new ByteArrayInputStream(jsonBytes), jsonBytes.length, -1)
                            .contentType("application/json")
                            .build());

            // Формируем URL для доступа
            String url = buildObjectUrl(studentId, testId, snapshot.getAttemptId());

            logger.info("Снимок успешно сохранён: {}", url);
            return url;

        } catch (Exception e) {
            logger.error("Ошибка при сохранении снимка попытки {}", snapshot.getAttemptId(), e);
            throw new RuntimeException("Не удалось сохранить снимок в MinIO", e);
        }
    }

    /**
     * Получает снимок попытки из MinIO.
     * <p>
     * Загружает JSON файл из MinIO и десериализует его в DTO.
     * 
     * @param studentId UUID студента
     * @param testId    UUID теста
     * @param attemptId UUID попытки
     * @return снимок попытки
     * @throws RuntimeException если снимок не найден или произошла ошибка
     */
    public AttemptSnapshotDto getSnapshot(UUID studentId, UUID testId, UUID attemptId) {
        try {
            String objectPath = buildObjectPath(studentId, testId, attemptId);

            logger.debug("Получение снимка попытки {} из MinIO: {}", attemptId, objectPath);

            // Скачиваем объект из MinIO
            InputStream stream = minioClient.getObject(
                    GetObjectArgs.builder()
                            .bucket(config.getBucket())
                            .object(objectPath)
                            .build());

            // Десериализуем JSON в DTO
            AttemptSnapshotDto snapshot = objectMapper.readValue(stream, AttemptSnapshotDto.class);

            logger.debug("Снимок успешно загружен");
            return snapshot;

        } catch (Exception e) {
            logger.error("Ошибка при получении снимка попытки {}", attemptId, e);
            throw new RuntimeException("Не удалось получить снимок из MinIO", e);
        }
    }

    /**
     * Удаляет снимок попытки из MinIO.
     * <p>
     * Используется при удалении попытки из системы.
     * 
     * @param studentId UUID студента
     * @param testId    UUID теста
     * @param attemptId UUID попытки
     * @return true если снимок был удалён, false если не найден
     */
    public boolean deleteSnapshot(UUID studentId, UUID testId, UUID attemptId) {
        try {
            String objectPath = buildObjectPath(studentId, testId, attemptId);

            logger.info("Удаление снимка попытки {} из MinIO: {}", attemptId, objectPath);

            minioClient.removeObject(
                    RemoveObjectArgs.builder()
                            .bucket(config.getBucket())
                            .object(objectPath)
                            .build());

            logger.info("Снимок успешно удалён");
            return true;

        } catch (Exception e) {
            logger.error("Ошибка при удалении снимка попытки {}", attemptId, e);
            return false;
        }
    }

    /**
     * Проверяет существование снимка в MinIO.
     * 
     * @param studentId UUID студента
     * @param testId    UUID теста
     * @param attemptId UUID попытки
     * @return true если снимок существует, false иначе
     */
    public boolean snapshotExists(UUID studentId, UUID testId, UUID attemptId) {
        try {
            String objectPath = buildObjectPath(studentId, testId, attemptId);

            minioClient.statObject(
                    io.minio.StatObjectArgs.builder()
                            .bucket(config.getBucket())
                            .object(objectPath)
                            .build());

            return true;

        } catch (Exception e) {
            return false;
        }
    }

    // ======================================================================
    // ВСПОМОГАТЕЛЬНЫЕ МЕТОДЫ
    // ======================================================================

    /**
     * Проверяет наличие bucket в MinIO и создаёт его при необходимости.
     * <p>
     * Вызывается автоматически при инициализации сервиса.
     */
    private void ensureBucketExists() {
        try {
            String bucketName = config.getBucket();

            boolean exists = minioClient.bucketExists(
                    BucketExistsArgs.builder()
                            .bucket(bucketName)
                            .build());

            if (!exists) {
                logger.info("Bucket {} не существует, создаём...", bucketName);

                minioClient.makeBucket(
                        MakeBucketArgs.builder()
                                .bucket(bucketName)
                                .build());

                logger.info("Bucket {} успешно создан", bucketName);
            } else {
                logger.info("Bucket {} уже существует", bucketName);
            }

        } catch (Exception e) {
            logger.error("Ошибка при проверке/создании bucket", e);
            throw new RuntimeException("Не удалось инициализировать MinIO bucket", e);
        }
    }

    /**
     * Формирует путь к объекту в MinIO.
     * <p>
     * Структура: {@code {student_id}/{test_id}/{attempt_id}.json}
     * 
     * @param studentId UUID студента
     * @param testId    UUID теста
     * @param attemptId UUID попытки
     * @return путь к объекту
     */
    private String buildObjectPath(UUID studentId, UUID testId, UUID attemptId) {
        return String.format("%s/%s/%s.json", studentId, testId, attemptId);
    }

    /**
     * Формирует URL для доступа к снимку.
     * <p>
     * Формат: {@code http(s)://endpoint/bucket/student_id/test_id/attempt_id.json}
     * 
     * <p>
     * <b>Примечание:</b> Для доступа к приватным объектам нужно использовать
     * presigned URL. В текущей реализации предполагается, что bucket публичный.
     * 
     * @param studentId UUID студента
     * @param testId    UUID теста
     * @param attemptId UUID попытки
     * @return полный URL к снимку
     */
    private String buildObjectUrl(UUID studentId, UUID testId, UUID attemptId) {
        String objectPath = buildObjectPath(studentId, testId, attemptId);
        return String.format("%s/%s", config.getBaseUrl(), objectPath);
    }

    /**
     * Генерирует временный presigned URL для безопасного доступа к снимку.
     * <p>
     * Presigned URL позволяет скачать приватный объект без аутентификации
     * в течение ограниченного времени.
     * 
     * @param studentId     UUID студента
     * @param testId        UUID теста
     * @param attemptId     UUID попытки
     * @param expirySeconds время действия URL в секундах
     * @return временный URL для скачивания
     * @throws RuntimeException если не удалось создать URL
     */
    public String generatePresignedUrl(UUID studentId, UUID testId, UUID attemptId,
            int expirySeconds) {
        try {
            String objectPath = buildObjectPath(studentId, testId, attemptId);

            String url = minioClient.getPresignedObjectUrl(
                    io.minio.GetPresignedObjectUrlArgs.builder()
                            .method(io.minio.http.Method.GET)
                            .bucket(config.getBucket())
                            .object(objectPath)
                            .expiry(expirySeconds)
                            .build());

            logger.debug("Создан presigned URL на {} секунд для попытки {}",
                    expirySeconds, attemptId);

            return url;

        } catch (Exception e) {
            logger.error("Ошибка при создании presigned URL для попытки {}", attemptId, e);
            throw new RuntimeException("Не удалось создать presigned URL", e);
        }
    }
}