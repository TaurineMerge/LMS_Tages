package models

type Category struct {
	BaseModel
	Title string `json:"title"`
}

// CategoryCreate - DTO для создания категории
type CategoryCreate struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

// CategoryUpdate - DTO для обновления категории
type CategoryUpdate struct {
	Title string `json:"title" validate:"required,min=1,max=255"`
}

// CategoryResponse - ответ с категорией
type CategoryResponse struct {
	Category
}

// CategoryListResponse - список категорий
type CategoryListResponse struct {
	Data  []CategoryResponse `json:"data"`
	Total int                `json:"total"`
}
