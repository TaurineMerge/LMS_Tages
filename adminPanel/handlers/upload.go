// Пакет handlers содержит обработчики HTTP-запросов для различных операций,
// включая загрузку файлов и управление категориями, курсами, уроками и здоровьем системы.
package handlers

import (
	"fmt"

	"adminPanel/middleware"
	"adminPanel/services"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// UploadHandler обрабатывает запросы на загрузку изображений в S3-совместимое хранилище.
type UploadHandler struct {
	s3Service *services.S3Service
}

// NewUploadHandler создает новый экземпляр UploadHandler с заданным сервисом S3.
func NewUploadHandler(s3Service *services.S3Service) *UploadHandler {
	return &UploadHandler{
		s3Service: s3Service,
	}
}

// RegisterRoutes регистрирует маршруты для загрузки изображений на переданном роутере.
func (h *UploadHandler) RegisterRoutes(upload fiber.Router) {
	upload.Post("/image", h.uploadImage)
	upload.Post("/image-from-url", h.uploadImageFromURL)
}

// UploadImageResponse представляет ответ на запрос загрузки изображения.
type UploadImageResponse struct {
	Status   string `json:"status"`
	ImageURL string `json:"image_url"`
	Message  string `json:"message"`
}

// uploadImage обрабатывает POST /upload/image.
// Загружает изображение из multipart формы в S3-совместимое хранилище.
func (h *UploadHandler) uploadImage(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span := trace.SpanFromContext(ctx)

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

// UploadImageFromURLRequest представляет запрос на загрузку изображения по URL.
type UploadImageFromURLRequest struct {
	URL string `json:"url" validate:"required,url"`
}

// uploadImageFromURL обрабатывает POST /upload/image-from-url.
// Загружает изображение по указанному URL в S3-совместимое хранилище.
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
