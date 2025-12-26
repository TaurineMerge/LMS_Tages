package viewmodel

import "github.com/TaurineMerge/LMS_Tages/publicSide/internal/dto/response"

type PaginationViewModel struct {
	BaseUrl string
	Page    int
	Limit   int
	Total   int
	Pages   int
}

func NewPaginationViewModel(pagination response.Pagination, baseUrl string) *PaginationViewModel{
	return &PaginationViewModel{
		BaseUrl: baseUrl,
		Page: pagination.Page,
		Limit: pagination.Limit,
		Total: pagination.Total,
		Pages: pagination.Pages,
	}
}
