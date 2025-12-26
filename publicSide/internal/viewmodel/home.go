// Package viewmodel содержит структуры, которые используются для передачи данных в шаблоны (views).
package viewmodel

// HomePageViewModel представляет данные для главной страницы.
type HomePageViewModel struct {
	Categories []CategoryViewModel // Срез категорий для отображения в карусели.
}

// NewHomePageViewModel создает новую модель представления для главной страницы.
func NewHomePageViewModel(
	categories []CategoryViewModel,
) *HomePageViewModel {
	return &HomePageViewModel{
		Categories: categories,
	}
}
