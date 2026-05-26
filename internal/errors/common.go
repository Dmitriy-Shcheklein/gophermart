package errors

import "errors"

const (
	InvalidContentTypeMsg = "invalid content-type"
	DecodeBodyErrMsg      = "error while decode body"
	ValidateBodyErrMsg    = "error while validate body"
)

var (
	ErrEmptyDep                = errors.New("empty dependency")
	ErrLoginDuplicate          = errors.New("login already exists")
	ErrInvalidAuthData         = errors.New("invalid auth data")
	ErrOrderAlreadyExists      = errors.New("order already exists")
	ErrOrderBelongsAnotherUser = errors.New("order belongs to another user")
	ErrOrderInvalidNumber      = errors.New("invalid order number")
)
