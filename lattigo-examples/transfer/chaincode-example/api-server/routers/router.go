// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"api-server/controllers"

	"github.com/astaxie/beego"
)

func init() {
	/*ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/object",
			beego.NSInclude(
				&controllers.ObjectController{},
			),
		),
		beego.NSNamespace("/user",
			beego.NSInclude(
				&controllers.UserController{},
			),
		),
		// beego.NSNamespace("/bank",
		// 	beego.NSInclude(
		// 		&controllers.BankController{},
		// 	),
		// ),
		// beego.NSNamespace("/account",
		// 	beego.NSInclude(
		// 		&controllers.AccountController{},
		// 	),
		// ),
	)
	beego.AddNamespace(ns)*/

	// todo 换成注解路由
	beego.Router("/v1/bank", &controllers.BankController{})
	beego.Router("/v1/account", &controllers.AccountController{})
	beego.Router("/v1/transfer", &controllers.TransferController{})
}
