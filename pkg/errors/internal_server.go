package errors

type InternalServer struct {
	description string
}

func NewInternalServer(description string) InternalServer {
	return InternalServer{description: description}
}

func (i InternalServer) Error() string {
	return i.description
}
