package errors

type Code string

const (
	ErrNotFound     Code = "NOT_FOUND"
	ErrDuplicate    Code = "DUPLICATE"
	ErrValidation   Code = "VALIDATION"
	ErrImageProcess Code = "IMAGE_PROCESS"
	ErrStorage      Code = "STORAGE"
	ErrDatabase     Code = "DATABASE"
	ErrInternal     Code = "INTERNAL"
)
