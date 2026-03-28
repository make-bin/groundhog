// @AI_GENERATED
package vo

// ValidationError is a value object representing a configuration validation error.
// It is immutable after creation.
type ValidationError struct {
	path    string
	message string
}

// NewValidationError creates a new ValidationError.
func NewValidationError(path, message string) ValidationError {
	return ValidationError{path: path, message: message}
}

// Path returns the configuration path where the error occurred.
func (e ValidationError) Path() string { return e.path }

// Message returns the validation error message.
func (e ValidationError) Message() string { return e.message }

// @AI_GENERATED: end
