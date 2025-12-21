package viewmodel

// CategoryCourseView представляет краткую информацию о курсе для страницы категорий.
type CategoryCourseView struct {
	ID    string
	Title string
}

// CategoryView представляет категорию с ограниченным списком курсов для отображения.
type CategoryView struct {
	ID             string
	Title          string
	TotalCourses   int
	Courses        []CategoryCourseView
	HasMoreCourses bool
}

// CategoriesPageViewModel содержит все данные, необходимые для отображения страницы
// со списком всех категорий.
type CategoriesPageViewModel struct {
	PageHeader PageHeaderViewModel
	Categories []CategoryView
}
