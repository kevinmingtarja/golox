package scanner

type Error struct {
}

// Error implements the error interface.
func (e Error) Error() string {
	return "error"
}
