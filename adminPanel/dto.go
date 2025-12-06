package main

import "time"

// DTO структуры повторяют доменные модели 1-в-1.
// Позже сюда можно будет добавить поля/теги, отличающиеся от внутренних моделей.

// HealthResponseDTO описывает ответ health-check.
type HealthResponseDTO struct {
	Status   string `json:"status"`
	Database string `json:"database"`
	Version  string `json:"version"`
}

// ErrorResponseDTO описывает стандартную ошибку API.
type ErrorResponseDTO struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

// CategoryDTO — категория знаний.
type CategoryDTO struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CategoryCreateDTO — тело запроса на создание категории.
type CategoryCreateDTO struct {
	Title string `json:"title"`
}

// CategoryUpdateDTO — тело запроса на обновление категории.
type CategoryUpdateDTO struct {
	Title string `json:"title"`
}

// CourseDTO — курс.
type CourseDTO struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Level       string    `json:"level"`
	CategoryID  string    `json:"category_id"`
	Visibility  string    `json:"visibility"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CourseCreateDTO — тело запроса на создание курса.
type CourseCreateDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	CategoryID  string `json:"category_id"`
	Visibility  string `json:"visibility"`
}

// CourseUpdateDTO — тело запроса на обновление курса.
type CourseUpdateDTO struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Level       string `json:"level"`
	CategoryID  string `json:"category_id"`
	Visibility  string `json:"visibility"`
}

// PaginatedCoursesDTO — страничный ответ по курсам.
type PaginatedCoursesDTO struct {
	Data  []CourseDTO `json:"data"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Pages int         `json:"pages"`
}

// LessonDTO — урок.
type LessonDTO struct {
	ID        string                 `json:"id"`
	Title     string                 `json:"title"`
	CourseID  string                 `json:"course_id"`
	Content   map[string]interface{} `json:"content"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

// LessonCreateDTO — тело запроса на создание урока.
type LessonCreateDTO struct {
	Title   string                 `json:"title"`
	Content map[string]interface{} `json:"content"`
}

// LessonUpdateDTO — тело запроса на обновление урока.
type LessonUpdateDTO struct {
	Title   string                 `json:"title"`
	Content map[string]interface{} `json:"content"`
}

// Мапперы между моделями и DTO. Пока поля 1-в-1, но это изолирует
// внешний контракт от внутренних структур.

func toHealthResponseDTO(h HealthResponse) HealthResponseDTO {
	return HealthResponseDTO(h)
}

func toErrorResponseDTO(e ErrorResponse) ErrorResponseDTO {
	return ErrorResponseDTO(e)
}

func toCategoryDTO(c Category) CategoryDTO {
	return CategoryDTO(c)
}

func toCategoryDTOs(list []Category) []CategoryDTO {
	res := make([]CategoryDTO, 0, len(list))
	for _, c := range list {
		res = append(res, toCategoryDTO(c))
	}
	return res
}

func toCourseDTO(c Course) CourseDTO {
	return CourseDTO(c)
}

func toCourseDTOs(list []Course) []CourseDTO {
	res := make([]CourseDTO, 0, len(list))
	for _, c := range list {
		res = append(res, toCourseDTO(c))
	}
	return res
}

func toPaginatedCoursesDTO(p PaginatedCourses) PaginatedCoursesDTO {
	return PaginatedCoursesDTO{
		Data:  toCourseDTOs(p.Data),
		Total: p.Total,
		Page:  p.Page,
		Limit: p.Limit,
		Pages: p.Pages,
	}
}

func toLessonDTO(l Lesson) LessonDTO {
	return LessonDTO(l)
}

func toLessonDTOs(list []Lesson) []LessonDTO {
	res := make([]LessonDTO, 0, len(list))
	for _, l := range list {
		res = append(res, toLessonDTO(l))
	}
	return res
}
