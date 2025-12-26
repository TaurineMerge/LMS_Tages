// Package template отвечает за настройку и конфигурацию шаблонизатора Handlebars.
package template

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/config"
	"github.com/TaurineMerge/LMS_Tages/publicSide/internal/viewmodel"
	"github.com/gofiber/template/handlebars/v2"
)

// NewEngine создает и настраивает новый экземпляр движка шаблонов Handlebars.
// Он настраивает перезагрузку и отладку на основе конфигурации приложения,
// а также регистрирует множество пользовательских вспомогательных функций.
func NewEngine(cfg *config.AppConfig) *handlebars.Engine {
	engine := handlebars.New("./templates", ".hbs")

	// Включает/выключает перезагрузку и отладку шаблонов
	engine.Reload(cfg.Dev)
	engine.Debug(cfg.Dev)

	// truncate обрезает текст до заданной длины, сохраняя целостность слов.
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

	// eq проверяет равенство двух целых чисел.
	engine.AddFunc("eq", func(a, b int) bool {
		return a == b
	})

	// neq проверяет неравенство двух целых чисел.
	engine.AddFunc("neq", func(a, b int) bool {
		return a != b
	})

	// streq проверяет равенство двух строк.
	engine.AddFunc("streq", func(a, b string) bool {
		return a == b
	})

	// gt проверяет, что a > b.
	engine.AddFunc("gt", func(a, b int) bool {
		return a > b
	})

	// lt проверяет, что a < b.
	engine.AddFunc("lt", func(a, b int) bool {
		return a < b
	})

	// add возвращает сумму двух целых чисел.
	engine.AddFunc("add", func(a, b int) int {
		return a + b
	})

	// subtract возвращает разность двух целых чисел.
	engine.AddFunc("subtract", func(a, b int) int {
		return a - b
	})

	// range создает срез целых чисел в заданном диапазоне [start, end].
	engine.AddFunc("range", func(start, end int) []int {
		result := make([]int, 0, end-start+1)
		for i := start; i <= end; i++ {
			result = append(result, i)
		}
		return result
	})

	// len возвращает количество рун (символов) в строке.
	engine.AddFunc("len", func(s string) int {
		return utf8.RuneCountInString(s)
	})

	// formatDate форматирует время в строку "дд.мм.гггг".
	engine.AddFunc("formatDate", func(t time.Time) string {
		return t.Format("02.01.2006")
	})

	// firstLessonRef возвращает URL первого урока из среза, или "#" если срез пуст.
	engine.AddFunc("firstLessonRef", func(slice []viewmodel.LessonViewModel) string {
		if len(slice) == 0 {
			return "#"
		}
		return slice[0].Ref
	})

	return engine
}
