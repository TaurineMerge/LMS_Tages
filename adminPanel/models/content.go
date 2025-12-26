package models

import (
	"encoding/json"
	"fmt"
)

// ContentTypeText константа для типа контента "text".
const (
	ContentTypeText  = "text"
	ContentTypeImage = "image"
)

// Content интерфейс для различных типов контента.
// Определяет метод Type(), возвращающий тип контента.
type Content interface {
	Type() string
}

// TextContent представляет текстовый контент.
// Содержит тип контента и данные текста с валидацией на максимальную длину 10000 символов.
type TextContent struct {
	ContentType string `json:"content_type"`
	Data        string `json:"data" validate:"max=10000"`
}

// Type возвращает тип контента для TextContent.
func (t TextContent) Type() string { return ContentTypeText }

// ImageContent представляет контент с изображением.
// Содержит тип контента, URL изображения и альтернативный текст с валидацией.
type ImageContent struct {
	ContentType string `json:"content_type"`
	URL         string `json:"url" validate:"required,url"`
	Alt         string `json:"alt" validate:"max=255"`
}

// Type возвращает тип контента для ImageContent.
func (i ImageContent) Type() string { return ContentTypeImage }

// ContentSlice представляет срез контентов.
// Реализует кастомный маршалинг и анмаршалинг JSON для поддержки полиморфизма.
type ContentSlice []Content

// UnmarshalJSON реализует кастомный анмаршалинг JSON для ContentSlice.
// Разбирает массив JSON, определяя тип каждого элемента по полю content_type.
func (cs *ContentSlice) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*cs = nil
		return nil
	}

	var rawMessages []json.RawMessage
	if err := json.Unmarshal(data, &rawMessages); err != nil {
		return fmt.Errorf("failed to unmarshal content array: %w", err)
	}

	*cs = make(ContentSlice, 0, len(rawMessages))

	for i, rawMessage := range rawMessages {
		if string(rawMessage) == "null" {
			continue
		}

		var base struct {
			ContentType string `json:"content_type"`
		}
		if err := json.Unmarshal(rawMessage, &base); err != nil {
			return fmt.Errorf("failed to read content_type at index %d: %w", i, err)
		}

		content, err := unmarshalContent(base.ContentType, rawMessage)
		if err != nil {
			return fmt.Errorf("failed to unmarshal content at index %d: %w", i, err)
		}

		*cs = append(*cs, content)
	}

	return nil
}

// unmarshalContent анмаршалит конкретный контент на основе типа.
// Принимает contentType и данные JSON, возвращает соответствующий Content.
func unmarshalContent(contentType string, data []byte) (Content, error) {
	switch contentType {
	case ContentTypeText:
		var content TextContent
		if err := json.Unmarshal(data, &content); err != nil {
			return nil, err
		}
		return content, nil

	case ContentTypeImage:
		var content ImageContent
		if err := json.Unmarshal(data, &content); err != nil {
			return nil, err
		}
		return content, nil

	default:
		return nil, fmt.Errorf("unknown content type: '%s'", contentType)
	}
}

// MarshalJSON реализует кастомный маршалинг JSON для ContentSlice.
// Преобразует срез в массив JSON.
func (cs ContentSlice) MarshalJSON() ([]byte, error) {
	if cs == nil {
		return []byte("null"), nil
	}

	arr := make([]interface{}, len(cs))
	for i, content := range cs {
		arr[i] = content
	}

	return json.Marshal(arr)
}
