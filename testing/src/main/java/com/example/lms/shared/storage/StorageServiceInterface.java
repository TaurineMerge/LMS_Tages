package com.example.lms.shared.storage;

import java.io.InputStream;
import java.time.LocalDate;
import java.util.List;
import java.util.Optional;

import com.example.lms.shared.storage.dto.FileMetadata;
import com.example.lms.shared.storage.dto.UploadResult;

/**
 * Интерфейс для работы с объектным хранилищем (MinIO).
 * <p>
 * Предоставляет методы для загрузки, скачивания и управления файлами
 * в объектном хранилище. Поддерживает работу со снепшотами попыток
 * тестирования и изображениями из вопросов/ответов.
 * 
 * <h2>Структура хранилища:</h2>
 * 
 * <pre>
 * bucket: test-attempts/
 *   └── snapshots/
 *       └── {student_id}/
 *           └── {test_id}/
 *               └── {attempt_id}.json
 * 
 * bucket: test-images/
 *   └── questions/
 *       └── {question_id}/
 *           └── {image_id}.{ext}
 *   └── answers/
 *       └── {answer_id}/
 *           └── {image_id}.{ext}
 * </pre>
 */
public interface StorageServiceInterface {

        /**
         * Загружает файл в хранилище.
         * 
         * @param bucketName  название bucket
         * @param objectPath  путь к объекту (включая имя файла)
         * @param inputStream поток данных файла
         * @param contentType MIME-тип файла (например, "application/json", "image/png")
         * @param size        размер файла в байтах
         * @return результат загрузки с метаданными
         * @throws com.example.lms.shared.storage.dto.StorageException если загрузка не
         *                                                             удалась
         */
        UploadResult upload(String bucketName, String objectPath, InputStream inputStream,
                        String contentType, long size);

        /**
         * Загружает снепшот попытки тестирования.
         * <p>
         * Автоматически формирует путь: snapshots/{studentId}/{testId}/{attemptId}.json
         * 
         * @param studentId    ID студента
         * @param testId       ID теста
         * @param attemptId    ID попытки
         * @param snapshotJson JSON-снепшот в виде строки
         * @return результат загрузки
         */
        UploadResult uploadSnapshot(String studentId, String testId, String attemptId, String snapshot,
                        String attemptVersionJson, LocalDate attemptDate);

        /**
         * Загружает изображение для вопроса или ответа.
         * 
         * @param entityType  тип сущности ("question" или "answer")
         * @param entityId    ID сущности (вопроса или ответа)
         * @param imageId     ID изображения
         * @param inputStream поток данных изображения
         * @param contentType MIME-тип изображения
         * @param size        размер файла в байтах
         * @return результат загрузки
         */
        UploadResult uploadImage(String entityType, String entityId, String imageId,
                        InputStream inputStream, String contentType, long size);

        /**
         * Скачивает файл из хранилища.
         * 
         * @param bucketName название bucket
         * @param objectPath путь к объекту
         * @return поток данных файла, или Optional.empty() если файл не найден
         */
        Optional<InputStream> download(String bucketName, String objectPath);

        /**
         * Скачивает снепшот попытки тестирования.
         * 
         * @param studentId ID студента
         * @param testId    ID теста
         * @param attemptId ID попытки
         * @return JSON-снепшот в виде строки, или Optional.empty() если не найден
         */
        Optional<String> downloadSnapshot(String studentId, String testId, String attemptId);

        /**
         * Скачивает изображение.
         * 
         * @param entityType тип сущности ("question" или "answer")
         * @param entityId   ID сущности
         * @param imageId    ID изображения
         * @return поток данных изображения, или Optional.empty() если не найдено
         */
        Optional<InputStream> downloadImage(String entityType, String entityId, String imageId);

        /**
         * Получает метаданные файла без скачивания.
         * 
         * @param bucketName название bucket
         * @param objectPath путь к объекту
         * @return метаданные файла, или Optional.empty() если файл не найден
         */
        Optional<FileMetadata> getMetadata(String bucketName, String objectPath);

        /**
         * Удаляет файл из хранилища.
         * 
         * @param bucketName название bucket
         * @param objectPath путь к объекту
         * @return true если файл был удален, false если файл не существовал
         */
        boolean delete(String bucketName, String objectPath);

        /**
         * Удаляет снепшот попытки.
         * 
         * @param studentId ID студента
         * @param testId    ID теста
         * @param attemptId ID попытки
         * @return true если снепшот был удален
         */
        boolean deleteSnapshot(String studentId, String testId, String attemptId);

        /**
         * Удаляет изображение.
         * 
         * @param entityType тип сущности ("question" или "answer")
         * @param entityId   ID сущности
         * @param imageId    ID изображения
         * @return true если изображение было удалено
         */
        boolean deleteImage(String entityType, String entityId, String imageId);

        /**
         * Получает список объектов в bucket с заданным префиксом.
         * 
         * @param bucketName название bucket
         * @param prefix     префикс пути (например, "snapshots/student-123/")
         * @return список путей к объектам
         */
        List<String> listObjects(String bucketName, String prefix);

        /**
         * Получает список всех снепшотов студента для конкретного теста.
         * 
         * @param studentId ID студента
         * @param testId    ID теста
         * @return список ID попыток
         */
        List<String> listSnapshots(String studentId, String testId);

        /**
         * Проверяет существование файла в хранилище.
         * 
         * @param bucketName название bucket
         * @param objectPath путь к объекту
         * @return true если файл существует
         */
        boolean exists(String bucketName, String objectPath);

        /**
         * Создает presigned URL для прямого доступа к файлу.
         * <p>
         * URL действителен в течение указанного времени и позволяет
         * скачать файл без дополнительной аутентификации.
         * 
         * @param bucketName       название bucket
         * @param objectPath       путь к объекту
         * @param expiresInSeconds время жизни URL в секундах
         * @return presigned URL для скачивания файла
         */
        String getPresignedUrl(String bucketName, String objectPath, int expiresInSeconds);

        /**
         * Создает presigned URL для изображения (по умолчанию на 1 час).
         * 
         * @param entityType тип сущности ("question" или "answer")
         * @param entityId   ID сущности
         * @param imageId    ID изображения
         * @return presigned URL
         */
        String getImageUrl(String entityType, String entityId, String imageId);

        /**
         * Получает расширение из Content-Type.
         * 
         * @param contentType Content-Type
         * @return extension
         */
        String getExtensionFromContentType(String contentType);
}