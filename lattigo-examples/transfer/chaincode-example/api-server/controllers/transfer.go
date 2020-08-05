package controllers

import (
	"api-server/models"
	"api-server/service"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type TransferController struct {
	beego.Controller
}

// @Title CreateTransfer
// @Description create transfer
// @Param	body		body 	models.TransferRequest	true		"parameters for transfer"
// @Success 200 {object} models.TransferResponse
// @Failure 403 body is empty
// @router / [post]
func (t *TransferController) Post() {
	logs.Debug("TransferController.Post")

	defer t.ServeJSON()

	var req models.TransferRequest
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
		t.Data["json"] = models.TransferResponse{
			FromBankID:    req.FromBankID,
			FromAccountID: req.FromAccountID,
			ToBankID:      req.ToBankID,
			ToAccountID:   req.ToAccountID,
			TxID:          txid,
			ValidCode:     validCode,
			Message:       err.Error(),
		}
	} else {
		t.Data["json"] = models.TransferResponse{
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
