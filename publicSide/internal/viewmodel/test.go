// Package viewmodel содержит структуры, которые используются для передачи данных в шаблоны (views).
package viewmodel

import (
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/domain"
)

// TestViewModel представляет данные для отображения информации о тесте на странице курса.
type TestViewModel struct {
	Title       string
	Ref         string // URL для перехода к прохождению теста.
	Description string
	MinPoint    int
}

// NewTestViewModel создает новую модель представления для теста.
// Возвращает nil, если переданный доменный объект теста равен nil.
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
