package apiconst


const (
	ParamCategoryID = "category_id"
	ParamCourseID   = "course_id"
	ParamLessonID   = "lesson_id"
)

const (
	PathCategory = "/:" + ParamCategoryID 
	PathCourse   = "/:" + ParamCourseID   
	PathLesson   = "/:" + ParamLessonID   
)
