package store

import "errors"

// ErrRecordNotFound представляет собой ошибку, указывающую на отсутствие записи в хранилище.
var ErrRecordNotFound = errors.New("record not found")
