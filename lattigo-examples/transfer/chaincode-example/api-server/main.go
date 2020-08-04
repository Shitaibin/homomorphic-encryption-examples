package main

import (
	_ "api-server/routers"
	"api-server/service"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	logs.SetLogger(logs.AdapterFile, `{"filename":"project.log","maxlines":0,"maxsize":0,"daily":true,"maxdays":10}`)

	// 创建客户端
	service.CLI = service.NewOrg1Peer1Client()

	beego.Run()

	service.CLI.Close()
	logs.Info("CLI closed.")
}
