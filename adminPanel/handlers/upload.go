package handlers

import (
	"fmt"

	"adminPanel/middleware"
	"adminPanel/services"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UploadHandler - HTTP обработчик для загрузки файлов
type UploadHandler struct {
	s3Service *services.S3Service
}

// NewUploadHandler создает новый HTTP обработчик для загрузки файлов
func NewUploadHandler(s3Service *services.S3Service) *UploadHandler {
	return &UploadHandler{
		s3Service: s3Service,
	}
}

// RegisterRoutes регистрирует маршруты для загрузки файлов
func (h *UploadHandler) RegisterRoutes(upload fiber.Router) {
	upload.Post("/image", h.uploadImage)
	upload.Post("/image-from-url", h.uploadImageFromURL)
}

// UploadImageResponse - структура ответа при успешной загрузке изображения
type UploadImageResponse struct {
	Status   string `json:"status"`
	ImageURL string `json:"image_url"`
	Message  string `json:"message"`
}

// uploadImage обрабатывает POST /api/v1/upload/image
// @Summary Загрузить изображение
// @Description Загружает изображение в S3 хранилище и возвращает публичный URL
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param image formData file true "Файл изображения (JPEG, PNG, GIF, WEBP, максимум 10 МБ)"
// @Success 200 {object} UploadImageResponse
// @Failure 400 {object} middleware.AppError "Неверный тип файла или размер"
// @Failure 500 {object} middleware.AppError "Ошибка загрузки"
// @Router /api/v1/upload/image [post]
func (h *UploadHandler) uploadImage(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span := trace.SpanFromContext(ctx)

	// Получаем файл из запроса
	file, err := c.FormFile("image")
	if err != nil {
		return middleware.NewAppError(
			fmt.Sprintf("Failed to read uploaded file: %v", err),
			400,
			"MISSING_FILE",
		)
	}

	span.SetAttributes(
		attribute.String("file.name", file.Filename),
		attribute.Int64("file.size", file.Size),
	)

	// Загружаем изображение в S3
	imageURL, err := h.s3Service.UploadImage(ctx, file)
	if err != nil {
		return err
	}

	span.AddEvent("image uploaded successfully", trace.WithAttributes(
		attribute.String("image.url", imageURL),
	))

	return c.JSON(UploadImageResponse{
		Status:   "success",
		ImageURL: imageURL,
		Message:  "Image uploaded successfully",
	})
}

// UploadImageFromURLRequest - структура запроса для загрузки изображения по URL
type UploadImageFromURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}

// uploadImageFromURL обрабатывает POST /api/v1/upload/image-from-url
// @Summary Загрузить изображение по URL
// @Description Скачивает изображение по URL и загружает в S3 хранилище
// @Tags Upload
// @Accept json
// @Produce json
// @Param body body UploadImageFromURLRequest true "URL изображения"
// @Success 200 {object} UploadImageResponse
// @Failure 400 {object} middleware.AppError "Неверный URL или тип файла"
// @Failure 500 {object} middleware.AppError "Ошибка загрузки"
// @Router /api/v1/upload/image-from-url [post]
func (h *UploadHandler) uploadImageFromURL(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span := trace.SpanFromContext(ctx)

	var req UploadImageFromURLRequest
	if err := c.BodyParser(&req); err != nil {
		return middleware.NewAppError(
			fmt.Sprintf("Invalid request body: %v", err),
			400,
			"VALIDATION_ERROR",
		)
	}

	if req.URL == "" {
		return middleware.NewAppError(
			"URL is required",
			400,
			"MISSING_URL",
		)
	}

	span.SetAttributes(attribute.String("source.url", req.URL))

	// Загружаем изображение по URL в S3
	imageURL, err := h.s3Service.UploadImageFromURL(ctx, req.URL)
	if err != nil {
		return err
	}

	span.AddEvent("image uploaded from URL", trace.WithAttributes(
		attribute.String("image.url", imageURL),
	))

	return c.JSON(UploadImageResponse{
		Status:   "success",
		ImageURL: imageURL,
		Message:  "Image uploaded successfully from URL",
	})
}
