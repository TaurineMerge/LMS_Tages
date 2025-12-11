// Package middleware содержит middleware-функции для обработки HTTP-запросов.
// Этот пакет предоставляет функции для работы с прокси, валидации, аутентификации и обработки ошибок.
package middleware

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"adminPanel/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ValidateStruct валидирует структуру DTO и возвращает ошибки валидации.
//
// Параметры:
//   - dto: интерфейс валидируемой структуры (должен быть указателем на структуру)
//
// Возвращает:
//   - map[string]string: карта ошибок валидации, где ключ - имя поля, значение - сообщение об ошибке
//   - error: ошибка, если произошла ошибка при валидации
func ValidateStruct(dto interface{}) (map[string]string, error) {
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

	return errors, nil
}

// ValidateRequest middleware для валидации запроса.
// Парсит тело запроса в DTO и проверяет его на соответствие правилам валидации.
//
// Параметры:
//   - dto: тип DTO для валидации (обычно используется как ValidateRequest(CreateCategoryDTO{}))
//
// Возвращает:
//   - fiber.Handler: middleware-функцию для использования в Fiber приложении
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
				Status: "error",
				Error: models.ErrorDetails{
					Code:    "VALIDATION_ERROR",
					Message: "Invalid request body",
				},
			})
		}

		// Валидируем DTO
		if validationErrors, err := ValidateStruct(dtoValue); err != nil {
			return c.Status(500).JSON(models.ErrorResponse{
				Status: "error",
				Error: models.ErrorDetails{
					Code:    "SERVER_ERROR",
					Message: "Validation error",
				},
			})
		} else if len(validationErrors) > 0 {
			return c.Status(422).JSON(models.ValidationErrorResponse{
				Status: "error",
				Error: models.ErrorDetails{
					Code:    "VALIDATION_ERROR",
					Message: "Validation failed",
				},
				Errors: validationErrors,
			})
		}

		// Сохраняем валидированные данные в контекст
		c.Locals("validatedDTO", dtoValue)

		return c.Next()
	}
}

// ValidateMiddleware middleware для валидации DTO.
// Аналогичен ValidateRequest, но с более простой реализацией.
//
// Параметры:
//   - dto: тип DTO для валидации
//
// Возвращает:
//   - fiber.Handler: middleware-функцию для использования в Fiber приложении
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
				Status: "error",
				Error: models.ErrorDetails{
					Code:    "VALIDATION_ERROR",
					Message: "Invalid request body",
				},
			})
		}

		// Валидируем DTO
		if validationErrors := validateStruct(dtoValue); len(validationErrors) > 0 {
			return c.Status(422).JSON(models.ValidationErrorResponse{
				Status: "error",
				Error: models.ErrorDetails{
					Code:    "VALIDATION_ERROR",
					Message: "Validation failed",
				},
				Errors: validationErrors,
			})
		}

		// Сохраняем валидированные данные в контекст
		c.Locals("validatedDTO", dtoValue)

		return c.Next()
	}
}

// validateStruct выполняет валидацию структуры.
//
// Параметры:
//   - dto: интерфейс валидируемой структуры
//
// Возвращает:
//   - map[string]string: карта ошибок валидации
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

// applyValidationRule применяет правило валидации к полю.
//
// Параметры:
//   - field: значение поля для валидации
//   - rule: правило валидации
//   - fieldName: имя поля (для сообщений об ошибках)
//
// Возвращает:
//   - error: ошибка валидации или nil, если валидация прошла успешно
func applyValidationRule(field reflect.Value, rule, fieldName string) error {
	switch {
	case rule == "required":
		return validateRequired(field, fieldName)
	case strings.HasPrefix(rule, "min="):
		minValue := parseRuleValue(rule)
		return validateMin(field, fieldName, minValue)
	case strings.HasPrefix(rule, "max="):
		maxValue := parseRuleValue(rule)
		return validateMax(field, fieldName, maxValue)
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

// validateRequired проверяет, что поле обязательно для заполнения.
//
// Параметры:
//   - field: значение поля
//   - fieldName: имя поля для сообщения об ошибке
//
// Возвращает:
//   - error: ошибка, если поле не заполнено
func validateRequired(field reflect.Value, fieldName string) error {
	if field.IsZero() {
		return fiber.NewError(400, fieldName+" is required")
	}
	return nil
}

// validateMin проверяет минимальную длину строки.
//
// Параметры:
//   - field: значение поля
//   - fieldName: имя поля для сообщения об ошибке
//   - minValue: минимально допустимое значение
//
// Возвращает:
//   - error: ошибка, если значение меньше минимального
func validateMin(field reflect.Value, fieldName string, minValue int) error {
	switch field.Kind() {
	case reflect.String:
		if field.Len() < minValue {
			return fiber.NewError(400, fmt.Sprintf("%s must be at least %d characters", fieldName, minValue))
		}
	}
	return nil
}

// validateMax проверяет максимальную длину строки.
//
// Параметры:
//   - field: значение поля
//   - fieldName: имя поля для сообщения об ошибке
//   - maxValue: максимально допустимое значение
//
// Возвращает:
//   - error: ошибка, если значение больше максимального
func validateMax(field reflect.Value, fieldName string, maxValue int) error {
	switch field.Kind() {
	case reflect.String:
		if field.Len() > maxValue {
			return fiber.NewError(400, fmt.Sprintf("%s must be at most %d characters", fieldName, maxValue))
		}
	}
	return nil
}

// validateUUID проверяет, что значение является валидным UUID.
//
// Параметры:
//   - field: значение поля
//   - fieldName: имя поля для сообщения об ошибке
//
// Возвращает:
//   - error: ошибка, если значение не является валидным UUID
func validateUUID(field reflect.Value, fieldName string) error {
	if !field.IsZero() {
		if _, err := uuid.Parse(field.String()); err != nil {
			return fiber.NewError(400, fieldName+" must be a valid UUID")
		}
	}
	return nil
}

// validateOneOf проверяет, что значение поля соответствует одному из допустимых вариантов.
//
// Параметры:
//   - field: значение поля
//   - fieldName: имя поля для сообщения об ошибке
//   - options: список допустимых значений
//
// Возвращает:
//   - error: ошибка, если значение не соответствует ни одному из вариантов
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

// validateEmail проверяет, что значение является валидным email-адресом.
//
// Параметры:
//   - field: значение поля
//   - fieldName: имя поля для сообщения об ошибке
//
// Возвращает:
//   - error: ошибка, если значение не является валидным email
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

// parseRuleValue извлекает числовое значение из правила валидации.
//
// Параметры:
//   - rule: правило валидации в формате "min=5" или "max=10"
//
// Возвращает:
//   - int: извлеченное числовое значение или 0, если извлечение не удалось
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

// getJSONFieldName извлекает имя поля из тега json структуры.
//
// Параметры:
//   - field: информация о поле структуры
//
// Возвращает:
//   - string: имя поля для использования в JSON (без опций вроде "omitempty")
func getJSONFieldName(field reflect.StructField) string {
	jsonTag := field.Tag.Get("json")
	if jsonTag == "" {
		return field.Name
	}

	// Убираем опции (например, "omitempty")
	parts := strings.Split(jsonTag, ",")
	return parts[0]
}
