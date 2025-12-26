// Package viewmodel содержит структуры, которые используются для передачи данных в шаблоны (views).
package viewmodel

import "github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"

// PaginationViewModel представляет данные, необходимые для рендеринга компонента пагинации.
type PaginationViewModel struct {
	BaseUrl string // Базовый URL для генерации ссылок на страницы.
	Page    int
	Limit   int
	Total   int
	Pages   int
}

// NewPaginationViewModel создает новую модель представления для пагинации.
func NewPaginationViewModel(pagination response.Pagination, baseUrl string) *PaginationViewModel {
	return &PaginationViewModel{
		BaseUrl: baseUrl,
		Page:    pagination.Page,
		Limit:   pagination.Limit,
		Total:   pagination.Total,
		Pages:   pagination.Pages,
	}
}
