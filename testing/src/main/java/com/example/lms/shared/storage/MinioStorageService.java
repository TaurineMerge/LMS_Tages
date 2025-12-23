package com.example.lms.shared.storage;

import java.io.ByteArrayInputStream;
import java.io.InputStream;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.Optional;
import java.util.concurrent.TimeUnit;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.config.MinioConfig;
import com.example.lms.shared.storage.dto.FileMetadata;
import com.example.lms.shared.storage.exceptions.StorageException;
import com.example.lms.shared.storage.dto.UploadResult;
import com.example.lms.shared.storage.validator.ImageValidator;

import io.minio.BucketExistsArgs;
import io.minio.GetObjectArgs;
import io.minio.GetPresignedObjectUrlArgs;
import io.minio.ListObjectsArgs;
import io.minio.MakeBucketArgs;
import io.minio.MinioClient;
import io.minio.PutObjectArgs;
import io.minio.RemoveObjectArgs;
import io.minio.Result;
import io.minio.StatObjectArgs;
import io.minio.StatObjectResponse;
import io.minio.http.Method;
import io.minio.messages.Item;

/**
 * Реализация StorageService для работы с MinIO.
 * <p>
 * Обеспечивает загрузку, скачивание и управление файлами в MinIO,
 * с автоматическим созданием buckets и валидацией файлов.
 */
public class MinioStorageService implements StorageServiceInterface {

    private static final Logger logger = LoggerFactory.getLogger(MinioStorageService.class);

    private final MinioClient minioClient;
    private final MinioConfig config;
    private final ImageValidator imageValidator;

    // Названия buckets
    private static final String SNAPSHOT_BUCKET = "test-attempts";
    private static final String IMAGE_BUCKET = "test-images";

    // Префиксы путей
    private static final String SNAPSHOT_PREFIX = "snapshots/";
    private static final String QUESTION_PREFIX = "questions/";
    private static final String ANSWER_PREFIX = "answers/";

    // Настройки presigned URLs
    private static final int DEFAULT_URL_EXPIRY_HOURS = 1;

    public MinioStorageService(MinioConfig config) {
        this.config = config;
        this.minioClient = config.createClient();
        this.imageValidator = new ImageValidator();

        // Инициализация buckets при запуске
        initializeBuckets();

        logger.info("MinioStorageService инициализирован");
    }

    /**
     * Создает необходимые buckets, если они не существуют.
     */
    private void initializeBuckets() {
        try {
            ensureBucketExists(SNAPSHOT_BUCKET);
            ensureBucketExists(IMAGE_BUCKET);
            logger.info("Все buckets проверены и готовы к использованию");
        } catch (Exception e) {
            logger.error("Ошибка при инициализации buckets", e);
            throw new StorageException("Не удалось инициализировать хранилище", e);
        }
    }

    /**
     * Проверяет существование bucket и создает его при необходимости.
     */
    private void ensureBucketExists(String bucketName) {
        try {
            boolean exists = minioClient.bucketExists(
                    BucketExistsArgs.builder().bucket(bucketName).build());

            if (!exists) {
                minioClient.makeBucket(
                        MakeBucketArgs.builder().bucket(bucketName).build());
                logger.info("Bucket '{}' создан", bucketName);
            } else {
                logger.debug("Bucket '{}' уже существует", bucketName);
            }
        } catch (Exception e) {
            throw new StorageException("Ошибка при проверке/создании bucket: " + bucketName, e);
        }
    }

