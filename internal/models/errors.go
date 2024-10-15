package models

import "errors"

var (
	ErrNoRows    = errors.New("row not found")
	ErrDuplicate = errors.New("duplicate item")
	ErrExist     = errors.New("item already exist")

	ErrSessionEmpty = errors.New("user session not found")
)
