package query

type Error struct {
	err error
}

func NewError(err error) error {
	return Error{err: err}
}

func (e Error) Error() string {
	return e.err.Error()
}
