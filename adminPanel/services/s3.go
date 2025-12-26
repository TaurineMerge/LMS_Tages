package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"adminPanel/config"
	"adminPanel/middleware"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// tracer трассировщик для сервиса S3.
// Используется для отслеживания операций с MinIO/S3.
var tracer = otel.Tracer("adminPanel/services")

// S3Service предоставляет методы для работы с MinIO/S3 хранилищем.
// Позволяет загружать, удалять и получать URL изображений.
type S3Service struct {
	client    *minio.Client
	bucket    string
	useSSL    bool
	publicURL string
}

// NewS3Service создает новый экземпляр S3Service на основе конфигурации MinIO.
// Инициализирует клиента MinIO и возвращает сервис.
func NewS3Service(cfg config.MinioConfig) (*S3Service, error) {
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

// EnsureBucketExists проверяет существование bucket и создает его, если необходимо.
// Устанавливает публичную политику доступа для чтения объектов.
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

// UploadImage загружает изображение из multipart.FileHeader в S3.
// Проверяет тип и размер файла, генерирует уникальное имя и возвращает публичный URL.
func (s *S3Service) UploadImage(ctx context.Context, file *multipart.FileHeader) (string, error) {
	ctx, span := tracer.Start(ctx, "S3Service.UploadImage")
	defer span.End()

	span.SetAttributes(
		attribute.String("file.name", file.Filename),
		attribute.Int64("file.size", file.Size),
	)

	contentType := file.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		return "", middleware.NewAppError(
			fmt.Sprintf("Invalid image type: %s. Only JPEG, PNG, GIF, and WEBP are allowed", contentType),
			400,
			"INVALID_IMAGE_TYPE",
		)
	}

	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		return "", middleware.NewAppError(
			fmt.Sprintf("Image size exceeds maximum allowed size of %d bytes", maxSize),
			400,
			"IMAGE_TOO_LARGE",
		)
	}

	src, err := file.Open()
	if err != nil {
		span.RecordError(err)
		return "", middleware.NewAppError(
			fmt.Sprintf("Failed to open uploaded file: %v", err),
			500,
			"FILE_OPEN_ERROR",
		)
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("go/%s/%s%s",
		time.Now().Format("2006/01/02"),
		uuid.New().String(),
		ext,
	)

	span.SetAttributes(attribute.String("object.name", objectName))

	_, err = s.client.PutObject(ctx, s.bucket, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		span.RecordError(err)
		return "", middleware.NewAppError(
			fmt.Sprintf("Failed to upload image to S3: %v", err),
			500,
			"S3_UPLOAD_ERROR",
		)
	}

	imageURL := s.GetImageURL(objectName)

	span.AddEvent("image uploaded", trace.WithAttributes(
		attribute.String("image.url", imageURL),
	))

	return imageURL, nil
}

// UploadImageKey загружает изображение из multipart.FileHeader в S3.
// Проверяет тип и размер файла, генерирует уникальное имя и возвращает ключ объекта.
func (s *S3Service) UploadImageKey(ctx context.Context, file *multipart.FileHeader) (string, error) {
	ctx, span := tracer.Start(ctx, "S3Service.UploadImageKey")
	defer span.End()

	span.SetAttributes(
		attribute.String("file.name", file.Filename),
		attribute.Int64("file.size", file.Size),
	)

	contentType := file.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		return "", middleware.NewAppError(
			fmt.Sprintf("Invalid image type: %s. Only JPEG, PNG, GIF, and WEBP are allowed", contentType),
			400,
			"INVALID_IMAGE_TYPE",
		)
	}

	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		return "", middleware.NewAppError(
			fmt.Sprintf("Image size exceeds maximum allowed size of %d bytes", maxSize),
			400,
			"IMAGE_TOO_LARGE",
		)
	}

	src, err := file.Open()
	if err != nil {
		span.RecordError(err)
		return "", middleware.NewAppError(
			fmt.Sprintf("Failed to open uploaded file: %v", err),
			500,
			"FILE_OPEN_ERROR",
		)
	}
	defer src.Close()

	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("go/%s/%s%s",
		time.Now().Format("2006/01/02"),
		uuid.New().String(),
		ext,
	)

	span.SetAttributes(attribute.String("object.name", objectName))

	_, err = s.client.PutObject(ctx, s.bucket, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		span.RecordError(err)
		return "", middleware.NewAppError(
			fmt.Sprintf("Failed to upload image to S3: %v", err),
			500,
			"S3_UPLOAD_ERROR",
		)
	}

	span.AddEvent("image uploaded", trace.WithAttributes(
		attribute.String("object.key", objectName),
	))

	return objectName, nil
}

