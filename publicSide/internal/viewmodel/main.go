package viewmodel

type MainViewModel struct {
	Title       string
}

func NewMain(title string) *MainViewModel {
	return &MainViewModel{
		Title: title,
	}
}