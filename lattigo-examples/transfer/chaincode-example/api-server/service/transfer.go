package service

import "github.com/pkg/errors"

func Transfer(FromBankID, FromAccountID, ToBankID, ToAccountID string, Amount float64) error {
	return errors.New("mock transfer")
}
