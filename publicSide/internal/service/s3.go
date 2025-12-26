package service

import (
	"fmt"
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

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

// GetImageURL возвращает публичный URL для изображения
func (s *S3Service) GetImageURL(objectName string) string {
	if objectName == "" {
		return ""
	}
	// Используем публичный URL для доступа через nginx/прокси
	return fmt.Sprintf("%s/%s/%s", strings.TrimRight(s.publicURL, "/"), s.bucket, objectName)
}
