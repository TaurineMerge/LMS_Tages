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

const TEST_PATH = "/testing/internal/categories/%s/courses/%s/test"

// Client is a client for the testing service.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	schema     *jsonschema.Schema
}

// NewClient creates a new client for the testing service.
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

// GetTest retrieves a test from the testing service.
func (c *Client) GetTest(ctx context.Context, categoryID, courseID string) (*TestData, error) {
	path := fmt.Sprintf(TEST_PATH, categoryID, courseID)
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
		return nil, fmt.Errorf("failed to unmarshal response for validation: %w", err)
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
