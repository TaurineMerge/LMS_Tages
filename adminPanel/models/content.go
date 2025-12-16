package models

import (
	"encoding/json"
	"fmt"
)

// Константы для типов контента
const (
	ContentTypeText  = "text"
	ContentTypeImage = "image"
)

// Content - интерфейс для различных типов контента урока.
// Требует от всех типов реализации метода Type().
type Content interface {
	Type() string
}

// TextContent - структура для текстового контента.
type TextContent struct {
	ContentType string `json:"content_type"`
	Data        string `json:"data" validate:"max=10000"`
}

func (t TextContent) Type() string { return ContentTypeText }

// ImageContent - структура для контента-изображения.
type ImageContent struct {
	ContentType string `json:"content_type"`
	URL         string `json:"url" validate:"required,url"`
	Alt         string `json:"alt" validate:"max=255"`
}

func (i ImageContent) Type() string { return ContentTypeImage }

// ContentSlice - кастомный тип для среза интерфейсов Content.
type ContentSlice []Content

// UnmarshalJSON - кастомный десериализатор для ContentSlice.
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

// MarshalJSON - кастомный сериализатор для ContentSlice.
func (cs ContentSlice) MarshalJSON() ([]byte, error) {
	if cs == nil {
		return []byte("null"), nil
	}

	// Просто сериализуем каждый элемент как есть
	arr := make([]interface{}, len(cs))
	for i, content := range cs {
		arr[i] = content
	}

	return json.Marshal(arr)
}
