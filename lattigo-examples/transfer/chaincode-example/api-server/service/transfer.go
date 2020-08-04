package service

import "github.com/pkg/errors"

func Transfer(FromBankID, FromAccountID, ToBankID, ToAccountID string, Amount float64) error {
	// todo
	return errors.New("mock transfer")
}
