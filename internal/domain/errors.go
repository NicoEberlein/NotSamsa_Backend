package domain

import "errors"

var (
	ErrNotFound         = errors.New("entity not found")
	ErrDuplicateEntity  = errors.New("duplicate entity")
	ErrPermissionDenied = errors.New("permission denied")
)
