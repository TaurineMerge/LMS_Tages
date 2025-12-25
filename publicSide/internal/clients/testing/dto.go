package testing

const (
	STATUS_OK = "success"
	STATUS_NOT_FOUND = "not_found"
)
// TestResponse is the top-level structure for the GetTest API response.
type TestResponse struct {
	Data   *TestData `json:"data"`
	Status string    `json:"status"`
}

// TestData contains the detailed information about a test.
type TestData struct {
	ID          string `json:"id"`
	CourseID    string `json:"courseId"`
	Title       string `json:"title"`
	MinPoint    int    `json:"min_point"`
	Description string `json:"description"`
}
