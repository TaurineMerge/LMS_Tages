package utils

import "strings"

const (
	AscendingDirection  = "ASC"
	DescendingDirection = "DESC"
)

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
