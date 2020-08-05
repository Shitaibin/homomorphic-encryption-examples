package models

import (
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type TransferRequest struct {
	FromBankID    string  `json:"fromBankId"`
	FromAccountID string  `json:"fromAccountId"`
	ToBankID      string  `json:"toBankId"`
	ToAccountID   string  `json:"toAccountId"`
	Amount        float64 `json:"amount"`
}

type TransferResponse struct {
	FromBankID    string                `json:"fromBankId"`
	FromAccountID string                `json:"fromAccountId"`
	ToBankID      string                `json:"toBankId"`
	ToAccountID   string                `json:"toAccountId"`
	TxID          fab.TransactionID     `json:"txId"`
	ValidCode     peer.TxValidationCode `json:"validCode"`
	Message       string                `json:"msg"` // 错误信息
}
