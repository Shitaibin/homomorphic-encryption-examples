package controllers

import (
	"api-server/service"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"

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
	BankID    string                `json:"bankId"`
	AccountID string                `json:"accountId"`
	Balance   uint64                `json:"balance"`
	TxID      fab.TransactionID     `json:"txId"`
	ValidCode peer.TxValidationCode `json:"validCode"`
	Message   string                `json:"msg"` // 错误信息
}

// @Title GetAccount
// @Description get account by bankId and accountId
// @Param	bankId		query 	string	true		"The key for get account"
// @Param	accountId		query 	string	true		"The key for get account"
// @Success 200 {object} AccountResponse
// @Failure 403 body is empty
// @router / [get]
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
		msg := fmt.Sprintf("GetAccountBalance error: %v", err)
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
			Message:   "success",
		}
	}
}

// @Title Set Account Balance
// @Description create users
// @Param	body		body 	AccountRequest	true		"set account balance parameters"
// @Success 200 {int} models.User.Id
// @Failure 403 body is empty
// @router / [post]
func (a *AccountController) Post() {
	a.Put()
}

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
		a.Data["json"] = AccountResponse{
			BankID:    req.BankID,
			AccountID: req.AccountID,
			Balance:   req.Balance,
			TxID:      txid,
			ValidCode: validCode,
			Message:   msg,
		}
	} else {
		a.Data["json"] = AccountResponse{
			BankID:    req.BankID,
			AccountID: req.AccountID,
			Balance:   req.Balance,
			TxID:      txid,
			ValidCode: validCode,
			Message:   "success",
		}
	}
}
