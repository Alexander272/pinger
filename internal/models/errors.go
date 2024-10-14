package models

import "errors"

var (
	ErrNoRows = errors.New("row not found")

	ErrSessionEmpty = errors.New("user session not found")
)
