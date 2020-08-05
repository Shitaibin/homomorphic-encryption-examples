package controllers

import (
	"api-server/models"
	"api-server/service"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
)

type AccountController struct {
	beego.Controller
}

// @Title GetAccount
// @Description get account by bankId and accountId
// @Param	bankId		path 	string	true		"The key for get account"
// @Param	accountId		path 	string	true		"The key for get account"
// @Success 200 {object} models.AccountResponse
// @Failure 403 body is empty
// @router /:accountId/bank/:bankId [get]
func (a *AccountController) Get() {
	logs.Debug("AccountController.Get")

	defer a.ServeJSON()

	// 方法1
	bid := a.GetString(":bankId")
	aid := a.GetString(":accountId")

	// 方法2
	// bid := a.Ctx.Input.Param(":bankId")
	// aid := a.Ctx.Input.Param(":accountId")

	if len(bid) <= 0 || len(aid) <= 0 {
		msg := fmt.Sprintf("not enough params for AccountController, bid = %v, aid = %v", bid, aid)
		logs.Error(msg)
		a.Data["json"] = msg
		return
	}

	bal, err := service.GetAccountBalance(bid, aid)
	if err != nil {
		msg := fmt.Sprintf("GetAccountBalance error: %v", err)
		logs.Error(msg)
		a.Data["json"] = models.AccountResponse{
			BankID:    bid,
			AccountID: aid,
			Balance:   bal,
			Message:   msg,
		}
	} else {
		a.Data["json"] = models.AccountResponse{
			BankID:    bid,
			AccountID: aid,
			Balance:   bal,
			Message:   "success",
		}
	}
}

// @Title Set Account Balance
// @Description create users
// @Param	body		body 	models.AccountRequest	true		"set account balance parameters"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (a *AccountController) Post() {
	a.Put()
}

func (a *AccountController) Put() {
	logs.Debug("AccountController.Put")

	defer a.ServeJSON()

	var req models.AccountRequest
	if err := json.Unmarshal(a.Ctx.Input.RequestBody, &req); err != nil {
		msg := fmt.Sprintf("unmarshal AccountRequest error: %s", err.Error())
		logs.Error(msg)
		a.Data["json"] = msg
		return
	}

	// if err := service.SetAccountBalance(req.BankID, req.AccountID, req.Balance); err != nil {
	// 	msg := fmt.Sprintf("SetAccountBalance error: %s", err.Error())
	// 	logs.Error(msg)
	// 	a.Data["json"] = msg
	// } else {
	// 	a.Data["json"] = "Set account balance success."
	// }

	txid, validCode, err := service.SetAccountBalance(req.BankID, req.AccountID, req.Balance)
	if err != nil {
		msg := fmt.Sprintf("GetAccountBalance error: %v", err)
		logs.Error(msg)
		a.Data["json"] = models.AccountResponse{
			BankID:    req.BankID,
			AccountID: req.AccountID,
			Balance:   req.Balance,
			TxID:      txid,
			ValidCode: validCode,
			Message:   msg,
		}
	} else {
		a.Data["json"] = models.AccountResponse{
			BankID:    req.BankID,
			AccountID: req.AccountID,
			Balance:   req.Balance,
			TxID:      txid,
			ValidCode: validCode,
			Message:   "success",
		}
	}
}
