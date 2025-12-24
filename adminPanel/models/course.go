package models

// Course - модель учебного курса
//
// Представляет учебный курс, который может содержать уроки.
// Курс принадлежит определенной категории и имеет уровень сложности.
//
// Поля:
//   - BaseModel: встроенная структура с общими полями (ID, CreatedAt, UpdatedAt)
//   - Title: название курса
//   - Description: описание курса
//   - Level: уровень сложности ("hard", "medium", "easy")
//   - CategoryID: уникальный идентификатор категории, к которой принадлежит курс
//   - Visibility: видимость курса ("draft", "public", "private")
//   - ImageKey: ключ изображения в S3 (опционально)
type Course struct {
	BaseModel
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	CategoryID  string `json:"category_id"`
	Visibility  string `json:"visibility"`
	ImageKey    string `json:"image_key"`
}
