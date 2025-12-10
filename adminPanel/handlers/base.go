// Package handlers contains HTTP route handlers.
package handlers

import (
	"github.com/google/uuid"
)

// isValidUUID проверяет валидность UUID
func isValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