    @Override
    public UploadResult upload(String bucketName, String objectPath, InputStream inputStream,
            String contentType, long size) {
        try {
            logger.debug("Загрузка файла: bucket={}, path={}, size={}", bucketName, objectPath, size);

            ensureBucketExists(bucketName);

            minioClient.putObject(
                    PutObjectArgs.builder()
                            .bucket(bucketName)
                            .object(objectPath)
                            .stream(inputStream, size, -1)
                            .contentType(contentType)
                            .build());

            String url = String.format("%s/%s/%s", config.getBaseUrl(), bucketName, objectPath);

            logger.info("Файл успешно загружен: {}", objectPath);
            return new UploadResult(objectPath, url, size, contentType, true);

        } catch (Exception e) {
            logger.error("Ошибка при загрузке файла: {}", objectPath, e);
            throw new StorageException("Не удалось загрузить файл: " + objectPath, e);
        }
    }

    @Override
    public UploadResult uploadSnapshot(String studentId, String testId, String attemptId, String snapshotJson) {
        String objectPath = buildSnapshotPath(studentId, testId, attemptId);
        byte[] bytes = snapshotJson.getBytes(StandardCharsets.UTF_8);

        try (InputStream inputStream = new ByteArrayInputStream(bytes)) {
            return upload(SNAPSHOT_BUCKET, objectPath, inputStream, "application/json", bytes.length);
        } catch (Exception e) {
            throw new StorageException("Ошибка при загрузке снепшота", e);
        }
    }

    @Override
    public UploadResult uploadImage(String entityType, String entityId, String imageId,
            InputStream inputStream, String contentType, long size) {
        // Валидация изображения
        imageValidator.validate(inputStream, contentType, size);

        String objectPath = buildImagePath(entityType, entityId, imageId);
        return upload(IMAGE_BUCKET, objectPath, inputStream, contentType, size);
    }

    @Override
    public Optional<InputStream> download(String bucketName, String objectPath) {
        try {
            logger.debug("Скачивание файла: bucket={}, path={}", bucketName, objectPath);

            InputStream stream = minioClient.getObject(
                    GetObjectArgs.builder()
                            .bucket(bucketName)
                            .object(objectPath)
                            .build());

            logger.info("Файл успешно скачан: {}", objectPath);
            return Optional.of(stream);

        } catch (Exception e) {
            logger.warn("Файл не найден или ошибка при скачивании: {}", objectPath, e);
            return Optional.empty();
        }
    }

    @Override
    public Optional<String> downloadSnapshot(String studentId, String testId, String attemptId) {
        String objectPath = buildSnapshotPath(studentId, testId, attemptId);

        try (InputStream stream = download(SNAPSHOT_BUCKET, objectPath).orElse(null)) {
            if (stream == null) {
                return Optional.empty();
            }

            String json = new String(stream.readAllBytes(), StandardCharsets.UTF_8);
            return Optional.of(json);

        } catch (Exception e) {
            logger.error("Ошибка при чтении снепшота", e);
            return Optional.empty();
        }
    }

    @Override
    public Optional<InputStream> downloadImage(String entityType, String entityId, String imageId) {
        String objectPath = buildImagePath(entityType, entityId, imageId);
        return download(IMAGE_BUCKET, objectPath);
    }

    @Override
    public Optional<FileMetadata> getMetadata(String bucketName, String objectPath) {
        try {
            StatObjectResponse stat = minioClient.statObject(
                    StatObjectArgs.builder()
                            .bucket(bucketName)
                            .object(objectPath)
                            .build());

            FileMetadata metadata = new FileMetadata(
                    objectPath,
                    stat.size(),
                    stat.contentType(),
                    stat.lastModified().toInstant());

            return Optional.of(metadata);

        } catch (Exception e) {
            logger.debug("Метаданные не найдены для: {}", objectPath);
            return Optional.empty();
        }
    }

    @Override
    public boolean delete(String bucketName, String objectPath) {
        try {
            logger.debug("Удаление файла: bucket={}, path={}", bucketName, objectPath);

            minioClient.removeObject(
                    RemoveObjectArgs.builder()
                            .bucket(bucketName)
                            .object(objectPath)
                            .build());

            logger.info("Файл успешно удален: {}", objectPath);
            return true;

        } catch (Exception e) {
            logger.warn("Ошибка при удалении файла: {}", objectPath, e);
            return false;
        }
    }

