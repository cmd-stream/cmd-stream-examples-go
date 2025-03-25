package hw

import "errors"

func ValidateLength(length int) (err error) {
	if length > 10 {
		return errors.New("length is too large")
	}
	return
}
