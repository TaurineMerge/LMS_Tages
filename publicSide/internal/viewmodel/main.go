// Package viewmodel содержит структуры, которые используются для передачи данных в шаблоны (views).
package viewmodel

// MainViewModel представляет основные данные для корневого шаблона `layouts/main.hbs`.
type MainViewModel struct {
	Title string // Заголовок страницы, который будет отображаться в теге <title>.
}

// NewMain создает новую модель представления для основного макета.
func NewMain(title string) *MainViewModel {
	return &MainViewModel{
		Title: title,
	}
}
