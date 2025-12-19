package models

// Category - модель категории курсов
//
// Представляет категорию, к которой могут принадлежать курсы.
// Категория содержит базовую информацию и наследует поля от BaseModel.
//
// Поля:
//   - BaseModel: встроенная структура с общими полями (ID, CreatedAt, UpdatedAt)
//   - Title: название категории
type Category struct {
	BaseModel
	Title string `json:"title"`
}
