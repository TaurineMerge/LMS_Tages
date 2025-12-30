package com.example.lms.shared.controller;

import java.io.InputStream;
import java.util.UUID;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import com.example.lms.shared.storage.StorageServiceInterface;
import com.example.lms.shared.storage.exceptions.StorageException;
import com.example.lms.shared.storage.dto.UploadResult;

import io.javalin.http.Context;
import io.javalin.http.UploadedFile;

/**
 * Контроллер для работы с изображениями вопросов и ответов.
 */
public class ImageController {

    private static final Logger logger = LoggerFactory.getLogger(ImageController.class);
    private final StorageServiceInterface storageService;

    public ImageController(StorageServiceInterface storageService) {
        this.storageService = storageService;
    }

    /**
     * POST /api/questions/{questionId}/images
     * Загрузка изображения для вопроса.
     */
    public void uploadQuestionImage(Context ctx) {
        String questionId = ctx.pathParam("questionId");
        UploadedFile file = ctx.uploadedFile("image");

        if (file == null) {
            ctx.status(400).json(new ErrorResponse("Файл не найден"));
            return;
        }

        try {
            String imageId = UUID.randomUUID().toString() +
                    storageService.getExtensionFromContentType(file.contentType());

            UploadResult result = storageService.uploadImage(
                    "question",
                    questionId,
                    imageId,
                    file.content(),
                    file.contentType(),
                    file.size());

            // Генерируем presigned URL для доступа к изображению
            String imageUrl = storageService.getImageUrl("question", questionId, imageId);

            ctx.status(201).json(new ImageUploadResponse(
                    imageId,
                    imageUrl,
                    result.getSize(),
                    result.getContentType()));

            logger.info("Изображение загружено для вопроса {}: {}", questionId, imageId);

        } catch (StorageException e) {
            logger.error("Ошибка загрузки изображения", e);
            ctx.status(400).json(new ErrorResponse(e.getMessage()));
        } catch (Exception e) {
            logger.error("Неожиданная ошибка", e);
            ctx.status(500).json(new ErrorResponse("Внутренняя ошибка сервера"));
        }
    }

    /**
     * POST /api/answers/{answerId}/images
     * Загрузка изображения для ответа.
     */
    public void uploadAnswerImage(Context ctx) {
        String answerId = ctx.pathParam("answerId");
        UploadedFile file = ctx.uploadedFile("image");

        if (file == null) {
            ctx.status(400).json(new ErrorResponse("Файл не найден"));
            return;
        }

        try {
            String imageId = UUID.randomUUID().toString() +
                    storageService.getExtensionFromContentType(file.contentType());

            UploadResult result = storageService.uploadImage(
                    "answer",
                    answerId,
                    imageId,
                    file.content(),
                    file.contentType(),
                    file.size());

            String imageUrl = storageService.getImageUrl("answer", answerId, imageId);

            ctx.status(201).json(new ImageUploadResponse(
                    imageId,
                    imageUrl,
                    result.getSize(),
                    result.getContentType()));

            logger.info("Изображение загружено для ответа {}: {}", answerId, imageId);

        } catch (StorageException e) {
            logger.error("Ошибка загрузки изображения", e);
            ctx.status(400).json(new ErrorResponse(e.getMessage()));
        } catch (Exception e) {
            logger.error("Неожиданная ошибка", e);
            ctx.status(500).json(new ErrorResponse("Внутренняя ошибка сервера"));
        }
    }

    /**
     * GET /api/questions/{questionId}/images/{imageId}
     * Получение ссылки на изображение вопроса.
     */
    public void getQuestionImage(Context ctx) {
        String questionId = ctx.pathParam("questionId");
        String imageId = ctx.pathParam("imageId");

        try {
            String imageUrl = storageService.getImageUrl("question", questionId, imageId);
            ctx.json(new ImageUrlResponse(imageUrl));
        } catch (Exception e) {
            logger.error("Ошибка получения URL изображения", e);
            ctx.status(404).json(new ErrorResponse("Изображение не найдено"));
        }
    }

    /**
     * GET /api/answers/{answerId}/images/{imageId}
     * Получение ссылки на изображение ответа.
     */
    public void getAnswerImage(Context ctx) {
        String answerId = ctx.pathParam("answerId");
        String imageId = ctx.pathParam("imageId");

        try {
            String imageUrl = storageService.getImageUrl("answer", answerId, imageId);
            ctx.json(new ImageUrlResponse(imageUrl));
        } catch (Exception e) {
            logger.error("Ошибка получения URL изображения", e);
            ctx.status(404).json(new ErrorResponse("Изображение не найдено"));
        }
    }

    /**
     * DELETE /api/questions/{questionId}/images/{imageId}
     * Удаление изображения вопроса.
     */
    public void deleteQuestionImage(Context ctx) {
        String questionId = ctx.pathParam("questionId");
        String imageId = ctx.pathParam("imageId");

        boolean deleted = storageService.deleteImage("question", questionId, imageId);

        if (deleted) {
            ctx.status(204);
            logger.info("Изображение удалено: question={}, image={}", questionId, imageId);
        } else {
            ctx.status(404).json(new ErrorResponse("Изображение не найдено"));
        }
    }

    /**
     * DELETE /api/answers/{answerId}/images/{imageId}
     * Удаление изображения ответа.
     */
    public void deleteAnswerImage(Context ctx) {
        String answerId = ctx.pathParam("answerId");
        String imageId = ctx.pathParam("imageId");

        boolean deleted = storageService.deleteImage("answer", answerId, imageId);

        if (deleted) {
            ctx.status(204);
            logger.info("Изображение удалено: answer={}, image={}", answerId, imageId);
        } else {
            ctx.status(404).json(new ErrorResponse("Изображение не найдено"));
        }
    }

    // ======================================================================
    // RESPONSE DTOs
    // ======================================================================

    private static class ImageUploadResponse {
        public final String imageId;
        public final String imageUrl;
        public final long size;
        public final String contentType;

        public ImageUploadResponse(String imageId, String imageUrl, long size, String contentType) {
            this.imageId = imageId;
            this.imageUrl = imageUrl;
            this.size = size;
            this.contentType = contentType;
        }
    }

    private static class ImageUrlResponse {
        public final String imageUrl;

        public ImageUrlResponse(String imageUrl) {
            this.imageUrl = imageUrl;
        }
    }

    private static class ErrorResponse {
        public final String error;

        public ErrorResponse(String error) {
            this.error = error;
        }
    }
}