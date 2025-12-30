// Package testing предоставляет клиент для взаимодействия с внешним сервисом тестирования.
package testing

// TestResponse представляет собой стандартную обертку ответа от сервиса тестирования.
type TestResponse struct {
	Data   *TestData `json:"data"`   // Полезная нагрузка (данные о тесте) или null.
	Status string    `json:"status"` // Статус ответа ("success", "not_found", и т.д.).
}

// TestData содержит детальную информацию о тесте.
type TestData struct {
	ID          string `json:"id"`          // Уникальный идентификатор теста.
	CourseID    string `json:"courseId"`    // ID курса, к которому относится тест.
	Title       string `json:"title"`       // Название теста.
	MinPoint    int    `json:"min_point"`   // Минимальный балл для прохождения.
	Description string `json:"description"` // Описание теста.
}
