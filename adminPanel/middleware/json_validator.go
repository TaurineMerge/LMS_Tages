// Package middleware содержит middleware-функции для обработки HTTP-запросов.
// Этот файл реализует валидацию с использованием JSON Schema.
package middleware

import (
	"adminPanel/handlers/dto/response"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/santhosh-tekuri/jsonschema/v5"
)

// Примечание: embed не поддерживает относительные пути с ..
// Схемы должны быть в поддиректории текущего пакета или использовать абсолютный путь
// Вместо этого будем загружать схемы из файловой системы напрямую

// SchemaValidator управляет компиляцией и кешированием JSON Schema.
type SchemaValidator struct {
	compiler *jsonschema.Compiler
	schemas  map[string]*jsonschema.Schema
	mu       sync.RWMutex
}

var (
	validator     *SchemaValidator
	validatorOnce sync.Once
)

// GetValidator возвращает синглтон экземпляр SchemaValidator.
func GetValidator() *SchemaValidator {
	validatorOnce.Do(func() {
		validator = &SchemaValidator{
			compiler: jsonschema.NewCompiler(),
			schemas:  make(map[string]*jsonschema.Schema),
		}
		validator.compiler.Draft = jsonschema.Draft7

		// Загружаем схемы из embed FS
		if err := validator.loadSchemas(); err != nil {
			panic(fmt.Sprintf("Failed to load schemas: %v", err))
		}
	})
	return validator
}

// loadSchemas загружает все JSON Schema из файловой системы.
func (v *SchemaValidator) loadSchemas() error {
	// Получаем путь к директории проекта
	schemasPath := filepath.Join("docs", "schemas")

	schemaFiles := []string{
		"category_schema.json",
		"category-create.json",
		"category-update.json",
		"course_schema.json",
		"course-create.json",
		"course-update.json",
		"lesson_schema.json",
		"lesson-create.json",
		"lesson-update.json",
	}

	for _, schemaFile := range schemaFiles {
		schemaPath := filepath.Join(schemasPath, schemaFile)

		data, err := os.ReadFile(schemaPath)
		if err != nil {
			return fmt.Errorf("failed to read schema %s: %w", schemaFile, err)
		}

		// Добавляем схему в компилятор
		if err := v.compiler.AddResource(schemaFile, bytes.NewReader(data)); err != nil {
			return fmt.Errorf("failed to add schema %s: %w", schemaFile, err)
		}
	}

	return nil
} // GetSchema возвращает скомпилированную схему по имени.
func (v *SchemaValidator) GetSchema(schemaName string) (*jsonschema.Schema, error) {
	v.mu.RLock()
	schema, exists := v.schemas[schemaName]
	v.mu.RUnlock()

	if exists {
		return schema, nil
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	// Проверяем еще раз после получения блокировки на запись
	if schema, exists := v.schemas[schemaName]; exists {
		return schema, nil
	}

	// Компилируем схему
	compiled, err := v.compiler.Compile(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema %s: %w", schemaName, err)
	}

	v.schemas[schemaName] = compiled
	return compiled, nil
}

// ValidateJSONSchema middleware для валидации запроса с использованием JSON Schema.
//
// Параметры:
//   - schemaName: имя файла схемы для валидации (например, "course_schema.json")
//
// Возвращает:
//   - fiber.Handler: middleware-функцию для использования в Fiber приложении
func ValidateJSONSchema(schemaName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Получаем валидатор
		v := GetValidator()

		// Получаем схему
		schema, err := v.GetSchema(schemaName)
		if err != nil {
			return c.Status(500).JSON(response.ErrorResponse{
				Status: "error",
				Error: response.ErrorDetails{
					Code:    "SCHEMA_ERROR",
					Message: fmt.Sprintf("Schema not found: %s", schemaName),
				},
			})
		}

		// Парсим тело запроса в interface{}
		var requestBody interface{}
		if err := json.Unmarshal(c.Body(), &requestBody); err != nil {
			return c.Status(400).JSON(response.ErrorResponse{
				Status: "error",
				Error: response.ErrorDetails{
					Code:    "INVALID_JSON",
					Message: "Invalid JSON format",
				},
			})
		}

		// Валидируем через JSON Schema
		if err := schema.Validate(requestBody); err != nil {
			validationErrors := make(map[string]string)

			if ve, ok := err.(*jsonschema.ValidationError); ok {
				// Извлекаем детальные ошибки валидации
				validationErrors = extractValidationErrors(ve)
			} else {
				validationErrors["_error"] = err.Error()
			}

			return c.Status(422).JSON(response.ValidationErrorResponse{
				Status: "error",
				Error: response.ErrorDetails{
					Code:    "VALIDATION_ERROR",
					Message: "Validation failed",
				},
				Errors: validationErrors,
			})
		}

		// Сохраняем валидированное тело в Locals
		c.Locals("validatedBody", requestBody)

		return c.Next()
	}
}

// extractValidationErrors извлекает ошибки валидации из ValidationError.
func extractValidationErrors(ve *jsonschema.ValidationError) map[string]string {
	errors := make(map[string]string)

	// Основная ошибка
	if ve.InstanceLocation != "" {
		fieldName := ve.InstanceLocation
		if fieldName == "" {
			fieldName = "root"
		}
		errors[fieldName] = ve.Message
	}

	// Вложенные ошибки
	for _, cause := range ve.Causes {
		fieldErrors := extractValidationErrors(cause)
		for field, msg := range fieldErrors {
			errors[field] = msg
		}
	}

	return errors
}
