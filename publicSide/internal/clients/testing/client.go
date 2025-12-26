// Package testing предоставляет клиент для взаимодействия с внешним сервисом тестирования.
package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/santhosh-tekuri/jsonschema/v5"
)

const (
	// TEST_API_PATH - путь для внутреннего API-взаимодействия для получения информации о тесте.
	TEST_API_PATH = "/testing/internal/categories/%s/courses/%s/test"
	// TEST_UI_PATH - путь для пользовательского интерфейса для прохождения теста.
	TEST_UI_PATH = "/testing/categories/%s/courses/%s/test"

	// STATUS_OK - строковый литерал для успешного статуса ответа.
	STATUS_OK = "success"
	// STATUS_NOT_FOUND - строковый литерал для статуса "не найдено".
	STATUS_NOT_FOUND = "not_found"
)

// Client инкапсулирует логику для отправки запросов к сервису тестирования.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	schema     *jsonschema.Schema
}

// NewClient создает новый экземпляр клиента для сервиса тестирования.
// `baseURL` - это базовый URL сервиса (например, "http://localhost:8081").
// `schemaPath` - путь к файлу JSON-схемы для валидации ответов.
func NewClient(baseURL string, schemaPath string) (*Client, error) {
	parsedBaseURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	schema, err := jsonschema.Compile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to compile schema: %w", err)
	}

	return &Client{
		baseURL: parsedBaseURL,
		httpClient: &http.Client{
			Timeout: 3 * time.Second,
		},
		schema: schema,
	}, nil
}

// GetTest запрашивает информацию о тесте для конкретного курса.
// Он выполняет GET-запрос, валидирует ответ по JSON-схеме и разбирает его.
// Возвращает `ErrTestNotFound`, если тест не найден, `ErrServiceUnavailable` при проблемах с сетью
// или `ErrInvalidResponse` при несоответствии ответа схеме.
func (c *Client) GetTest(ctx context.Context, categoryID, courseID string) (*TestData, error) {
	path := fmt.Sprintf(TEST_API_PATH, categoryID, courseID)
	requestURL := c.baseURL.ResolveReference(&url.URL{Path: path})

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrServiceUnavailable, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var v interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response for validation: %w: %v", ErrInvalidResponse, err)
	}
	if err := c.schema.Validate(v); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidResponse, err)
	}

	var testResponse TestResponse
	if err := json.Unmarshal(body, &testResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response into DTO: %w", err)
	}

	if testResponse.Status == STATUS_NOT_FOUND {
		return nil, ErrTestNotFound
	}

	return testResponse.Data, nil
}

// GetUITestURL генерирует полный URL для страницы прохождения теста.
func GetUITestURL(baseURL, categoryId, courseId string) string {
	url := fmt.Sprintf("%s/%s", baseURL, TEST_UI_PATH)
	return fmt.Sprintf(url, categoryId, courseId)
}
