package model

import "errors"

var (
	ErrInvalidRequest = errors.New("invalid request")

	ErrAccountNotFound     = errors.New("account not found")
	ErrInsufficientFunds   = errors.New("insufficient funds")
	ErrSameAccountTransfer = errors.New("cannot transfer to the same account")

	ErrDBOperation = errors.New("database operation failed")
)
