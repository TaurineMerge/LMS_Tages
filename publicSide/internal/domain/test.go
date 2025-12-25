package domain

// Test represents the domain model for a test.
type Test struct {
	ID          string
	CourseID    string
	Title       string
	MinPoint    int
	Description string
}
