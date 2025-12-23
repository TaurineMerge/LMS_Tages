package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"adminPanel/config"
	"adminPanel/exceptions"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tracer = otel.Tracer("adminPanel/services")

// S3Service предоставляет методы для работы с S3-совместимым хранилищем (MinIO)
type S3Service struct {
	client    *minio.Client
	bucket    string
	useSSL    bool
	publicURL string
}

// NewS3Service создает новый сервис для работы с S3
func NewS3Service(cfg config.MinioConfig) (*S3Service, error) {
	// Инициализируем MinIO клиент
	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	return &S3Service{
		client:    minioClient,
		bucket:    cfg.Bucket,
		useSSL:    cfg.UseSSL,
		publicURL: cfg.PublicURL,
	}, nil
}

// EnsureBucketExists проверяет существование бакета и создает его при необходимости
func (s *S3Service) EnsureBucketExists(ctx context.Context) error {
	ctx, span := tracer.Start(ctx, "S3Service.EnsureBucketExists")
	defer span.End()

	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{})
		if err != nil {
			span.RecordError(err)
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		span.AddEvent("bucket created", trace.WithAttributes(
			attribute.String("bucket", s.bucket),
		))
	}

	// Устанавливаем политику публичного чтения для изображений
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, s.bucket)

	err = s.client.SetBucketPolicy(ctx, s.bucket, policy)
	if err != nil {
		span.RecordError(err)
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}

	return nil
}

// UploadImage загружает изображение в S3 и возвращает его URL
func (s *S3Service) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	ctx, span := tracer.Start(ctx, "S3Service.UploadImage")
	defer span.End()

	span.SetAttributes(
		attribute.String("file.name", file.Filename),
		attribute.Int64("file.size", file.Size),
	)

	// Проверяем тип файла
	contentType := file.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		return "", exceptions.NewAppError(
			fmt.Sprintf("Invalid image type: %s. Only JPEG, PNG, GIF, and WEBP are allowed", contentType),
			400,
			"INVALID_IMAGE_TYPE",
		)
	}

	// Проверяем размер файла (максимум 10 МБ)
	maxSize := int64(10 * 1024 * 1024) // 10 MB
	if file.Size > maxSize {
		return "", exceptions.NewAppError(
			fmt.Sprintf("Image size exceeds maximum allowed size of %d bytes", maxSize),
			400,
			"IMAGE_TOO_LARGE",
		)
	}

	// Открываем файл
	src, err := file.Open()
	if err != nil {
		span.RecordError(err)
		return "", exceptions.NewAppError(
			fmt.Sprintf("Failed to open uploaded file: %v", err),
			500,
			"FILE_OPEN_ERROR",
		)
	}
	defer src.Close()

	// Генерируем уникальное имя файла
	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("go/%s/%s%s",
		time.Now().Format("2006/01/02"),
		uuid.New().String(),
		ext,
	)

	span.SetAttributes(attribute.String("object.name", objectName))

	// Загружаем файл в MinIO
	_, err = s.client.PutObject(ctx, s.bucket, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		span.RecordError(err)
		return "", exceptions.NewAppError(
			fmt.Sprintf("Failed to upload image to S3: %v", err),
			500,
			"S3_UPLOAD_ERROR",
		)
	}

	// Формируем публичный URL
	imageURL := s.GetImageURL(objectName)

	span.AddEvent("image uploaded", trace.WithAttributes(
		attribute.String("image.url", imageURL),
	))

	return imageURL, nil
}

// DeleteImage удаляет изображение из S3
func (s *S3Service) DeleteImage(ctx context.Context, imageURL string) error {
	ctx, span := tracer.Start(ctx, "S3Service.DeleteImage")
	defer span.End()

	span.SetAttributes(attribute.String("image.url", imageURL))

	// Извлекаем имя объекта из URL
	objectName := s.extractObjectNameFromURL(imageURL)
	if objectName == "" {
		return exceptions.NewAppError(
			"Invalid image URL",
			400,
			"INVALID_IMAGE_URL",
		)
	}

	span.SetAttributes(attribute.String("object.name", objectName))

	// Удаляем объект из MinIO
	err := s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		span.RecordError(err)
		return exceptions.NewAppError(
			fmt.Sprintf("Failed to delete image from S3: %v", err),
			500,
			"S3_DELETE_ERROR",
		)
	}

	span.AddEvent("image deleted")

	return nil
}

// GetImageURL возвращает публичный URL для изображения
func (s *S3Service) GetImageURL(objectName string) string {
	// Используем публичный URL для доступа через nginx/прокси
	return fmt.Sprintf("%s/%s/%s", strings.TrimRight(s.publicURL, "/"), s.bucket, objectName)
}

// extractObjectNameFromURL извлекает имя объекта из публичного URL
func (s *S3Service) extractObjectNameFromURL(imageURL string) string {
	// Ожидаемый формат: http://endpoint/bucket/objectName
	parts := strings.SplitN(imageURL, "/"+s.bucket+"/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// isValidImageType проверяет, является ли тип контента допустимым для изображения
func isValidImageType(contentType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/gif",
		"image/webp",
	}

	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

// UploadImageFromReader загружает изображение из io.Reader
func (s *S3Service) UploadImageFromReader(ctx context.Context, reader io.Reader, filename string, size int64, contentType string) (string, error) {
	ctx, span := tracer.Start(ctx, "S3Service.UploadImageFromReader")
	defer span.End()

	span.SetAttributes(
		attribute.String("file.name", filename),
		attribute.Int64("file.size", size),
		attribute.String("content.type", contentType),
	)

	// Проверяем тип файла
	if !isValidImageType(contentType) {
		return "", exceptions.NewAppError(
			fmt.Sprintf("Invalid image type: %s. Only JPEG, PNG, GIF, and WEBP are allowed", contentType),
			400,
			"INVALID_IMAGE_TYPE",
		)
	}

	// Генерируем уникальное имя файла
	ext := filepath.Ext(filename)
	objectName := fmt.Sprintf("go/%s/%s%s",
		time.Now().Format("2006/01/02"),
		uuid.New().String(),
		ext,
	)

	span.SetAttributes(attribute.String("object.name", objectName))

	// Загружаем файл в MinIO
	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		span.RecordError(err)
		return "", exceptions.NewAppError(
			fmt.Sprintf("Failed to upload image to S3: %v", err),
			500,
			"S3_UPLOAD_ERROR",
		)
	}

	// Формируем публичный URL
	imageURL := s.GetImageURL(objectName)

	span.AddEvent("image uploaded", trace.WithAttributes(
		attribute.String("image.url", imageURL),
	))

	return imageURL, nil
}
