package model

import (
	"errors"
	"fmt"
)

type AccountID int64

type Account struct {
	ID      AccountID
	Balance int64
}

func (a *Account) String() string {
	return fmt.Sprintf("Account(%d, %d)", a.ID, a.Balance)
}

var (
	ErrNoAccountsWithBalance = errors.New("no accounts with balance")
)
