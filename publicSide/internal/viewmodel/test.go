package viewmodel

import (
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
func NewTestViewModel(test *domain.Test, testURL string) *TestViewModel {
	if test == nil {
		return nil
	}

	return &TestViewModel{
		Title:       test.Title,
		Description: test.Description,
		MinPoint:    test.MinPoint,
		Ref:         testURL,
	}
}
