// Package template provides template engine configuration for the application.
package template

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	"github.com/gofiber/template/handlebars/v2"
)

// NewEngine creates a new Handlebars template engine with custom helpers.
func NewEngine(cfg *config.Config) *handlebars.Engine {
	engine := handlebars.New("./templates", ".hbs")

	// Reload and Debug should be enabled based on development mode
	engine.Reload(cfg.Dev)
	engine.Debug(cfg.Dev)

	// Register custom helper for truncating text
	engine.AddFunc("truncate", func(text string, length int) string {
		if utf8.RuneCountInString(text) <= length {
			return text
		}

		truncated := ""
		count := 0
		for _, r := range text {
			if count >= length {
				break
			}
			truncated += string(r)
			count++
		}

		lastSpace := strings.LastIndex(truncated, " ")
		if lastSpace > 0 {
			truncated = truncated[:lastSpace]
		}

		return truncated + "..."
	})

	engine.AddFunc("eqs", func(a, b string) bool {
		return a == b
	})

	// Register custom helper for equality comparison
	engine.AddFunc("eq", func(a, b int) bool {
		return a == b
	})

	// Register custom helper for string equality comparison
	engine.AddFunc("streq", func(a, b string) bool {
		return a == b
	})

	// Register helper for greater than comparison
	engine.AddFunc("gt", func(a, b int) bool {
		return a > b
	})

	// Register helper for less than comparison
	engine.AddFunc("lt", func(a, b int) bool {
		return a < b
	})

	// Register helper for addition
	engine.AddFunc("add", func(a, b int) int {
		return a + b
	})

	// Register helper for subtraction
	engine.AddFunc("subtract", func(a, b int) int {
		return a - b
	})

	// Register helper for range generation
	engine.AddFunc("range", func(start, end int) []int {
		result := make([]int, 0, end-start+1)
		for i := start; i <= end; i++ {
			result = append(result, i)
		}
		return result
	})

	// Register helper for date formatting
	engine.AddFunc("formatDate", func(t time.Time) string {
		return t.Format("02.01.2006")
	})

	return engine
}
