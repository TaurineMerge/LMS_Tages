package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// toString преобразует значение в строку.
// Обрабатывает []byte, [16]byte, uuid.UUID, string и другие типы.
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

// parseTime преобразует значение в time.Time.
// Обрабатывает string в формате RFC3339 и time.Time, возвращает zero time при ошибке.
func parseTime(value interface{}) time.Time {
	switch val := value.(type) {
	case string:
		if t, err := time.Parse(time.RFC3339, val); err == nil {
			return t
		}
	case time.Time:
		return val
	}
	return time.Time{}
}
