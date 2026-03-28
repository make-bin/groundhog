// @AI_GENERATED
package bcode

// BCode represents a business error code with HTTP status mapping.
type BCode struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
	Details    string `json:"details,omitempty"`
}

// WithDetails returns a copy of the BCode with the given details attached.
func (b BCode) WithDetails(details string) BCode {
	b.Details = details
	return b
}

// Predefined business error codes.
var (
	Success             = BCode{Code: 0, Message: "success", HTTPStatus: 200}
	ErrInvalidRequest   = BCode{Code: 1000, Message: "invalid request", HTTPStatus: 400}
	ErrValidationFailed = BCode{Code: 1001, Message: "validation failed", HTTPStatus: 400}
	ErrUnauthorized     = BCode{Code: 2000, Message: "unauthorized", HTTPStatus: 401}
	ErrTokenExpired     = BCode{Code: 2002, Message: "token expired", HTTPStatus: 401}
	ErrTokenInvalid     = BCode{Code: 2003, Message: "token invalid", HTTPStatus: 401}
	ErrForbidden        = BCode{Code: 3000, Message: "forbidden", HTTPStatus: 403}
	ErrNotFound         = BCode{Code: 4000, Message: "resource not found", HTTPStatus: 404}
	ErrConflict         = BCode{Code: 4002, Message: "resource conflict", HTTPStatus: 409}
	ErrInternal         = BCode{Code: 5000, Message: "internal server error", HTTPStatus: 500}
	ErrDatabase         = BCode{Code: 5001, Message: "database error", HTTPStatus: 500}
	ErrConfig           = BCode{Code: 5003, Message: "configuration error", HTTPStatus: 500}
)

// @AI_GENERATED: end
