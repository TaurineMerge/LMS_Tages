package handlers

import (
	"fmt"

	"adminPanel/exceptions"
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
// @Failure 400 {object} exceptions.AppError "Неверный тип файла или размер"
// @Failure 500 {object} exceptions.AppError "Ошибка загрузки"
// @Router /api/v1/upload/image [post]
func (h *UploadHandler) uploadImage(c *fiber.Ctx) error {
	ctx := c.UserContext()
	span := trace.SpanFromContext(ctx)

	// Получаем файл из запроса
	file, err := c.FormFile("image")
	if err != nil {
		return exceptions.NewAppError(
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
