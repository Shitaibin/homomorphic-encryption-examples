package controllers

import (
	"api-server/models"
	"api-server/service"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type BankController struct {
	beego.Controller
}

// @Title CreateBank
// @Description create bank
// @Param	body		body 	models.AddPublicRequest	true		"parameters for create bank account"
// @Success 200 {object} models.Bank
// @Failure 403 body is empty
// @router / [post]
func (b *BankController) Post() {
	// defer func() {
	// 	// 发生宕机时，获取panic传递的上下文并打印
	// 	err := recover()
	// 	switch err.(type) {
	// 	case runtime.Error: // 运行时错误
	// 		logs.Error("runtime error:", err)
	// 	default: // 非运行时错误
	// 		logs.Error("error:", err)
	// 	}
	// }()

	logs.Info("BankController.Post")

	defer b.ServeJSON()

	var req models.AddPublicRequest
	if err := json.Unmarshal(b.Ctx.Input.RequestBody, &req); err != nil {
		msg := fmt.Sprintf("unmarshal AddPublicRequest error: %s", err.Error())
		logs.Error(msg)
		b.Data["json"] = msg
		return
	}

	// 	创建go的client，调用链码，把链码结果生成json返回
	bank, err := service.NewBank(req.BankID)
	if err != nil {
		msg := fmt.Sprintf("new bank error: %s", err.Error())
		b.Data["json"] = msg
	} else {
		b.Data["json"] = bank
	}
}

// @Title GetBankKeys
// @Description get keys of bank
// @Param	bankId		path 	string	true		"bankId for get bank keys"
// @Success 200 {object} models.User
// @Failure 403 :bankId is empty
// @router /:bankId/key [get]
func (b *BankController) Keys() {
	bid := b.GetString(":bankId")
	if bid == "" {
		b.Ctx.WriteString("bid is empty")
	}

	keyName := bid + ".keys"
	fp := "./keys/" + keyName

	if err := service.SaveBankKeys(bid, fp); err != nil {
		b.Ctx.WriteString(fmt.Sprintf("Download bank keys error: %s", err.Error()))
	}

	b.Ctx.Output.Download(fp, keyName)
}