// DeleteImage удаляет изображение из S3 по публичному URL.
// Извлекает имя объекта из URL и удаляет его.
func (s *S3Service) DeleteImage(ctx context.Context, imageURL string) error {
	ctx, span := tracer.Start(ctx, "S3Service.DeleteImage")
	defer span.End()

	span.SetAttributes(attribute.String("image.url", imageURL))

	objectName := s.extractObjectNameFromURL(imageURL)
	if objectName == "" {
		return middleware.NewAppError(
			"Invalid image URL",
			400,
			"INVALID_IMAGE_URL",
		)
	}

	span.SetAttributes(attribute.String("object.name", objectName))

	err := s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		span.RecordError(err)
		return middleware.NewAppError(
			fmt.Sprintf("Failed to delete image from S3: %v", err),
			500,
			"S3_DELETE_ERROR",
		)
	}

	span.AddEvent("image deleted")

	return nil
}

// GetImageURL формирует публичный URL для объекта по его имени.
// Использует publicURL, bucket и objectName.
func (s *S3Service) GetImageURL(objectName string) string {
	return fmt.Sprintf("%s/%s/%s", strings.TrimRight(s.publicURL, "/"), s.bucket, objectName)
}

// extractObjectNameFromURL извлекает имя объекта из публичного URL.
// Разбирает URL и возвращает часть после bucket.
func (s *S3Service) extractObjectNameFromURL(imageURL string) string {
	parts := strings.SplitN(imageURL, "/"+s.bucket+"/", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// isValidImageType проверяет, является ли contentType допустимым типом изображения.
// Поддерживает JPEG, PNG, GIF, WEBP.
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

// UploadImageFromReader загружает изображение из io.Reader в S3.
// Принимает reader, имя файла, размер и тип контента, возвращает публичный URL.
func (s *S3Service) UploadImageFromReader(ctx context.Context, reader io.Reader, filename string, size int64, contentType string) (string, error) {
	ctx, span := tracer.Start(ctx, "S3Service.UploadImageFromReader")
	defer span.End()

	span.SetAttributes(
		attribute.String("file.name", filename),
		attribute.Int64("file.size", size),
		attribute.String("content.type", contentType),
	)

	if !isValidImageType(contentType) {
		return "", middleware.NewAppError(
			fmt.Sprintf("Invalid image type: %s. Only JPEG, PNG, GIF, and WEBP are allowed", contentType),
			400,
			"INVALID_IMAGE_TYPE",
		)
	}

	ext := filepath.Ext(filename)
	objectName := fmt.Sprintf("go/%s/%s%s",
		time.Now().Format("2006/01/02"),
		uuid.New().String(),
		ext,
	)

	span.SetAttributes(attribute.String("object.name", objectName))

	_, err := s.client.PutObject(ctx, s.bucket, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		span.RecordError(err)
		return "", middleware.NewAppError(
			fmt.Sprintf("Failed to upload image to S3: %v", err),
			500,
			"S3_UPLOAD_ERROR",
		)
	}

	imageURL := s.GetImageURL(objectName)

	span.AddEvent("image uploaded", trace.WithAttributes(
		attribute.String("image.url", imageURL),
	))

	return imageURL, nil
}

// UploadImageFromURL скачивает изображение по URL и загружает в S3.
// Проверяет тип контента, генерирует имя и возвращает публичный URL.
func (s *S3Service) UploadImageFromURL(ctx context.Context, imageURL string) (string, error) {
	ctx, span := tracer.Start(ctx, "S3Service.UploadImageFromURL")
	defer span.End()

	span.SetAttributes(attribute.String("source.url", imageURL))

	resp, err := http.Get(imageURL)
	if err != nil {
		span.RecordError(err)
		return "", middleware.NewAppError(
			fmt.Sprintf("Failed to download image from URL: %v", err),
			400,
			"IMAGE_DOWNLOAD_ERROR",
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", middleware.NewAppError(
			fmt.Sprintf("Failed to download image: HTTP %d", resp.StatusCode),
			400,
			"IMAGE_DOWNLOAD_ERROR",
		)
	}

	contentType := resp.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		return "", middleware.NewAppError(
			fmt.Sprintf("Invalid image type from URL: %s", contentType),
			400,
			"INVALID_IMAGE_TYPE",
		)
	}

	ext := filepath.Ext(imageURL)
	if ext == "" || len(ext) > 5 {
		switch contentType {
		case "image/jpeg", "image/jpg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		default:
			ext = ".jpg"
		}
	}

	objectName := fmt.Sprintf("go/%s/%s%s",
		time.Now().Format("2006/01/02"),
		uuid.New().String(),
		ext,
	)

	span.SetAttributes(attribute.String("object.name", objectName))

	_, err = s.client.PutObject(ctx, s.bucket, objectName, resp.Body, resp.ContentLength, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		span.RecordError(err)
		return "", middleware.NewAppError(
			fmt.Sprintf("Failed to upload image to S3: %v", err),
			500,
			"S3_UPLOAD_ERROR",
		)
	}

	s3URL := s.GetImageURL(objectName)

	span.AddEvent("image uploaded from URL", trace.WithAttributes(
		attribute.String("s3.url", s3URL),
	))

	return s3URL, nil
}
