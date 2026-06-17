package errors

type Conflict struct {
	description string
}

func NewConflict(description string) Conflict {
	return Conflict{description: description}
}

func (e Conflict) Error() string {
	return e.description
}
