package viewmodel

import (
	"fmt"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/clients/testing"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
)

type TestViewModel struct {
	Title       string
	Ref         string
	Description string
	MinPoint    int
}

// NewTestViewModel creates a new TestViewModel from a domain model.
// It constructs the Ref URL to the external testing service.
func NewTestViewModel(test *domain.Test, testingServiceBaseURL string) *TestViewModel {
	if test == nil {
		return nil
	}

	// Construct the full URL to the test on the testing service's frontend
	// Assuming a path structure like /courses/{courseID}/test
	// This might need adjustment based on the actual frontend routing of the testing service.
	refURL := fmt.Sprintf(testing.TEST_PATH, testingServiceBaseURL, test.CourseID)

	return &TestViewModel{
		Title:       test.Title,
		Description: test.Description,
		MinPoint:    test.MinPoint,
		Ref:         refURL,
	}
}