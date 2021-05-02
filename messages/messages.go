package messages

import "errors"

var (
	ErrFailedDecoding  = errors.New("Failed during decoding request")
	ErrAlreadyLiked    = errors.New("Already liked post")
	ErrInvalid         = errors.New("Failed validation")
	ErrInvalidUsername = errors.New("Username failed")
	ErrInvalidCheck    = errors.New("Failed to convert to type")
)
