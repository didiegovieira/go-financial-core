package errors

type NotFound struct {
	description string
}

func NewNotFound(description string) NotFound {
	return NotFound{description: description}
}

func (e NotFound) Error() string {
	return e.description
}

func TagNotFound(message string, availableTags any) map[string]interface{} {
	return map[string]interface{}{
		"message":        message,
		"available_tags": availableTags,
	}
}
