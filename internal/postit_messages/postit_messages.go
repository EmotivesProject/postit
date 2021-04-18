package postit_messages

import "errors"

const (
	MsgHealthOK = "Health ok"
)

var (
	ErrFailedDecoding  = errors.New("Failed during decoding request")
	ErrAlreadyLiked    = errors.New("Already liked post")
	ErrInvalid         = errors.New("Failed validation")
	ErrInvalidUsername = errors.New("Username failed")
)
