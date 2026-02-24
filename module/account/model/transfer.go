package model

import (
	"fmt"
	"time"
)

type Transfer struct {
	ID            int64     `json:"id"`
	FromAccountID int64     `json:"from_account_id"`
	ToAccountID   int64     `json:"to_account_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

type Entry struct {
	ID         int64     `json:"id"`
	AccountID  int64     `json:"account_id"`
	TransferID int64     `json:"transfer_id"`
	Amount     int64     `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

type CreateTransferRequest struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

func (r *CreateTransferRequest) Validate() error {
	if r.FromAccountID <= 0 {
		return fmt.Errorf("from_account_id must be > 0")
	}

	if r.ToAccountID <= 0 {
		return fmt.Errorf("to_account_id must be > 0")
	}

	if r.FromAccountID == r.ToAccountID {
		return fmt.Errorf("from_account_id and to_account_id must be different")
	}

	if r.Amount <= 0 {
		return fmt.Errorf("amount must be > 0")
	}

	return nil
}
