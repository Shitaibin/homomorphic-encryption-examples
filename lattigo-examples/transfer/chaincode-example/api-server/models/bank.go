package models

import (
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
)

type AddPublicRequest struct {
	BankID string `json:"bankId"`
	// PublicKey string `json:"publicKey"`
}

type Bank struct {
	BankID    string                `json:"bankId"`
	TxID      fab.TransactionID     `json:"txId"`
	ValidCode peer.TxValidationCode `json:"validCode"`
	Message   string                `json:"msg"` // 错误信息
}
