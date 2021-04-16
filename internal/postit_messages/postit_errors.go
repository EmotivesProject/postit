package postit_messages

import "errors"

var (
	ErrFailedDecoding = errors.New("Failed during decoding request")
	ErrAlreadyLiked   = errors.New("Already liked post")
	ErrInvalid        = errors.New("Failed validation")
)
