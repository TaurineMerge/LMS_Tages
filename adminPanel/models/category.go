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
	Title string `json:"title" validate:"omitempty,min=1,max=255"`
}

// CategoryResponse - ответ с категорией
type CategoryResponse struct {
	Status string   `json:"status"`
	Data   Category `json:"data"`
}

// PaginatedCategoriesResponse - пагинированный список категорий
type PaginatedCategoriesResponse struct {
	Status string `json:"status"`
	Data   struct {
		Items      []Category `json:"items"`
		Pagination Pagination `json:"pagination"`
	} `json:"data"`
}
