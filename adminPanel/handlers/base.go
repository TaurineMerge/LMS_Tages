// Package handlers содержит HTTP-обработчики для всех маршрутов приложения.
// Этот пакет предоставляет обработчики для health check, категорий, курсов и уроков.
package handlers

import (
	"github.com/google/uuid"
)

// isValidUUID проверяет валидность UUID строки
//
// Функция используется для валидации параметров маршрутов,
// содержащих идентификаторы сущностей.
//
// Параметры:
//   - u: строка для проверки
//
// Возвращает:
//   - bool: true, если строка является валидным UUID
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
