package errors

type Forbidden struct {
	description string
}

func NewForbidden(description string) Forbidden {
	return Forbidden{description: description}
}

func (e Forbidden) Error() string {
	return e.description
}
