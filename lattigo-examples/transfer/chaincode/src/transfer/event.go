package main

import (
	"encoding/json"
	"fmt"
)

const ChainCodeEventName_Transfer = "ChainCodeEvent_Transfer"

type ChainCodeEvent_Transfer struct {
	From        string
	FromAccount string
	To          string
	ToAccount   string
	Amount      string
	Msg         string // 转账描述信息
}

func NewMarshaledTransferEvent(from, facc, to, tacc, amount string) []byte {
	e := ChainCodeEvent_Transfer{From: from, FromAccount: facc, To: to, ToAccount: tacc, Amount: amount}
	b, err := json.Marshal(e)
	if err != nil {
		fmt.Println("error:", err)
	}
	return b
}
