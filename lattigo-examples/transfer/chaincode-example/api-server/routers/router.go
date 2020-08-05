// @APIVersion 1.0.0
// @Title Cross Bank Transfer API
// @Description 基于同态加密的跨行转账API
// @Contact hz_stb@163.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"api-server/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/v1",
		// beego.NSNamespace("/object",
		// 	beego.NSInclude(
		// 		&controllers.ObjectController{},
		// 	),
		// ),
		// beego.NSNamespace("/user",
		// 	beego.NSInclude(
		// 		&controllers.UserController{},
		// 	),
		// ),
		beego.NSNamespace("/bank",
			beego.NSInclude(
				&controllers.BankController{},
			),
		),
		beego.NSNamespace("/account",
			beego.NSInclude(
				&controllers.AccountController{},
			),
		),
		beego.NSNamespace("/transfer",
			beego.NSInclude(
				&controllers.TransferController{},
			),
		),
	)
	beego.AddNamespace(ns)
}
