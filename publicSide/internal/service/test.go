// Package service предоставляет бизнес-логику приложения.
package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/clients/testing"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
	"github.com/TaurineMerge/LMS_Tages/publicSide/pkg/apperrors"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// TestService определяет интерфейс для бизнес-логики, связанной с тестами.
type TestService interface {
	// GetTest получает информацию о тесте для указанного курса.
	GetTest(ctx context.Context, categoryID, courseID string) (*domain.Test, error)
}

// testService является реализацией TestService.
type testService struct {
	testingClient *testing.Client
}

// NewTestService создает новый экземпляр testService.
func NewTestService(testingClient *testing.Client) TestService {
	return &testService{
		testingClient: testingClient,
	}
}

// GetTest обращается к клиенту сервиса тестирования для получения данных о тесте.
// Он обрабатывает различные типы ошибок от клиента (не найдено, сервис недоступен)
// и преобразует их в стандартизированные ошибки приложения.
func (s *testService) GetTest(ctx context.Context, categoryID, courseID string) (*domain.Test, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "testService.GetTest")
	defer span.End()

	span.SetAttributes(
		attribute.String("category_id", categoryID),
		attribute.String("course_id", courseID),
	)

	testDTO, err := s.testingClient.GetTest(ctx, categoryID, courseID)
	if err != nil {
		if errors.Is(err, testing.ErrTestNotFound) {
			return nil, apperrors.NewNotFound("Test")
		}
		if errors.Is(err, testing.ErrServiceUnavailable) ||
			errors.Is(err, testing.ErrInvalidResponse) {
			slog.Error("Testing service is unavailable", "error", err)
			return nil, apperrors.NewServiceUnavailable("Testing")
		}

		slog.Error("Failed to get test from client", "error", err)
		return nil, err
	}

	domainTest := mapTestDataToDomain(testDTO)

	return domainTest, nil
}

// mapTestDataToDomain преобразует DTO от клиента в доменную модель Test.
func mapTestDataToDomain(dto *testing.TestData) *domain.Test {
	if dto == nil {
		return nil
	}
	return &domain.Test{
		ID:          dto.ID,
		CourseID:    dto.CourseID,
		Title:       dto.Title,
		MinPoint:    dto.MinPoint,
		Description: dto.Description,
	}
}
