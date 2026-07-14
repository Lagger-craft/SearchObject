package errors

import "fmt"

type AppError struct {
	Code    Code
	Message string
	Op      string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s [%s]: %v", e.Op, e.Message, e.Code, e.Err)
	}
	return fmt.Sprintf("%s: %s [%s]", e.Op, e.Message, e.Code)
}

func (e *AppError) Unwrap() error { return e.Err }

func New(code Code, msg, op string) *AppError {
	return &AppError{Code: code, Message: msg, Op: op}
}

func Wrap(err error, code Code, msg, op string) *AppError {
	return &AppError{Code: code, Message: msg, Op: op, Err: err}
}

func IsCode(err error, code Code) bool {
	if e, ok := err.(*AppError); ok {
		return e.Code == code
	}
	return false
}
