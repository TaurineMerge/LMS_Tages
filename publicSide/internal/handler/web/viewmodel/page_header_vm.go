package viewmodel

// Breadcrumb представляет один элемент "хлебных крошек".
type Breadcrumb struct {
    Text string
    URL  string // URL может быть пустым для последнего, неактивного элемента
}

// PageHeaderViewModel содержит все данные для partial'а page_header.hbs.
type PageHeaderViewModel struct {
    Title       string
    Breadcrumbs []Breadcrumb
}
