package services

import (
	"fmt"

	"github.com/google/uuid"
)

// toString normalizes DB values to string (handles []byte/uuid and string).
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