    @Override
    public boolean deleteSnapshot(String studentId, String testId, String attemptId) {
        String objectPath = buildSnapshotPath(studentId, testId, attemptId);
        return delete(SNAPSHOT_BUCKET, objectPath);
    }

    @Override
    public boolean deleteImage(String entityType, String entityId, String imageId) {
        String objectPath = buildImagePath(entityType, entityId, imageId);
        return delete(IMAGE_BUCKET, objectPath);
    }

    @Override
    public List<String> listObjects(String bucketName, String prefix) {
        List<String> objects = new ArrayList<>();

        try {
            Iterable<Result<Item>> results = minioClient.listObjects(
                    ListObjectsArgs.builder()
                            .bucket(bucketName)
                            .prefix(prefix)
                            .recursive(true)
                            .build());

            for (Result<Item> result : results) {
                Item item = result.get();
                objects.add(item.objectName());
            }

            logger.debug("Найдено {} объектов с префиксом: {}", objects.size(), prefix);
            return objects;

        } catch (Exception e) {
            logger.error("Ошибка при получении списка объектов", e);
            throw new StorageException("Не удалось получить список объектов", e);
        }
    }

    @Override
    public List<String> listSnapshots(String studentId, String testId) {
        String prefix = String.format("%s%s/%s/", SNAPSHOT_PREFIX, studentId, testId);
        List<String> paths = listObjects(SNAPSHOT_BUCKET, prefix);

        // Извлекаем только ID попыток из полных путей
        return paths.stream()
                .map(path -> path.substring(path.lastIndexOf('/') + 1))
                .map(filename -> filename.replace(".json", ""))
                .toList();
    }

    @Override
    public boolean exists(String bucketName, String objectPath) {
        return getMetadata(bucketName, objectPath).isPresent();
    }

    @Override
    public String getPresignedUrl(String bucketName, String objectPath, int expiresInSeconds) {
        try {
            String url = minioClient.getPresignedObjectUrl(
                    GetPresignedObjectUrlArgs.builder()
                            .method(Method.GET)
                            .bucket(bucketName)
                            .object(objectPath)
                            .expiry(expiresInSeconds, TimeUnit.SECONDS)
                            .build());

            logger.debug("Создан presigned URL для: {}", objectPath);
            return url;

        } catch (Exception e) {
            logger.error("Ошибка при создании presigned URL", e);
            throw new StorageException("Не удалось создать presigned URL", e);
        }
    }

    @Override
    public String getImageUrl(String entityType, String entityId, String imageId) {
        String objectPath = buildImagePath(entityType, entityId, imageId);
        int expirySeconds = (int) TimeUnit.HOURS.toSeconds(DEFAULT_URL_EXPIRY_HOURS);
        return getPresignedUrl(IMAGE_BUCKET, objectPath, expirySeconds);
    }

    @Override
    public String getExtensionFromContentType(String contentType) {
        return imageValidator.getExtensionFromContentType(contentType);
    }

    // ======================================================================
    // HELPER METHODS
    // ======================================================================

    /**
     * Формирует путь к снепшоту в хранилище.
     * Формат: snapshots/{studentId}/{testId}/{attemptId}.json
     */
    private String buildSnapshotPath(String studentId, String testId, String attemptId) {
        return String.format("%s%s/%s/%s.json", SNAPSHOT_PREFIX, studentId, testId, attemptId);
    }

    /**
     * Формирует путь к изображению в хранилище.
     * Формат: {entityType}/{entityId}/{imageId}.{ext}
     */
    private String buildImagePath(String entityType, String entityId, String imageId) {
        String prefix = entityType.equals("question") ? QUESTION_PREFIX : ANSWER_PREFIX;
        return String.format("%s%s/%s", prefix, entityId, imageId);
    }
}