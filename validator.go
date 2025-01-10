package exmpls

import "errors"

var ErrTooLargeString = errors.New("too large string")

// StringLengthValidator is used to validate the content of commands.
func StringLengthValidator(length int) (err error) {
	if length > 20 {
		err = ErrTooLargeString
	}
	return
}
