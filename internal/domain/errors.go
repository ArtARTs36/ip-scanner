package domain

import (
	"errors"
	"fmt"
)

var ErrNotFound = errors.New("not found")

type InvalidIPError struct {
	err error
}

func NewErrIPInvalid(err error) *InvalidIPError {
	return &InvalidIPError{err: err}
}

func (e InvalidIPError) Error() string {
	return fmt.Sprintf("ip invalid: %s", e.err.Error())
}
