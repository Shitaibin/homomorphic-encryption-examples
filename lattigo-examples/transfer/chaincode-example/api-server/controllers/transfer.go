package controllers

import (
	"api-server/service"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type TransferController struct {
	beego.Controller
}

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

// @Title CreateTransfer
// @Description create transfer
// @Param	body		body 	TransferRequest	true		"parameters for transfer"
// @Success 200 {object} TransferResponse
// @Failure 403 body is empty
// @router / [post]
func (t *TransferController) Post() {
	logs.Debug("TransferController.Post")

	defer t.ServeJSON()

	var req TransferRequest
	if err := json.Unmarshal(t.Ctx.Input.RequestBody, &req); err != nil {
		msg := fmt.Sprintf("unmarshal AccountRequest error: %s", err.Error())
		logs.Error(msg)
		t.Data["json"] = msg
		return
	}

	txid, validCode, err := service.Transfer(req.FromBankID, req.FromAccountID, req.ToBankID, req.ToAccountID, req.Amount)
	if err != nil {
		msg := fmt.Sprintf("TransferController error: %s", err)
		logs.Error(msg)
		t.Data["json"] = TransferResponse{
			FromBankID:    req.FromBankID,
			FromAccountID: req.FromAccountID,
			ToBankID:      req.ToBankID,
			ToAccountID:   req.ToAccountID,
			TxID:          txid,
			ValidCode:     validCode,
			Message:       err.Error(),
		}
	} else {
		t.Data["json"] = TransferResponse{
			FromBankID:    req.FromBankID,
			FromAccountID: req.FromAccountID,
			ToBankID:      req.ToBankID,
			ToAccountID:   req.ToAccountID,
			TxID:          txid,
			ValidCode:     validCode,
			Message:       "success",
		}
	}
}
