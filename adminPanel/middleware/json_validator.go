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

// SchemaValidator управляет загрузкой и компиляцией JSON-схем для валидации.
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
// Инициализирует валидатор при первом вызове, загружая все схемы.
func GetValidator() *SchemaValidator {
	validatorOnce.Do(func() {
		validator = &SchemaValidator{
			compiler: jsonschema.NewCompiler(),
			schemas:  make(map[string]*jsonschema.Schema),
		}
		validator.compiler.Draft = jsonschema.Draft7

		if err := validator.loadSchemas(); err != nil {
			panic(fmt.Sprintf("Failed to load schemas: %v", err))
		}
	})
	return validator
}

// loadSchemas загружает и добавляет JSON-схемы в компилятор из директории docs/schemas.
func (v *SchemaValidator) loadSchemas() error {
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

		if err := v.compiler.AddResource(schemaFile, bytes.NewReader(data)); err != nil {
			return fmt.Errorf("failed to add schema %s: %w", schemaFile, err)
		}
	}

	return nil
}

// GetSchema возвращает скомпилированную схему по имени.
// Кеширует схемы для повторного использования.
func (v *SchemaValidator) GetSchema(schemaName string) (*jsonschema.Schema, error) {
	v.mu.RLock()
	schema, exists := v.schemas[schemaName]
	v.mu.RUnlock()

	if exists {
		return schema, nil
	}

	v.mu.Lock()
	defer v.mu.Unlock()

	if schema, exists := v.schemas[schemaName]; exists {
		return schema, nil
	}

	compiled, err := v.compiler.Compile(schemaName)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema %s: %w", schemaName, err)
	}

	v.schemas[schemaName] = compiled
	return compiled, nil
}

// ValidateJSONSchema возвращает промежуточное ПО для валидации тела запроса по JSON-схеме.
// В случае ошибки валидации возвращает 422 с деталями ошибок.
func ValidateJSONSchema(schemaName string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		v := GetValidator()

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

		if err := schema.Validate(requestBody); err != nil {
			validationErrors := make(map[string]string)

			if ve, ok := err.(*jsonschema.ValidationError); ok {
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

		c.Locals("validatedBody", requestBody)

		return c.Next()
	}
}

// extractValidationErrors извлекает ошибки валидации из ValidationError в карту полей.
func extractValidationErrors(ve *jsonschema.ValidationError) map[string]string {
	errors := make(map[string]string)

	if ve.InstanceLocation != "" {
		fieldName := ve.InstanceLocation
		if fieldName == "" {
			fieldName = "root"
		}
		errors[fieldName] = ve.Message
	}

	for _, cause := range ve.Causes {
		fieldErrors := extractValidationErrors(cause)
		for field, msg := range fieldErrors {
			errors[field] = msg
		}
	}

	return errors
}
