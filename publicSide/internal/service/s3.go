// Package service предоставляет бизнес-логику приложения.
package service

import (
	"fmt"
	"strings"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// S3Service инкапсулирует логику для работы с S3-совместимым хранилищем (MinIO).
// На данный момент используется только для получения публичных URL-адресов объектов.
type S3Service struct {
	client    *minio.Client
	bucket    string
	useSSL    bool
	publicURL string
}


// NewS3Service создает новый экземпляр S3Service.
// Он инициализирует клиент MinIO на основе предоставленной конфигурации.
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

// GetImageURL генерирует публичный, общедоступный URL для объекта в хранилище.
// `objectName` — это ключ (имя файла) объекта в бакете.
func (s *S3Service) GetImageURL(objectName string) string {
	if objectName == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s/%s", strings.TrimRight(s.publicURL, "/"), s.bucket, objectName)
}
