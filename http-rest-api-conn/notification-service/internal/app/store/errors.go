package store

import "errors"

// ErrRecordNotFound представляет ошибку "запись не найдена".
var (
	ErrRecordNotFound = errors.New("record not found")
)
