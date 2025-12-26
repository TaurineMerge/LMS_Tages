// Package handlers содержит HTTP-обработчики для API adminPanel.
// Включает обработчики для категорий, курсов, уроков, загрузки файлов и проверки здоровья.
package handlers

import (
	"github.com/google/uuid"
)

// isValidUUID проверяет, является ли строка валидным UUID.
// Возвращает true, если строка может быть распарсена как UUID.
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
