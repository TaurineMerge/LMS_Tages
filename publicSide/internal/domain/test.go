// Package domain определяет основные бизнес-сущности и модели данных,
// которые используются во всем приложении.
package domain

// Test представляет собой итоговый тест по курсу.
type Test struct {
	ID          string // Уникальный идентификатор
	CourseID    string // ID курса, к которому относится тест
	Title       string // Название теста
	MinPoint    int    // Минимальный балл для прохождения
	Description string // Описание теста
}
