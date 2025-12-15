// Package services содержит бизнес-логику приложения.
// Этот пакет предоставляет сервисы для работы с категориями, курсами и уроками.
package services

import (
	"fmt"

	"github.com/google/uuid"
)

// toString нормализует значения из базы данных в строку.
// Обрабатывает различные типы данных: []byte, uuid.UUID, [16]byte и string.
//
// Особенности:
//   - Если []byte имеет длину 16 байт, пытается преобразовать в UUID
//   - Для [16]byte также пытается преобразовать в UUID
//   - Для uuid.UUID возвращает строковое представление
//   - Для string возвращает значение как есть
//   - Для остальных типов использует fmt.Sprintf("%v", v)
//
// Параметры:
//   - v: интерфейс значения из базы данных
//
// Возвращает:
//   - string: нормализованное строковое представление значения
func toString(v interface{}) string {
	switch val := v.(type) {
	case []byte:
		if len(val) == 16 {
			if u, err := uuid.FromBytes(val); err == nil {
				return u.String()
			}
		}
		return string(val)
	case [16]byte:
		if u, err := uuid.FromBytes(val[:]); err == nil {
			return u.String()
		}
	case uuid.UUID:
		return val.String()
	case string:
		return val
	}
	return fmt.Sprintf("%v", v)
}
