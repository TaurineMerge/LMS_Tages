package models

// Category представляет категорию в системе.
// Встраивает BaseModel и содержит поле Title для названия категории.
type Category struct {
	BaseModel
	Title string `json:"title"`
}
