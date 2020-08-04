package controllers

import (
	"encoding/json"
	"fmt"

	"github.com/astaxie/beego"
)

type BankController struct {
	*beego.Controller
}

type AddPublicRequest struct {
	BankID    string `json:"bankId"`
	PublicKey string `json:"publicKey"`
}

// @router /v1/bank/publickey [post]
func (b *BankController) PublicKey() {
	var req AddPublicRequest
	if err := json.Unmarshal(b.Ctx.Input.RequestBody, &req); err != nil {
		b.Ctx.WriteString(fmt.Sprintf("unmarshal AddPublicRequest error: %s", err.Error()))
	}

	// 	创建go的client，调用链码，把链码结果生成json返回
	b.ServeJSON()
}
