package main

// Централизованные коды ошибок API.
// Эти значения используются в ErrorResponse / ErrorResponseDTO.

const (
	ErrCodeBadRequest       = "BAD_REQUEST"
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeNotFound         = "NOT_FOUND"
	ErrCodeConflict         = "CONFLICT"
	ErrCodeInternal         = "INTERNAL_ERROR"
	ErrCodeUnauthorized     = "UNAUTHORIZED"
)

// Утилиты для создания DTO ошибок с единообразными кодами.

func newBadRequestError(message string) ErrorResponseDTO {
	return toErrorResponseDTO(ErrorResponse{
		Error: message,
		Code:  ErrCodeBadRequest,
	})
}

func newValidationError(message string) ErrorResponseDTO {
	return toErrorResponseDTO(ErrorResponse{
		Error: message,
		Code:  ErrCodeValidation,
	})
}

func newNotFoundError(message string) ErrorResponseDTO {
	return toErrorResponseDTO(ErrorResponse{
		Error: message,
		Code:  ErrCodeNotFound,
	})
}

func newConflictError(message string) ErrorResponseDTO {
	return toErrorResponseDTO(ErrorResponse{
		Error: message,
		Code:  ErrCodeConflict,
	})
}

func newInternalError(message string) ErrorResponseDTO {
	return toErrorResponseDTO(ErrorResponse{
		Error: message,
		Code:  ErrCodeInternal,
	})
}

func newUnauthorizedError(message string) ErrorResponseDTO {
	return toErrorResponseDTO(ErrorResponse{
		Error: message,
		Code:  ErrCodeUnauthorized,
	})
}
