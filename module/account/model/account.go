package model

import (
	"fmt"
	"strings"
	"time"
)

type Account struct {
	ID        int64     `json:"id"`
	OwnerName string    `json:"owner_name"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAccountRequest struct {
	OwnerName string `json:"owner_name"`
	Balance   int64  `json:"balance"`
}

func (r *CreateAccountRequest) Validate() error {
	r.OwnerName = strings.TrimSpace(r.OwnerName)

	if r.OwnerName == "" {
		return fmt.Errorf("owner_name is required")
	}

	if r.Balance < 0 {
		return fmt.Errorf("balance must be >= 0")
	}

	return nil
}

type ListAccountsQuery struct {
	Limit  int32 `form:"limit"`
	Offset int32 `form:"offset"`
}

func (q *ListAccountsQuery) Normalize() error {
	if q.Limit == 0 {
		q.Limit = 10
	}

	if q.Limit < 0 {
		return fmt.Errorf("limit must be >= 0")
	}
	if q.Offset < 0 {
		return fmt.Errorf("offset must be >= 0")
	}

	if q.Limit > 100 {
		q.Limit = 100
	}

	return nil
}
