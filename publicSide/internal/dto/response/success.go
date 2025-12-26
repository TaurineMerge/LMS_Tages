// Package response содержит структуры данных для формирования HTTP-ответов.
package response

const (
	StatusSuccess = "success" // Статус для успешного ответа.
	StatusError   = "error"   // Статус для ответа с ошибкой.
)

// SuccessResponse представляет собой стандартную обертку для успешных JSON-ответов.
type SuccessResponse struct {
	Status string      `json:"status"` // Статус ответа, всегда "success".
	Data   interface{} `json:"data"`   // Полезная нагрузка (данные) ответа.
}
