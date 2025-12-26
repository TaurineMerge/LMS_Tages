// Package utils предоставляет общие вспомогательные функции, используемые в разных частях приложения.
package utils

import "strings"

const (
	AscendingDirection  = "ASC"  // Направление сортировки по возрастанию.
	DescendingDirection = "DESC" // Направление сортировки по убыванию.
)

// UnpackSort разбирает строку `sortBy` для определения поля и направления сортировки.
//
// `sortBy`: строка, задающая сортировку. Префикс "-" означает сортировку по убыванию (например, "-created_at").
// `defaultColumn`: поле для сортировки по умолчанию, если `sortBy` пуста или недопустима.
// `defaultDirection`: направление сортировки по умолчанию.
// `allowedColumn`: карта с разрешенными для сортировки полями. Если пуста, проверка не производится.
//
// Возвращает имя колонки и направление сортировки ("ASC" или "DESC").
func UnpackSort(
	sortBy string,
	defaultColumn string,
	defaultDirection string,
	allowedColumn map[string]bool,
) (string, string) {
	if sortBy == "" || (len(allowedColumn) > 0 && !allowedColumn[sortBy]) {
		return defaultColumn, defaultDirection
	}

	direction := AscendingDirection
	if strings.HasPrefix(sortBy, "-") {
		direction = DescendingDirection
		sortBy = strings.TrimPrefix(sortBy, "-")
	}

	return sortBy, direction
}
