package controllers

import (
	"api-server/service"
	"encoding/json"
	"fmt"

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
	FromBankID    string `json:"fromBankId"`
	FromAccountID string `json:"fromAccountId"`
	ToBankID      string `json:"toBankId"`
	ToAccountID   string `json:"toAccountId"`
	Message       string `json:"msg"` // 错误信息
}

// @router /v1/transfer [post]
func (t *TransferController) Post() {
	logs.Debug("AccountController.Get")

	defer t.ServeJSON()

	var req TransferRequest
	if err := json.Unmarshal(t.Ctx.Input.RequestBody, &req); err != nil {
		msg := fmt.Sprintf("unmarshal AccountRequest error: %s", err.Error())
		logs.Error(msg)
		t.Data["json"] = msg
	}

	err := service.Transfer(req.FromBankID, req.FromAccountID, req.ToBankID, req.ToAccountID, req.Amount)
	if err != nil {
		msg := fmt.Sprintf("GetAccountBalance error: %s", err)
		logs.Error(msg)
		t.Data["json"] = msg
	} else {
		t.Data["json"] = TransferResponse{
			FromBankID:    req.FromBankID,
			FromAccountID: req.FromAccountID,
			ToBankID:      req.ToBankID,
			ToAccountID:   req.ToAccountID,
			Message:       err.Error(),
		}
	}
}
