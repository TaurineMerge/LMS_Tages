package middleware

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"adminPanel/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ValidateStruct - валидирует структуру и возвращает ошибки
func ValidateStruct(dto interface{}) (map[string]string, error) {
	errors := make(map[string]string)
	val := reflect.ValueOf(dto).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		validateTag := structField.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		rules := strings.Split(validateTag, ",")
		fieldName := getJSONFieldName(structField)

		for _, rule := range rules {
			if err := applyValidationRule(field, rule, fieldName); err != nil {
				errors[fieldName] = err.Error()
				break
			}
		}
	}

	return errors, nil
}

// ValidateRequest - middleware для валидации запроса
func ValidateRequest(dto interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Создаем экземпляр DTO
		dtoType := reflect.TypeOf(dto)
		if dtoType.Kind() == reflect.Ptr {
			dtoType = dtoType.Elem()
		}
		dtoValue := reflect.New(dtoType).Interface()

		// Парсим тело запроса
		if err := c.BodyParser(dtoValue); err != nil {
			return c.Status(400).JSON(models.ErrorResponse{
				Error: "Invalid request body",
				Code:  "BAD_REQUEST",
			})
		}

		// Валидируем DTO
		if validationErrors, err := ValidateStruct(dtoValue); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{
				Error: "Validation error",
				Code:  "INTERNAL_ERROR",
			})
		} else if len(validationErrors) > 0 {
			return c.Status(422).JSON(models.ValidationErrorResponse{
				Error:  "Validation failed",
				Code:   "VALIDATION_ERROR",
				Errors: validationErrors,
			})
		}

		// Сохраняем валидированные данные в контекст
		c.Locals("validatedDTO", dtoValue)

		return c.Next()
	}
}

// ValidateMiddleware - middleware для валидации DTO
func ValidateMiddleware(dto interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Создаем экземпляр DTO
		dtoType := reflect.TypeOf(dto)
		if dtoType.Kind() == reflect.Ptr {
			dtoType = dtoType.Elem()
		}
		dtoValue := reflect.New(dtoType).Interface()

		// Парсим тело запроса
		if err := c.BodyParser(dtoValue); err != nil {
			return c.Status(400).JSON(models.ErrorResponse{
				Error: "Invalid request body",
				Code:  "BAD_REQUEST",
			})
		}

		// Валидируем DTO
		if validationErrors := validateStruct(dtoValue); len(validationErrors) > 0 {
			return c.Status(422).JSON(models.ValidationErrorResponse{
				Error:  "Validation failed",
				Code:   "VALIDATION_ERROR",
				Errors: validationErrors,
			})
		}

		// Сохраняем валидированные данные в контекст
		c.Locals("validatedDTO", dtoValue)

		return c.Next()
	}
}

// validateStruct выполняет валидацию структуры
func validateStruct(dto interface{}) map[string]string {
	errors := make(map[string]string)
	val := reflect.ValueOf(dto).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		// Получаем тег validate
		validateTag := structField.Tag.Get("validate")
		if validateTag == "" {
			continue
		}

		// Парсим правила валидации
		rules := strings.Split(validateTag, ",")
		fieldName := getJSONFieldName(structField)

		for _, rule := range rules {
			if err := applyValidationRule(field, rule, fieldName); err != nil {
				errors[fieldName] = err.Error()
				break
			}
		}
	}

	return errors
}

// applyValidationRule применяет правило валидации
func applyValidationRule(field reflect.Value, rule, fieldName string) error {
	switch {
	case rule == "required":
		return validateRequired(field, fieldName)
	case strings.HasPrefix(rule, "min="):
		min := parseRuleValue(rule)
		return validateMin(field, fieldName, min)
	case strings.HasPrefix(rule, "max="):
		max := parseRuleValue(rule)
		return validateMax(field, fieldName, max)
	case rule == "uuid4":
		return validateUUID(field, fieldName)
	case strings.HasPrefix(rule, "oneof="):
		options := strings.Split(strings.TrimPrefix(rule, "oneof="), " ")
		return validateOneOf(field, fieldName, options)
	case rule == "email":
		return validateEmail(field, fieldName)
	default:
		return nil
	}
}

// Вспомогательные функции валидации
func validateRequired(field reflect.Value, fieldName string) error {
	if field.IsZero() {
		return fiber.NewError(400, fieldName+" is required")
	}
	return nil
}

func validateMin(field reflect.Value, fieldName string, min int) error {
	switch field.Kind() {
	case reflect.String:
		if field.Len() < min {
			return fiber.NewError(400, fieldName+" must be at least "+strconv.Itoa(min)+" characters")
		}
	}
	return nil
}

func validateMax(field reflect.Value, fieldName string, max int) error {
	switch field.Kind() {
	case reflect.String:
		if field.Len() > max {
			return fiber.NewError(400, fieldName+" must be at most "+strconv.Itoa(max)+" characters")
		}
	}
	return nil
}

func validateUUID(field reflect.Value, fieldName string) error {
	if !field.IsZero() {
		if _, err := uuid.Parse(field.String()); err != nil {
			return fiber.NewError(400, fieldName+" must be a valid UUID")
		}
	}
	return nil
}

func validateOneOf(field reflect.Value, fieldName string, options []string) error {
	if !field.IsZero() {
		value := field.String()
		for _, option := range options {
			if value == option {
				return nil
			}
		}
		return fiber.NewError(400, fieldName+" must be one of: "+strings.Join(options, ", "))
	}
	return nil
}

func validateEmail(field reflect.Value, fieldName string) error {
	if !field.IsZero() {
		emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
		matched, _ := regexp.MatchString(emailRegex, field.String())
		if !matched {
			return fiber.NewError(400, fieldName+" must be a valid email")
		}
	}
	return nil
}

func parseRuleValue(rule string) int {
	parts := strings.Split(rule, "=")
	if len(parts) == 2 {
		var value int
		if _, err := fmt.Sscanf(parts[1], "%d", &value); err == nil {
			return value
		}
	}
	return 0
}

func getJSONFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return field.Name
	}

	// Убираем опции (например, "omitempty")
	parts := strings.Split(jsonTag, ",")
	return parts[0]
}
