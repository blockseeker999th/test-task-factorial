package utils

type ErrorResponse struct {
	Error string `json:"error"`
}

var (
	ErrFailedToDecode   = "failed to decode request"
	ErrValidationFailed = "incorrect input"
	ErrSavingCalc       = "error saving calculation"
)

func NewErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Error: message,
	}
}
