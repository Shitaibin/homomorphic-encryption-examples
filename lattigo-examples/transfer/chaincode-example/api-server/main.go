package main

import (
	_ "api-server/routers"
	"api-server/service"
	"os"

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
	if len(os.Args) <= 1 || os.Args[1] == "1" {
		service.CLI = service.NewOrg1Peer0Client()
	} else {
		service.CLI = service.NewOrg2Peer0Client()
	}

	beego.Run()

	service.CLI.Close()
	logs.Info("CLI closed.")
}
