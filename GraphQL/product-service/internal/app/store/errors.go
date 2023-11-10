package store

import "errors"

// ErrRecordNotFound представляет собой ошибку, возвращаемую при попытке получения записи, которая не существует в хранилище.
var ErrRecordNotFound = errors.New("record not found")
