package viewmodel

type HomePageViewModel struct {
	Categories []CategoryViewModel
}

func NewHomePageViewModel(
	categories []CategoryViewModel,
) *HomePageViewModel {
	return &HomePageViewModel{
		Categories: categories,
	}
}
