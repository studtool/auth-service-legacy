package errs

//go:generate easyjson

const (
	Internal      = 0
	BadFormat     = 1
	InvalidFormat = 2
	Conflict      = 3
	NotFound      = 4
)

//easyjson:json
type Error struct {
	Type    int    `json:"-"`
	Message string `json:"message"`
}

func NewError(t int, message string) *Error {
	return &Error{
		Type:    t,
		Message: message,
	}
}

func NewInternalError(message string) *Error {
	return NewError(Internal, message)
}

func NewBadFormatError(message string) *Error {
	return NewError(BadFormat, message)
}

func NewInvalidFormatError(message string) *Error {
	return NewError(InvalidFormat, message)
}

func NewConflictError(message string) *Error {
	return NewError(Conflict, message)
}

func NewNotFoundError(message string) *Error {
	return NewError(NotFound, message)
}
