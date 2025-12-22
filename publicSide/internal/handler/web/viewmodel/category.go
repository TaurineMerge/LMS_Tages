package viewmodel

type CategoryViewModel struct {
	Title string
	Ref string
	CoursesAmount int
	Courses []CourseViewModel
	CoursesRef string
}