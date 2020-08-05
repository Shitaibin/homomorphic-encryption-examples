package models

import (
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type AccountRequest struct {
	BankID    string `json:"bankId"`
	AccountID string `json:"accountId"`
	Balance   uint64 `json:"balance"`
}

type AccountResponse struct {
	BankID    string                `json:"bankId"`
	AccountID string                `json:"accountId"`
	Balance   uint64                `json:"balance"`
	TxID      fab.TransactionID     `json:"txId"`
	ValidCode peer.TxValidationCode `json:"validCode"`
	Message   string                `json:"msg"` // 错误信息
}
