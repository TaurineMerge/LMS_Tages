package com.example.lms.shared.storage.validator;

import java.io.IOException;
import java.io.InputStream;
import java.util.Arrays;
import java.util.List;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.shared.storage.exceptions.StorageException;

/**
 * Валидатор для изображений.
 * <p>
 * Проверяет MIME-тип, размер и формат изображений перед загрузкой в хранилище.
 */
public class ImageValidator {

    private static final Logger logger = LoggerFactory.getLogger(ImageValidator.class);

    // Допустимые MIME-типы
    private static final List<String> ALLOWED_MIME_TYPES = Arrays.asList(
            "image/jpeg",
            "image/png",
            "image/svg+xml");

    // Допустимые расширения файлов
    private static final List<String> ALLOWED_EXTENSIONS = Arrays.asList(
            ".jpg", ".jpeg", ".png", ".svg");

    // Максимальный размер изображения (5 МБ)
    private static final long MAX_FILE_SIZE = 5 * 1024 * 1024;

    // Минимальный размер изображения (1 КБ)
    private static final long MIN_FILE_SIZE = 1024;

    // Magic bytes для определения формата
    private static final byte[] JPEG_MAGIC = new byte[] { (byte) 0xFF, (byte) 0xD8, (byte) 0xFF };
    private static final byte[] PNG_MAGIC = new byte[] { (byte) 0x89, 0x50, 0x4E, 0x47 };
    private static final byte[] SVG_MAGIC = "<svg".getBytes(); // Упрощенная проверка

    /**
     * Валидирует изображение перед загрузкой.
     * 
     * @param inputStream поток данных изображения
     * @param contentType MIME-тип
     * @param size        размер файла
     * @throws StorageException если валидация не прошла
     */
    public void validate(InputStream inputStream, String contentType, long size) {
        validateContentType(contentType);
        validateSize(size);
        validateMagicBytes(inputStream, contentType);
    }

    /**
     * Валидирует MIME-тип изображения.
     */
    private void validateContentType(String contentType) {
        if (contentType == null || contentType.isEmpty()) {
            throw new StorageException("Content-Type не указан");
        }

        // Нормализуем MIME-тип (убираем параметры типа charset)
        String normalizedType = contentType.split(";")[0].trim().toLowerCase();

        if (!ALLOWED_MIME_TYPES.contains(normalizedType)) {
            logger.warn("Недопустимый MIME-тип: {}", contentType);
            throw new StorageException(
                    String.format("Недопустимый тип файла: %s. Разрешены: %s",
                            contentType, String.join(", ", ALLOWED_MIME_TYPES)));
        }
    }

    /**
     * Валидирует размер файла.
     */
    private void validateSize(long size) {
        if (size < MIN_FILE_SIZE) {
            throw new StorageException(
                    String.format("Файл слишком мал: %d байт (минимум: %d байт)",
                            size, MIN_FILE_SIZE));
        }

        if (size > MAX_FILE_SIZE) {
            throw new StorageException(
                    String.format("Файл слишком велик: %d байт (максимум: %d байт)",
                            size, MAX_FILE_SIZE));
        }
    }

    /**
     * Проверяет magic bytes для определения реального формата файла.
     * <p>
     * Защита от подмены расширения файла.
     */
    private void validateMagicBytes(InputStream inputStream, String contentType) {
        try {
            // Помечаем позицию для возврата
            inputStream.mark(20);

            byte[] header = new byte[20];
            int bytesRead = inputStream.read(header);

            // Возвращаемся в начало потока
            inputStream.reset();

            if (bytesRead < 4) {
                throw new StorageException("Файл слишком мал для определения формата");
            }

            String normalizedType = contentType.split(";")[0].trim().toLowerCase();

            switch (normalizedType) {
                case "image/jpeg":
                    if (!startsWith(header, JPEG_MAGIC)) {
                        throw new StorageException("Файл не является валидным JPEG");
                    }
                    break;

                case "image/png":
                    if (!startsWith(header, PNG_MAGIC)) {
                        throw new StorageException("Файл не является валидным PNG");
                    }
                    break;

                case "image/svg+xml":
                    // SVG - текстовый формат, проверяем наличие XML или <svg
                    String headerStr = new String(header, 0, Math.min(bytesRead, 20));
                    if (!headerStr.contains("<svg") && !headerStr.contains("<?xml")) {
                        throw new StorageException("Файл не является валидным SVG");
                    }
                    break;

                default:
                    throw new StorageException("Неизвестный тип файла: " + contentType);
            }

            logger.debug("Валидация magic bytes пройдена для типа: {}", contentType);

        } catch (IOException e) {
            logger.error("Ошибка при чтении magic bytes", e);
            throw new StorageException("Не удалось прочитать заголовок файла", e);
        }
    }

    /**
     * Проверяет, начинается ли массив байт с заданного префикса.
     */
    private boolean startsWith(byte[] array, byte[] prefix) {
        if (array.length < prefix.length) {
            return false;
        }
        for (int i = 0; i < prefix.length; i++) {
            if (array[i] != prefix[i]) {
                return false;
            }
        }
        return true;
    }

    /**
     * Валидирует имя файла изображения.
     * 
     * @param filename имя файла
     * @return true если имя валидно
     */
    public boolean isValidFilename(String filename) {
        if (filename == null || filename.isEmpty()) {
            return false;
        }

        String lowerFilename = filename.toLowerCase();
        return ALLOWED_EXTENSIONS.stream().anyMatch(lowerFilename::endsWith);
    }

    /**
     * Извлекает расширение из MIME-типа.
     * 
     * @param contentType MIME-тип
     * @return расширение файла (например, ".jpg")
     */
    public String getExtensionFromContentType(String contentType) {
        String normalizedType = contentType.split(";")[0].trim().toLowerCase();

        switch (normalizedType) {
            case "image/jpeg":
                return ".jpg";
            case "image/png":
                return ".png";
            case "image/svg+xml":
                return ".svg";
            default:
                throw new StorageException("Неизвестный MIME-тип: " + contentType);
        }
    }
}