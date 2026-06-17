package errors

import "fmt"

func NewError(err error, msg string) error {
	return fmt.Errorf(msg+": %w", err)
}

func ErrorResponse(message string, statusCode int) map[string]interface{} {
	return map[string]interface{}{
		"message":     message,
		"status_code": statusCode,
	}
}
