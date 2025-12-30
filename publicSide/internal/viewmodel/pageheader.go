// Package viewmodel содержит структуры, которые используются для передачи данных в шаблоны (views).
package viewmodel

// PageHeaderViewModel представляет данные для стандартного заголовка страницы,
// включающего заголовок и "хлебные крошки".
type PageHeaderViewModel struct {
	Title       string
	Breadcrumbs []Breadcrumb
}

// NewPageHeaderViewModel создает новую модель представления для заголовка страницы.
func NewPageHeaderViewModel(title string, breadcrumbs []Breadcrumb) *PageHeaderViewModel {
	return &PageHeaderViewModel{
		Title:       title,
		Breadcrumbs: breadcrumbs,
	}
}
