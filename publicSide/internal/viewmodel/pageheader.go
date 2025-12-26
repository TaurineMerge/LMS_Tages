package viewmodel

// PageHeaderViewModel содержит все данные для partial'а page_header.hbs.
type PageHeaderViewModel struct {
	Title       string
	Breadcrumbs []Breadcrumb
}

// NewPageHeaderViewModel создает новый PageHeaderViewModel.
func NewPageHeaderViewModel(title string, breadcrumbs []Breadcrumb) *PageHeaderViewModel {
	return &PageHeaderViewModel{
		Title:       title,
		Breadcrumbs: breadcrumbs,
	}
}
