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

	"github.com/astaxie/beego/plugins/cors"

	"github.com/astaxie/beego"
)

func init() {
	// beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
	// 	// AllowAllOrigins:  true,
	// 	AllowMethods:     []string{"POST", "GET"},
	// 	AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "x-requested-with", "no-referrer-when-downgrade"},
	// 	ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type", "Access-Control-Allow-Origin"},
	// 	AllowCredentials: true,
	// 	AllowOrigins:     []string{"*"},
	// }))

	//InsertFilter是提供一个过滤函数
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		// 允许访问所有源
		AllowAllOrigins: true,
		// 可选参数"GET", "POST", "PUT", "DELETE", "OPTIONS" (*为所有)
		// 其中Options跨域复杂请求预检
		AllowMethods: []string{"*"},
		// 指的是允许的Header的种类
		AllowHeaders: []string{"*"},
		// 公开的HTTP标头列表
		ExposeHeaders: []string{"Content-Length"},
		// 如果设置，则允许共享身份验证凭据，例如cookie
		AllowCredentials: true,
	}))

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
