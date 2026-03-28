// @AI_GENERATED
package errors

// ValidationError represents an input validation error.
type ValidationError struct {
	Message string
	Err     error
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) Unwrap() error {
	return e.Err
}

// NotFoundError represents a resource not found error.
type NotFoundError struct {
	Message string
	Err     error
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *NotFoundError) Unwrap() error {
	return e.Err
}

// ConflictError represents a resource conflict error.
type ConflictError struct {
	Message string
	Err     error
}

func (e *ConflictError) Error() string {
	return e.Message
}

func (e *ConflictError) Unwrap() error {
	return e.Err
}

// UnauthorizedError represents an authentication failure error.
type UnauthorizedError struct {
	Message string
	Err     error
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}

func (e *UnauthorizedError) Unwrap() error {
	return e.Err
}

// ForbiddenError represents a permission denied error.
type ForbiddenError struct {
	Message string
	Err     error
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func (e *ForbiddenError) Unwrap() error {
	return e.Err
}

// InternalError represents an unexpected internal error.
type InternalError struct {
	Message string
	Err     error
}

func (e *InternalError) Error() string {
	return e.Message
}

func (e *InternalError) Unwrap() error {
	return e.Err
}

// @AI_GENERATED: end
