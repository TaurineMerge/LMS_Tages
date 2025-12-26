package models

// Course представляет курс в системе.
// Встраивает BaseModel и содержит поля для заголовка, описания, уровня сложности,
// ID категории, видимости и ключа изображения.
type Course struct {
	BaseModel
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	CategoryID  string `json:"category_id"`
	Visibility  string `json:"visibility"`
	ImageKey    string `json:"image_key"`
}
