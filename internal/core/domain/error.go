package domain

import "errors"

var (
	ErrInternal        = errors.New("internal error")
	ErrConflictingData = errors.New("data conflicts with existing data in unique column")
	ErrDataNotFound    = errors.New("data not found")
	ErrNoUpdatedData   = errors.New("no data to updated")
	ErrInvalidData     = errors.New("invalid data")
)
