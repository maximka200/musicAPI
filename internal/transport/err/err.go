package err

import "errors"

var (
	ErrAlreadyExist = errors.New("already exist")
	ErrNotFound     = errors.New("not found")
)
