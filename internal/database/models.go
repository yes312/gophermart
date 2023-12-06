package db

import (
	"errors"
	"time"
)

type User struct {
	Login string
	Hash  string
}

type OrderStatus struct {
	Number     string    `json:"order"`
	Status     string    `json:"status"`
	Accrual    int       `json:"accrual ,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type Balance struct {
	Current  int `json:"current"`
	Withdraw int `json:"withdraw"`
}

type OrderSum struct {
	OrderNumber string `json:"order"`
	Sum         int    `json:"sum"`
}

type OrderUserID struct {
	OrderNumber string `json:"order"`
	UserID      string `json:"user_id"`
}

type Withdrawal struct {
	OrderNumber string    `json:"order"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}

var ErrNotEnoughFunds = errors.New("not enough funds on balance")
