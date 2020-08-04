package controllers

import (
	"api-server/service"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
)

type AccountController struct {
	beego.Controller
}

type AccountRequest struct {
	BankID    string `json:"bankId"`
	AccountID string `json:"accountId"`
	Balance   uint64 `json:"balance"`
}

type AccountResponse struct {
	BankID    string `json:"bankId"`
	AccountID string `json:"accountId"`
	Balance   uint64 `json:"balance"`
	Message   string `json:"msg"` // 错误信息
}

// @router /v1/account/:id [get]
func (a *AccountController) Get() {
	logs.Debug("AccountController.Get")

	defer a.ServeJSON()

	bid := a.GetString("bankId")
	aid := a.GetString("accountId")

	if len(bid) <= 0 || len(aid) <= 0 {
		msg := fmt.Sprintf("not enough params for AccountController, bid = %v, aid = %v", bid, aid)
		logs.Error(msg)
		a.Data["json"] = msg
		return
	}

	bal, err := service.GetAccountBalance(bid, aid)
	if err != nil {
		msg := fmt.Sprintf("GetAccountBalance error: %s", err)
		logs.Error(msg)
		a.Data["json"] = AccountResponse{
			BankID:    bid,
			AccountID: aid,
			Balance:   bal,
			Message:   msg,
		}
	} else {
		a.Data["json"] = AccountResponse{
			BankID:    bid,
			AccountID: aid,
			Balance:   bal,
		}
	}
}

// @router /v1/account [post]
func (a *AccountController) Post() {
	a.Put()
}

// @router /v1/account [put]
func (a *AccountController) Put() {
	logs.Debug("AccountController.Put")

	defer a.ServeJSON()

	var req AccountRequest
	if err := json.Unmarshal(a.Ctx.Input.RequestBody, &req); err != nil {
		msg := fmt.Sprintf("unmarshal AccountRequest error: %s", err.Error())
		logs.Error(msg)
		a.Data["json"] = msg
		return
	}

	if err := service.SetAccountBalance(req.BankID, req.AccountID, req.Balance); err != nil {
		msg := fmt.Sprintf("SetAccountBalance error: %s", err.Error())
		logs.Error(msg)
		a.Data["json"] = msg
	} else {
		a.Data["json"] = "Set account balance success."
	}
}
