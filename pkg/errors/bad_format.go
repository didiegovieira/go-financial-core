package errors

type BadFormat struct {
	description string
}

func NewBadFormat(description string) BadFormat {
	return BadFormat{description: description}
}

func (e BadFormat) Error() string {
	return e.description
}
