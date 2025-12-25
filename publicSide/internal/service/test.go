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

// TestService defines the interface for test-related business logic.
type TestService interface {
	GetTest(ctx context.Context, categoryID, courseID string) (*domain.Test, error)
}

type testService struct {
	testingClient *testing.Client
}

// NewTestService creates a new instance of the test service.
func NewTestService(testingClient *testing.Client) TestService {
	return &testService{
		testingClient: testingClient,
	}
}

// GetTest retrieves test details, calling the client and mapping the result to a domain model.
func (s *testService) GetTest(ctx context.Context, categoryID, courseID string) (*domain.Test, error) {
	tracer := otel.Tracer("service")
	ctx, span := tracer.Start(ctx, "testService.GetTest")
	defer span.End()

	span.SetAttributes(
		attribute.String("category_id", categoryID),
		attribute.String("course_id", courseID),
	)

	// Call the client to get the test DTO
	testDTO, err := s.testingClient.GetTest(ctx, categoryID, courseID)
	if err != nil {
		if errors.Is(err, testing.ErrTestNotFound) {
			return nil, apperrors.NewNotFound("Test")
		}
		if errors.Is(err, testing.ErrServiceUnavailable) {
			slog.Error("Testing service is unavailable", "error", err)
			return nil, apperrors.NewServiceUnavailable("Testing")
		}
		// For other errors like invalid response, log and return a generic error
		slog.Error("Failed to get test from client", "error", err)
		return nil, err
	}

	// Map DTO to domain model
	domainTest := mapTestDataToDomain(testDTO)

	return domainTest, nil
}

// mapTestDataToDomain converts a testing.TestData DTO to a domain.Test model.
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
