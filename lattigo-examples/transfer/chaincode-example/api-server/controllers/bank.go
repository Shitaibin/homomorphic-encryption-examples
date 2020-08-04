package controllers

import (
	"api-server/service"
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type BankController struct {
	beego.Controller
}

type AddPublicRequest struct {
	BankID string `json:"bankId"`
	// PublicKey string `json:"publicKey"`
}

// @Title CreateUser
// @Description create users
// @Param	body		body 	models.User	true		"body for user content"
// @Success 200 {int} models.User.Id
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

	var req AddPublicRequest
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
